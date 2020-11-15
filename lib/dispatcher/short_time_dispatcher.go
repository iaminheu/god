package dispatcher

import (
	"git.zc0901.com/go/god/lib/syncx"
	"git.zc0901.com/go/god/lib/timex"
	"time"
)

// ShortTimeDispatcher 短时分发器，仅在指定时间段内执行分发任务
type ShortTimeDispatcher struct {
	duration time.Duration
	lastTime *syncx.AtomicDuration
}

// NewShortTimeDispatcher 仅在指定时间段内执行分发任务
func NewShortTimeDispatcher(duration time.Duration) *ShortTimeDispatcher {
	return &ShortTimeDispatcher{
		duration: duration,
		lastTime: syncx.NewAtomicDuration(),
	}
}

func (d *ShortTimeDispatcher) DoOrDiscard(execute func()) bool {
	now := timex.Now()
	lastTime := d.lastTime.Load()
	if lastTime == 0 || lastTime+d.duration < now {
		d.lastTime.Set(now)
		execute()
		return true
	}

	return false
}
