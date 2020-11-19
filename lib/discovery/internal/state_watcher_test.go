package internal

import (
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/connectivity"
	"sync"
	"testing"
)

func TestStateWatcher_watch(t *testing.T) {
	watcher := newStateWatcher()
	var wg sync.WaitGroup
	wg.Add(1)
	watcher.addListener(func() {
		wg.Done()
	})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	conn := NewMocketcdConn(ctrl)
	conn.EXPECT().GetState().Return(connectivity.Ready)
	conn.EXPECT().GetState().Return(connectivity.TransientFailure)
	conn.EXPECT().GetState().Return(connectivity.Ready).AnyTimes()
	conn.EXPECT().WaitForStateChange(gomock.Any(), gomock.Any()).Return(true).AnyTimes()
	go watcher.watch(conn)
	wg.Wait()
}
