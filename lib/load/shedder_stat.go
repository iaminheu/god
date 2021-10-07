package load

import (
	"sync/atomic"
	"time"

	"git.zc0901.com/go/god/lib/logx"
	"git.zc0901.com/go/god/lib/stat"
)

type (
	// ShedderStat 用于降低负载的统计项
	ShedderStat struct {
		name  string
		total int64
		pass  int64
		drop  int64
	}

	snapshot struct {
		Total int64
		Pass  int64
		Drop  int64
	}
)

// NewShedderStat 新建降载统计项
func NewShedderStat(name string) *ShedderStat {
	st := &ShedderStat{name: name}
	go st.run()
	return st
}

func (s *ShedderStat) IncrTotal() {
	atomic.AddInt64(&s.total, 1)
}

func (s *ShedderStat) IncrPass() {
	atomic.AddInt64(&s.pass, 1)
}

func (s *ShedderStat) IncrDrop() {
	atomic.AddInt64(&s.drop, 1)
}

func (s *ShedderStat) loop(c <-chan time.Time) {
	for range c {
		snapshot := s.reset()

		if !enabled.True() {
			continue
		}

		cpuUsage := stat.CpuUsage()
		if snapshot.Drop == 0 {
			logx.Statf("(%s) 负载统计 [1m], CPU: %d, 总请求: %d, 已通过: %d, 已丢弃: %d",
				s.name, cpuUsage, snapshot.Total, snapshot.Pass, snapshot.Drop)
		} else {
			logx.Statf("(%s) 降载统计 [1m], CPU: %d, 总请求: %d, 已通过: %d, 已丢弃: %d",
				s.name, cpuUsage, snapshot.Total, snapshot.Pass, snapshot.Drop)
		}
	}
}

func (s *ShedderStat) reset() snapshot {
	return snapshot{
		Total: atomic.SwapInt64(&s.total, 0),
		Pass:  atomic.SwapInt64(&s.pass, 0),
		Drop:  atomic.SwapInt64(&s.drop, 0),
	}
}

func (s *ShedderStat) run() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	s.loop(ticker.C)
}
