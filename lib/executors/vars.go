package executors

import "time"

const defaultFlushInterval = time.Second

// Execute 是执行任务的方法
type Execute func(tasks []interface{})
