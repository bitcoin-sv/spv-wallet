package taskmanager

import (
	"context"

	taskq "github.com/vmihailenco/taskq/v3"
)

// TaskEngine is the taskmanager client interface
type TaskEngine interface {
	RegisterTask(name string, handler interface{}) error
	ResetCron()
	RunTask(ctx context.Context, options *TaskRunOptions) error
	Tasks() map[string]*taskq.Task
	CronJobsInit(cronJobsMap CronJobs) error
	Close(ctx context.Context) error
	Factory() Factory
	GetTxnCtx(ctx context.Context) context.Context
	IsNewRelicEnabled() bool
}
