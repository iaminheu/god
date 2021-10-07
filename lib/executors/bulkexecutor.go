package executors

import (
	"time"
)

const defaultBulkTasks = 1000

type (
	// BulkExecutor 批量执行器是基于下述条件的执行器：
	// 1. 任务数量达到指定大小
	// 2. 时间触达指定的间隔
	BulkExecutor struct {
		executor  *PeriodicalExecutor
		container *bulkContainer
	}

	// BulkOption 是 BulkExecutor的 自定义方法。
	BulkOption func(options *bulkOptions)

	bulkOptions struct {
		batchSize     int           // 批次数量
		flushInterval time.Duration // 批量执行时间间隔
	}
)

func NewBulkExecutor(execute Execute, opts ...BulkOption) *BulkExecutor {
	options := newBulkOptions()
	for _, opt := range opts {
		opt(&options)
	}

	container := &bulkContainer{
		execute:   execute,
		batchSize: options.batchSize,
	}
	executor := &BulkExecutor{
		executor:  NewPeriodicalExecutor(options.flushInterval, container),
		container: container,
	}

	return executor
}

// Add 添加任务到批量执行器。
func (b *BulkExecutor) Add(task interface{}) error {
	b.executor.Add(task)
	return nil
}

// Flush 强制刷新并执行任务。
func (b *BulkExecutor) Flush() {
	b.executor.Flush()
}

// Wait 等待任务执行完成。
func (b *BulkExecutor) Wait() {
	b.executor.Wait()
}

func WithBulkSize(size int) BulkOption {
	return func(options *bulkOptions) {
		options.batchSize = size
	}
}

func WithBulkInterval(duration time.Duration) BulkOption {
	return func(options *bulkOptions) {
		options.flushInterval = duration
	}
}

func newBulkOptions() bulkOptions {
	return bulkOptions{
		batchSize:     defaultBulkTasks,
		flushInterval: defaultFlushInterval,
	}
}

type bulkContainer struct {
	tasks     []interface{}
	execute   Execute
	batchSize int
}

func (b *bulkContainer) Add(task interface{}) bool {
	b.tasks = append(b.tasks, task)
	return len(b.tasks) >= b.batchSize
}

func (b *bulkContainer) Execute(tasks interface{}) {
	v := tasks.([]interface{})
	b.execute(v)
}

func (b *bulkContainer) RemoveAll() interface{} {
	tasks := b.tasks
	b.tasks = nil
	return tasks
}
