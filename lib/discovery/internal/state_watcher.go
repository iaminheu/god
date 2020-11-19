//go:generate mockgen -package internal -destination state_watcher_mock.go -source state_watcher.go etcdConn
package internal

import (
	"context"
	"google.golang.org/grpc/connectivity"
	"sync"
)

type (
	etcdConn interface {
		GetState() connectivity.State
		WaitForStateChange(ctx context.Context, sourceState connectivity.State) bool
	}

	stateWatcher struct {
		disconnected bool
		currentState connectivity.State
		listeners    []func()
		lock         sync.Mutex
	}
)

func newStateWatcher() *stateWatcher {
	return new(stateWatcher)
}

func (sw *stateWatcher) addListener(l func()) {
	sw.lock.Lock()
	sw.listeners = append(sw.listeners, l)
	sw.lock.Unlock()
}

func (sw *stateWatcher) watch(conn etcdConn) {
	sw.currentState = conn.GetState()
	for {
		if conn.WaitForStateChange(context.Background(), sw.currentState) {
			newState := conn.GetState()
			sw.lock.Lock()
			sw.currentState = newState

			switch newState {
			case connectivity.TransientFailure, connectivity.Shutdown:
				sw.disconnected = true
			case connectivity.Ready:
				if sw.disconnected {
					sw.disconnected = false
					for _, listener := range sw.listeners {
						listener()
					}
				}
			}

			sw.lock.Unlock()
		}
	}
}
