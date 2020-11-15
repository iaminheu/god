package logx

import (
	"git.zc0901.com/go/god/lib/syncx"
	"git.zc0901.com/go/god/lib/timex"
	"sync/atomic"
	"time"
)

type shortTimeLogger struct {
	duration  time.Duration
	lastTime  *syncx.AtomicDuration
	discarded uint32
}

func newShortTimeLogger(milliseconds int) *shortTimeLogger {
	return &shortTimeLogger{
		duration: time.Duration(milliseconds) * time.Millisecond,
		lastTime: syncx.NewAtomicDuration(),
	}
}

func (le *shortTimeLogger) logOrDiscard(execute func()) {
	if le == nil || le.duration <= 0 {
		execute()
		return
	}

	now := timex.Now()
	if now-le.lastTime.Load() <= le.duration {
		atomic.AddUint32(&le.discarded, 1)
	} else {
		le.lastTime.Set(now)
		discarded := atomic.SwapUint32(&le.discarded, 0)
		if discarded > 0 {
			Errorf("放弃 %d 个错误信息", discarded)
		}

		execute()
	}
}
