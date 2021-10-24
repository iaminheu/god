package syncx

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestCachedCalls_Get(t *testing.T) {
	calls := NewSingleFlight()
	v, hit, err := calls.Do("ping", func() (val interface{}, err error) {
		return "pong", nil
	})
	got := fmt.Sprintf("%v (%T) %t", v, v, hit)
	expect := "pong (string) false"
	if got != expect {
		t.Errorf("Do = %v; expect %v", got, expect)
	}
	if err != nil {
		t.Errorf("Do error = %v", err)
	}
}

func TestCachedCalls_GetErr(t *testing.T) {
	calls := NewSingleFlight()
	someErr := errors.New("一些错误")
	v, _, err := calls.Do("ping", func() (interface{}, error) {
		return nil, someErr
	})
	if err != someErr {
		t.Errorf("返回 error = %v; 期望 someErr", someErr)
	}
	if v != nil {
		t.Errorf("不期望的非空值 %#v", v)
	}
}

func TestExclusiveCallDoDupSuppress(t *testing.T) {
	calls := NewSingleFlight()
	c := make(chan string)
	var callCounter int32
	fn := func() (interface{}, error) {
		atomic.AddInt32(&callCounter, 1)
		return <-c, nil
	}

	const n = 10
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			v, _, err := calls.Do("ping", fn)
			if err != nil {
				t.Errorf("执行失败：%val", err)
			}
			if v.(string) != "pong" {
				t.Errorf("得到 %q; 期望 %q", v, "pong")
			}
			wg.Done()
		}()
	}
	time.Sleep(100 * time.Millisecond) // 阻塞上述 goroutines
	c <- "pong"
	wg.Wait()
	if got := atomic.LoadInt32(&callCounter); got != 1 {
		t.Errorf("实际调用次数 = %d; 期望 1", got)
	}
}

func TestExclusiveCallDoDiffDupSuppress(t *testing.T) {
	g := NewSingleFlight()
	broadcast := make(chan struct{})
	var calls int32
	tests := []string{"e", "a", "e", "a", "b", "c", "b", "a", "c", "d", "b", "c", "d"}

	var wg sync.WaitGroup
	for _, key := range tests {
		wg.Add(1)
		go func(k string) {
			<-broadcast // get all goroutines ready
			_, _, err := g.Do(k, func() (interface{}, error) {
				atomic.AddInt32(&calls, 1)
				time.Sleep(10 * time.Millisecond)
				return nil, nil
			})
			if err != nil {
				t.Errorf("Do error: %v", err)
			}
			wg.Done()
		}(key)
	}

	time.Sleep(100 * time.Millisecond) // let goroutines above block
	close(broadcast)
	wg.Wait()

	if got := atomic.LoadInt32(&calls); got != 5 {
		// five letters
		t.Errorf("number of calls = %d; want 5", got)
	}
}

func TestExclusiveCallDoExDupSuppress(t *testing.T) {
	g := NewSingleFlight()
	c := make(chan string)
	var calls int32
	fn := func() (interface{}, error) {
		atomic.AddInt32(&calls, 1)
		return <-c, nil
	}

	const n = 10
	var wg sync.WaitGroup
	var hits int32
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			v, hit, err := g.Do("key", fn)
			if err != nil {
				t.Errorf("Do error: %v", err)
			}
			if hit {
				atomic.AddInt32(&hits, 1)
			}
			if v.(string) != "bar" {
				t.Errorf("got %q; want %q", v, "bar")
			}
			wg.Done()
		}()
	}
	time.Sleep(100 * time.Millisecond) // let goroutines above block
	c <- "bar"
	wg.Wait()
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Errorf("number of calls = %d; want 1", got)
	}
	if got := atomic.LoadInt32(&hits); got != 9 {
		t.Errorf("hits = %d; want 1", got)
	}
}
