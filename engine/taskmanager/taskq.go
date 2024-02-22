package taskmanager

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis_rate/v9"
	taskq "github.com/vmihailenco/taskq/v3"
	"github.com/vmihailenco/taskq/v3/memqueue"
	"github.com/vmihailenco/taskq/v3/redisq"
)

var mutex sync.Mutex

// TasqOps allow functional options to be supplied
type TasqOps func(*taskq.QueueOptions)

// WithRedis will set the redis client for the TaskQ engine
// Note: Because we use redis/v8, we need to use Redis lower than 7.2.0
func WithRedis(addr string) TasqOps {
	return func(queueOptions *taskq.QueueOptions) {
		queueOptions.Redis = redis.NewClient(&redis.Options{
			Addr: strings.Replace(addr, "redis://", "", -1),
		})
	}
}

// DefaultTaskQConfig will return a QueueOptions with specified name and functional options applied
func DefaultTaskQConfig(name string, opts ...TasqOps) *taskq.QueueOptions {
	queueOptions := &taskq.QueueOptions{
		BufferSize:           10,                      // Size of the buffer where reserved messages are stored.
		ConsumerIdleTimeout:  6 * time.Hour,           // ConsumerIdleTimeout Time after which the consumer need to be deleted.
		Handler:              nil,                     // Optional message handler. The default is the global Tasks registry.
		MaxNumFetcher:        0,                       // Maximum number of goroutines fetching messages.
		MaxNumWorker:         10,                      // Maximum number of goroutines processing messages.
		MinNumWorker:         1,                       // Minimum number of goroutines processing messages.
		Name:                 name,                    // Queue name.
		PauseErrorsThreshold: 100,                     // Number of consecutive failures after which queue processing is paused.
		RateLimit:            redis_rate.Limit{},      // Processing rate limit.
		RateLimiter:          nil,                     // Optional rate limiter. The default is to use Redis.
		Redis:                nil,                     // Redis client that is used for storing metadata.
		ReservationSize:      10,                      // Number of messages reserved by a fetcher in the queue in one request.
		ReservationTimeout:   60 * time.Second,        // Time after which the reserved message is returned to the queue.
		Storage:              taskq.NewLocalStorage(), // Optional storage interface. The default is to use Redis.
		WaitTimeout:          3 * time.Second,         // Time that a long polling receive call waits for a message to become available before returning an empty response.
		WorkerLimit:          0,                       // Global limit of concurrently running workers across all servers. Overrides MaxNumWorker.
	}

	for _, opt := range opts {
		opt(queueOptions)
	}

	return queueOptions
}

// TaskRunOptions are the options for running a task
type TaskRunOptions struct {
	Arguments      []interface{} // Arguments for the task
	RunEveryPeriod time.Duration // Cron job!
	TaskName       string        // Name of the task
}

func (runOptions *TaskRunOptions) runImmediately() bool {
	return runOptions.RunEveryPeriod == 0
}

// loadTaskQ will load TaskQ based on the Factory Type and configuration set by the client loading
func (c *TaskManager) loadTaskQ(ctx context.Context) error {
	// Check for a valid config (set on client creation)
	factoryType := c.Factory()
	if factoryType == FactoryEmpty {
		return fmt.Errorf("missing factory type to load taskq")
	}

	var factory taskq.Factory
	if factoryType == FactoryMemory {
		factory = memqueue.NewFactory()
	} else if factoryType == FactoryRedis {
		factory = redisq.NewFactory()
	}

	// Set the queue
	q := factory.RegisterQueue(c.options.taskq.config)
	c.options.taskq.queue = q
	if factoryType == FactoryRedis {
		if err := q.Consumer().Start(ctx); err != nil {
			return err
		}
	}

	// turn off logger for now
	// NOTE: having issues with logger with system resources
	// taskq.SetLogger(nil)

	return nil
}

// RegisterTask will register a new task to handle asynchronously
func (c *TaskManager) RegisterTask(name string, handler interface{}) (err error) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = fmt.Errorf(fmt.Sprintf("registering task panic: %v", panicErr))
			c.options.logger.Error().Msg(err.Error())
		}
	}()

	mutex.Lock()
	defer mutex.Unlock()

	if t := taskq.Tasks.Get(name); t != nil {
		// if already registered - register the task locally
		c.options.taskq.tasks[name] = t
	} else {
		// Register and store the task
		c.options.taskq.tasks[name] = taskq.RegisterTask(&taskq.TaskOptions{
			Name:       name,
			Handler:    handler,
			RetryLimit: 1,
		})
	}

	c.options.logger.Debug().Msgf("registering task: %s...", c.options.taskq.tasks[name].Name())
	return nil
}

// RunTask will run a task using TaskQ
func (c *TaskManager) RunTask(ctx context.Context, options *TaskRunOptions) error {
	c.options.logger.Info().Msgf("executing task: %s", options.TaskName)

	// Try to get the task
	task, ok := c.options.taskq.tasks[options.TaskName]
	if !ok {
		return fmt.Errorf("task %s not registered", options.TaskName)
	}

	// Task message will be used to add to the queue
	taskMessage := task.WithArgs(ctx, options.Arguments...)

	if options.runImmediately() {
		return c.options.taskq.queue.Add(taskMessage)
	}
	// Note: The first scheduled run will be after the period has passed
	return c.scheduleTaskWithCron(ctx, task, taskMessage, options.RunEveryPeriod)
}

func (c *TaskManager) scheduleTaskWithCron(ctx context.Context, task *taskq.Task, taskMessage *taskq.Message, runEveryPeriod time.Duration) error {
	// When using Redis, we need to use a distributed timed lock to prevent the addition of the same task to the queue by multiple instances.
	// With this approach, only one instance will add the task to the queue within a given period.
	var tryLock func() bool
	if c.Factory() == FactoryRedis {
		key := fmt.Sprintf("taskq_cronlock_%s", task.Name())

		// The runEveryPeriod should be greater than 1 second
		if runEveryPeriod < 1*time.Second {
			return fmt.Errorf("runEveryPeriod should be greater than 1 second")
		}

		// Lock time is the period minus 500ms to allow for some clock drift
		lockTime := runEveryPeriod - 500*time.Millisecond

		tryLock = func() bool {
			boolCmd := c.options.taskq.config.Redis.SetNX(ctx, key, "1", lockTime)
			return boolCmd.Val()
		}
	}

	// handler will be called by cron every runEveryPeriod seconds
	handler := func() {
		if tryLock != nil && !tryLock() {
			return
		}
		_ = c.options.taskq.queue.Add(taskMessage)
	}
	_, err := c.options.cronService.AddFunc(
		fmt.Sprintf("@every %ds", int(runEveryPeriod.Seconds())),
		handler,
	)
	return err
}
