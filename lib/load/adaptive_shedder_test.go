package load

import (
	"git.zc0901.com/go/god/lib/collection"
	"git.zc0901.com/go/god/lib/logx"
	"git.zc0901.com/go/god/lib/mathx"
	"git.zc0901.com/go/god/lib/stat"
	"git.zc0901.com/go/god/lib/syncx"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	buckets        = 10
	bucketDuration = time.Millisecond * 50
)

func init() {
	stat.SetReporter(nil)
}

func TestNewAdaptiveShedder(t *testing.T) {
	var wg sync.WaitGroup
	var drop int64
	shedder := NewAdaptiveShedder(WithWindow(bucketDuration), WithBuckets(buckets), WithCpuThreshold(100))
	proba := mathx.NewProba()
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 30; i++ {
				promise, err := shedder.Allow()
				if err != nil {
					atomic.AddInt64(&drop, 1)
				} else {
					count := rand.Intn(5)
					time.Sleep(time.Millisecond * time.Duration(count))
					if proba.TrueOnProba(0.01) {
						promise.Fail()
					} else {
						promise.Pass()
					}
				}
			}
		}()
	}
	wg.Wait()
}

func TestPromise_Pass(t *testing.T) {
	passCounter := newRollingWindow()
	for i := 0; i <= 10; i++ {
		passCounter.Add(float64(i * 100))
		time.Sleep(bucketDuration)
	}
	shedder := &adaptiveShedder{
		droppedRecently: syncx.NewAtomicBool(),
		passCounter:     passCounter,
	}
	assert.Equal(t, int64(1000), shedder.maxPass())

	// 重置计数器 - 默认 maxPass 为 1
	passCounter = newRollingWindow()
	shedder = &adaptiveShedder{
		droppedRecently: syncx.NewAtomicBool(),
		passCounter:     passCounter,
	}
	assert.Equal(t, int64(1), shedder.maxPass())
}

func TestPromise_MinRt(t *testing.T) {
	rtCounter := newRollingWindow()
	for i := 0; i < 10; i++ {
		if i > 0 {
			time.Sleep(bucketDuration)
		}
		for j := i*10 + 1; j <= i*10+10; j++ {
			rtCounter.Add(float64(j))
		}
	}
	shedder := &adaptiveShedder{rtCounter: rtCounter}
	assert.Equal(t, float64(6), shedder.minRt())
}

func TestPromise_MaxFlight(t *testing.T) {
	passCounter := newRollingWindow()
	rtCounter := newRollingWindow()
	for i := 0; i < 10; i++ {
		if i > 0 {
			time.Sleep(bucketDuration)
		}
		passCounter.Add(float64((i + 1) * 100))
		for j := i*10 + 1; j <= i*10+10; j++ {
			rtCounter.Add(float64(j))
		}
	}
	shedder := &adaptiveShedder{
		windows:         buckets,
		droppedRecently: syncx.NewAtomicBool(),
		passCounter:     passCounter,
		rtCounter:       rtCounter,
	}
	assert.Equal(t, int64(54), shedder.maxFlight())
}

func TestAdaptiveShedderShouldDrop(t *testing.T) {
	logx.Disable()
	passCounter := newRollingWindow()
	rtCounter := newRollingWindow()
	for i := 0; i < 10; i++ {
		if i > 0 {
			time.Sleep(bucketDuration)
		}
		passCounter.Add(float64((i + 1) * 100))
		for j := i*10 + 1; j <= i*10+10; j++ {
			rtCounter.Add(float64(j))
		}
	}
	shedder := &adaptiveShedder{
		passCounter:     passCounter,
		rtCounter:       rtCounter,
		windows:         buckets,
		dropTime:        syncx.NewAtomicDuration(),
		droppedRecently: syncx.NewAtomicBool(),
	}
	// cpu >=  800, inflight < maxPass
	systemOverloadChecker = func(int64) bool {
		return true
	}
	shedder.avgFlying = 50
	assert.False(t, shedder.shouldDrop())

	// cpu >=  800, inflight > maxPass
	shedder.avgFlying = 80
	shedder.flying = 50
	assert.False(t, shedder.shouldDrop())

	// cpu >=  800, inflight > maxPass
	shedder.avgFlying = 80
	shedder.flying = 80
	assert.True(t, shedder.shouldDrop())

	// cpu < 800, inflight > maxPass
	systemOverloadChecker = func(int64) bool {
		return false
	}
	shedder.avgFlying = 80
	assert.False(t, shedder.shouldDrop())

	// cpu >=  800, inflight < maxPass
	systemOverloadChecker = func(int64) bool {
		return true
	}
	shedder.avgFlying = 80
	shedder.flying = 80
	_, err := shedder.Allow()
	assert.NotNil(t, err)
}

func TestAdaptiveShedderStillHot(t *testing.T) {
	logx.Disable()
	passCounter := newRollingWindow()
	rtCounter := newRollingWindow()
	for i := 0; i < 10; i++ {
		if i > 0 {
			time.Sleep(bucketDuration)
		}
		passCounter.Add(float64((i + 1) * 100))
		for j := i*10 + 1; j <= i*10+10; j++ {
			rtCounter.Add(float64(j))
		}
	}
	shedder := &adaptiveShedder{
		passCounter:     passCounter,
		rtCounter:       rtCounter,
		windows:         buckets,
		dropTime:        syncx.NewAtomicDuration(),
		droppedRecently: syncx.ForAtomicBool(true),
	}
	assert.False(t, shedder.stillHot())
	shedder.dropTime.Set(-coolOffDuration * 2)
	assert.False(t, shedder.stillHot())
}

func BenchmarkAdaptiveShedder_Allow(b *testing.B) {
	logx.Disable()

	bench := func(b *testing.B) {
		var shedder = NewAdaptiveShedder()
		proba := mathx.NewProba()
		for i := 0; i < 6000; i++ {
			p, err := shedder.Allow()
			if err == nil {
				time.Sleep(time.Millisecond)
				if proba.TrueOnProba(0.01) {
					p.Fail()
				} else {
					p.Pass()
				}
			}
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			p, err := shedder.Allow()
			if err == nil {
				p.Pass()
			}
		}
	}

	systemOverloadChecker = func(int64) bool {
		return true
	}
	b.Run("high load", bench)
	systemOverloadChecker = func(int64) bool {
		return false
	}
	b.Run("low load", bench)
}

func newRollingWindow() *collection.RollingWindow {
	return collection.NewRollingWindow(buckets, bucketDuration, collection.IgnoreCurrentBucket())
}
