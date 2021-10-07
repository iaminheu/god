package load

import (
	"errors"
	"fmt"
	"math"
	"sync/atomic"
	"time"

	"git.zc0901.com/go/god/lib/collection"
	"git.zc0901.com/go/god/lib/logx"
	"git.zc0901.com/go/god/lib/stat"
	"git.zc0901.com/go/god/lib/syncx"
	"git.zc0901.com/go/god/lib/timex"
)

const (
	defaultBuckets = 50
	defaultWindow  = 5 * time.Second

	// 1000m 计数法，900m 差不多相当于80%
	defaultCpuThreshold = 900

	// 默认最小响应时间
	defaultMinRt = float64(time.Second / time.Millisecond)

	flyingBeta      = 0.9         // 用于动态计算的超参
	coolOffDuration = time.Second // 冷却时长
)

var (
	ErrServiceOverloaded = errors.New("service overloaded - 服务超载")

	// 是否启用自适应负载泄流器，默认启用
	enabled = syncx.ForAtomicBool(true)

	// 系统超载检测函数（判断CPU用量是否超过预置的阈值）
	systemOverloadChecker = func(cpuThreshold int64) bool {
		return stat.CpuUsage() >= cpuThreshold
	}
)

type (
	Promise interface {
		Pass()
		Fail()
	}

	promise struct {
		start   time.Duration
		shedder *adaptiveShedder
	}

	// Shedder 负载泄流器
	Shedder interface {
		Allow() (Promise, error)
	}

	shedderOptions struct {
		window       time.Duration
		buckets      int
		cpuThreshold int64
	}

	ShedderOption func(opts *shedderOptions)

	// 自适应泄流器
	adaptiveShedder struct {
		cpuThreshold    int64
		windows         int64 // 每秒的buckets
		flying          int64 // 飞行架次
		avgFlying       float64
		avgFlyingLock   syncx.SwitchLock
		dropTime        *syncx.AtomicDuration
		droppedRecently *syncx.AtomicBool
		passCounter     *collection.RollingWindow // 请求通过计数器
		rtCounter       *collection.RollingWindow // 响应时间计数器
	}
)

func NewAdaptiveShedder(opts ...ShedderOption) Shedder {
	if !enabled.True() {
		return newNopShedder()
	}

	options := shedderOptions{
		window:       defaultWindow,
		buckets:      defaultBuckets,
		cpuThreshold: defaultCpuThreshold,
	}
	for _, opt := range opts {
		opt(&options)
	}
	bucketDuration := options.window / time.Duration(options.buckets)
	return &adaptiveShedder{
		cpuThreshold:    options.cpuThreshold,
		windows:         int64(time.Second / bucketDuration),
		dropTime:        syncx.NewAtomicDuration(),
		droppedRecently: syncx.NewAtomicBool(),
		passCounter:     collection.NewRollingWindow(options.buckets, bucketDuration, collection.IgnoreCurrentBucket()),
		rtCounter:       collection.NewRollingWindow(options.buckets, bucketDuration, collection.IgnoreCurrentBucket()),
	}
}

func WithBuckets(buckets int) ShedderOption {
	return func(opts *shedderOptions) {
		opts.buckets = buckets
	}
}

func WithWindow(window time.Duration) ShedderOption {
	return func(opts *shedderOptions) {
		opts.window = window
	}
}

func WithCpuThreshold(cpuThreshold int64) ShedderOption {
	return func(opts *shedderOptions) {
		opts.cpuThreshold = cpuThreshold
	}
}

func Disable() {
	enabled.Set(false)
}

// Allow 判断是否接受请求并进行相关处理。
func (as *adaptiveShedder) Allow() (Promise, error) {
	if as.shouldDrop() {
		as.dropTime.Set(timex.Now())
		as.droppedRecently.Set(true)

		return nil, ErrServiceOverloaded
	}

	as.addFlying(1)

	return &promise{
		start:   timex.Now(),
		shedder: as,
	}, nil
}

// shouldDrop 判断是否删除该请求，若删除则记录日志
func (as *adaptiveShedder) shouldDrop() bool {
	if as.systemOverloaded() || as.stillHot() {
		if as.highThru() {
			flying := atomic.LoadInt64(&as.flying)
			as.avgFlyingLock.Lock()
			avgFlying := as.avgFlying
			as.avgFlyingLock.Unlock()
			msg := fmt.Sprintf("丢弃请求，CPU: %d, maxPass: %d, minRt: %.2f, hot: %t, flying: %d, avgFlying: %.2f",
				stat.CpuUsage(), as.maxPass(), as.minRt(), as.stillHot(), flying, avgFlying)
			logx.Error(msg)
			stat.Report(msg)
			return true
		}
	}

	return false
}

func (as *adaptiveShedder) systemOverloaded() bool {
	return systemOverloadChecker(as.cpuThreshold)
}

func (as *adaptiveShedder) stillHot() bool {
	if !as.droppedRecently.True() {
		return false
	}

	dropTime := as.dropTime.Load()
	if dropTime == 0 {
		return false
	}

	hot := timex.Since(dropTime) < coolOffDuration
	if !hot {
		as.droppedRecently.Set(false)
	}

	return hot
}

func (as *adaptiveShedder) highThru() bool {
	as.avgFlyingLock.Lock()
	avgFlying := as.avgFlying
	as.avgFlyingLock.Unlock()
	maxFlight := as.maxFlight()
	return int64(avgFlying) > maxFlight && atomic.LoadInt64(&as.flying) > maxFlight
}

func (as *adaptiveShedder) maxFlight() int64 {
	// windows = buckets per second
	// maxQPS = maxPASS * windows
	// minRT = 最小平均响应时间(毫秒)
	// maxQPS = minRT / 每秒的毫秒数
	return int64(math.Max(1, float64(as.maxPass()*as.windows)*(as.minRt()/1e3)))
}

// maxPass 最大请求数
func (as *adaptiveShedder) maxPass() int64 {
	var result float64 = 1

	as.passCounter.Reduce(func(b *collection.Bucket) {
		if b.Accepts > result {
			result = b.Accepts
		}
	})

	return int64(result)
}

// minRt 最小平均响应时间(毫秒)
func (as *adaptiveShedder) minRt() float64 {
	result := defaultMinRt

	as.rtCounter.Reduce(func(b *collection.Bucket) {
		if b.Requests <= 0 {
			return
		}

		avg := math.Round(b.Accepts / float64(b.Requests))
		if avg < result {
			result = avg
		}
	})

	return result
}

func (as *adaptiveShedder) addFlying(delta int64) {
	flying := atomic.AddInt64(&as.flying, delta)
	// 当请求完成，更新 avgFlying
	// 该策略让 avgFlying 和 flying 稍有延迟，更加平滑：
	// 1. 当飞行请求快速增加时，avgFlying 增加较慢，可接受更多请求，
	// 2. 当飞行请求大量被删时，avgFlying 慢慢地删，让无效请求更少，
	// 如此，让服务可以尽可能多的接收更多请求。

	if delta < 0 {
		as.avgFlyingLock.Lock()
		as.avgFlying = as.avgFlying*flyingBeta + float64(flying)*(1-flyingBeta)
		as.avgFlyingLock.Unlock()
	}
}

func (p *promise) Pass() {
	rt := float64(timex.Since(p.start)) / float64(time.Millisecond)
	p.shedder.addFlying(-1)
	p.shedder.rtCounter.Add(math.Ceil(rt))
	p.shedder.passCounter.Add(1)
}

func (p *promise) Fail() {
	p.shedder.addFlying(-1)
}
