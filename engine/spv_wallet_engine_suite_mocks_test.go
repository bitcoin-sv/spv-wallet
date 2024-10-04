package engine

import (
	"context"

	taskq "github.com/vmihailenco/taskq/v3"

	"github.com/bitcoin-sv/spv-wallet/engine/taskmanager"
)

// taskManagerMock is a base for an empty task manager
type taskManagerMockBase struct{}

func (tm *taskManagerMockBase) Info(context.Context, string, ...interface{}) {}

func (tm *taskManagerMockBase) RegisterTask(string, interface{}) error {
	return nil
}

func (tm *taskManagerMockBase) ResetCron() {}

func (tm *taskManagerMockBase) RunTask(context.Context, *taskmanager.TaskRunOptions) error {
	return nil
}

func (tm *taskManagerMockBase) Tasks() map[string]*taskq.Task {
	return nil
}

func (tm *taskManagerMockBase) Close(context.Context) error {
	return nil
}

func (tm *taskManagerMockBase) Debug(bool) {}

func (tm *taskManagerMockBase) Factory() taskmanager.Factory {
	return taskmanager.FactoryEmpty
}

func (tm *taskManagerMockBase) GetTxnCtx(ctx context.Context) context.Context {
	return ctx
}

func (tm *taskManagerMockBase) IsDebug() bool {
	return false
}

func (tm *taskManagerMockBase) CronJobsInit(taskmanager.CronJobs) error {
	return nil
}

// Sets custom task manager only for testing
func withTaskManagerMockup() ClientOps {
	return func(c *clientOptions) {
		c.taskManager.TaskEngine = &taskManagerMockBase{}
	}
}
