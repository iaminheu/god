package executors

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBulkExecutor_Add(t *testing.T) {
	var values []int
	var lock sync.Mutex

	executor := NewBulkExecutor(func(tasks []interface{}) {
		lock.Lock()
		values = append(values, len(tasks))
		lock.Unlock()
	}, WithBulkSize(10), WithBulkInterval(time.Minute))

	for i := 0; i < 50; i++ {
		executor.Add(1)
		time.Sleep(time.Millisecond)
	}

	lock.Lock()
	assert.True(t, len(values) > 0)
	for i := 0; i < len(values); i++ {
		assert.Equal(t, 10, values[i])
	}
	lock.Unlock()
}

func TestBulkExecutor_Flush(t *testing.T) {
	const caches = 10
	const size = 5
	var wait sync.WaitGroup

	wait.Add(1)
	executor := NewBulkExecutor(func(tasks []interface{}) {
		assert.Equal(t, size, len(tasks))
		wait.Done()
	}, WithBulkSize(caches), WithBulkInterval(time.Millisecond))

	for i := 0; i < size; i++ {
		executor.Add(1)
	}
	wait.Wait()
}

func TestBulkExecutor_FlushSlowTasks(t *testing.T) {
	const total = 1500
	lock := new(sync.Mutex)
	result := make([]interface{}, 0, 10000)
	executor := NewBulkExecutor(func(tasks []interface{}) {
		time.Sleep(time.Millisecond * 100)
		lock.Lock()
		defer lock.Unlock()
		result = append(result, tasks...)
	}, WithBulkSize(1000))
	for i := 0; i < total; i++ {
		assert.Nil(t, executor.Add(i))
	}

	executor.Flush()
	executor.Wait()
	assert.Equal(t, total, len(result))
}

func BenchmarkBulkExecutor(b *testing.B) {
	b.ReportAllocs()

	executor := NewBulkExecutor(func(tasks []interface{}) {
		time.Sleep(time.Millisecond * time.Duration(len(tasks)))
	})
	for i := 0; i < b.N; i++ {
		time.Sleep(time.Microsecond * 200)
		executor.Add(1)
	}
	executor.Flush()
}
