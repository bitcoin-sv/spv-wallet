/*
Package taskmanager is the task/job management service layer for concurrent and asynchronous tasks with cron scheduling.
*/
package taskmanager

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/logging"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	taskq "github.com/vmihailenco/taskq/v3"
)

type (

	// TaskManager implements the TaskEngine interface
	TaskManager struct {
		options *options
	}

	options struct {
		cronService     *cron.Cron      // Internal cron job client
		logger          *zerolog.Logger // Internal logging
		newRelicEnabled bool            // If NewRelic is enabled (parent application)
		taskq           *taskqOptions   // All configuration and options for using TaskQ
	}

	// taskqOptions holds all the configuration for the TaskQ engine
	taskqOptions struct {
		config *taskq.QueueOptions    // Configuration for the TaskQ engine
		queue  taskq.Queue            // Queue for TaskQ
		tasks  map[string]*taskq.Task // Registered tasks
	}
)

// NewTaskManager creates a new client for all TaskManager functionality
// If no options are given, it will use local memory for the queue.
// ctx may contain a NewRelic txn (or one will be created)
func NewTaskManager(ctx context.Context, opts ...TaskManagerOptions) (TaskEngine, error) {
	// Create a new tm with defaults
	tm := &TaskManager{options: &options{
		newRelicEnabled: false,
		taskq: &taskqOptions{
			tasks:  make(map[string]*taskq.Task),
			config: DefaultTaskQConfig("taskq"),
		},
	}}

	// Overwrite defaults with any set by user
	for _, opt := range opts {
		opt(tm.options)
	}

	if tm.options.logger == nil {
		tm.options.logger = logging.GetDefaultLogger()
	}

	// Use NewRelic if it's enabled (use existing txn if found on ctx)
	// ctx = tm.options.getTxnCtx(ctx)

	if err := tm.loadTaskQ(ctx); err != nil {
		return nil, err
	}

	tm.ResetCron()

	return tm, nil
}

// Close the client and any open connections
func (tm *TaskManager) Close(ctx context.Context) error {
	if txn := newrelic.FromContext(ctx); txn != nil {
		defer txn.StartSegment("close_taskmanager").End()
	}
	if tm != nil && tm.options != nil {

		// Stop the cron scheduler
		if tm.options.cronService != nil {
			tm.options.cronService.Stop()
			tm.options.cronService = nil
		}

		// Close the taskq queue
		if err := tm.options.taskq.queue.Close(); err != nil {
			return err
		}

		// Empty all values and reset
		tm.options.taskq.config = nil
		tm.options.taskq.queue = nil
	}

	return nil
}

// ResetCron will reset the cron scheduler and all loaded tasks
func (tm *TaskManager) ResetCron() {
	if tm.options.cronService != nil {
		tm.options.cronService.Stop()
	}
	tm.options.cronService = cron.New()
	tm.options.cronService.Start()
}

// IsNewRelicEnabled will return if new relic is enabled
func (tm *TaskManager) IsNewRelicEnabled() bool {
	return tm.options.newRelicEnabled
}

// Tasks will return the list of tasks
func (tm *TaskManager) Tasks() map[string]*taskq.Task {
	return tm.options.taskq.tasks
}

// Factory will return the factory that is set
func (tm *TaskManager) Factory() Factory {
	if tm.options == nil || tm.options.taskq == nil {
		return FactoryEmpty
	}
	if tm.options.taskq.config.Redis != nil {
		return FactoryRedis
	}
	return FactoryMemory
}

// GetTxnCtx will check for an existing transaction
func (tm *TaskManager) GetTxnCtx(ctx context.Context) context.Context {
	if tm.options.newRelicEnabled {
		txn := newrelic.FromContext(ctx)
		if txn != nil {
			ctx = newrelic.NewContext(ctx, txn)
		}
	}
	return ctx
}
