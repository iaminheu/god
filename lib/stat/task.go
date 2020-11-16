package stat

import "time"

// Task 统计任务项
type Task struct {
	Drop        bool
	Duration    time.Duration
	Description string
}
