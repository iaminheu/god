package load

import (
	"git.zc0901.com/go/god/lib/logx"
	"git.zc0901.com/go/god/lib/stat"
	"sync/atomic"
	"time"
)

type (
	// ShedderStat 负载卸流统计
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

// NewShedderStat 新建负载卸流统计
func NewShedderStat(name string) *ShedderStat {
	st := &ShedderStat{name: name}
	go st.run()
	return st
}

func (s ShedderStat) run() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		usage := stat.CpuUsage()
		ss := s.reset()
		if ss.Drop == 0 {
			logx.Statf("(%s) 负载卸流器统计 [1m], cpu: %d, total: %d, pass: %d, drop: %d",
				s.name, usage, ss.Total, ss.Pass, ss.Drop)
		} else {
			logx.Statf("(%s) 负载卸流统计 [1m], cpu: %d, total: %d, pass: %d, drop: %d",
				s.name, usage, ss.Total, ss.Pass, ss.Drop)
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

func (s *ShedderStat) IncrTotal() {
	atomic.AddInt64(&s.total, 1)
}

func (s *ShedderStat) IncrPass() {
	atomic.AddInt64(&s.pass, 1)
}

func (s *ShedderStat) IncrDrop() {
	atomic.AddInt64(&s.drop, 1)
}
