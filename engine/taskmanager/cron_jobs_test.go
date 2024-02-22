package taskmanager

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCronTasks(t *testing.T) {
	makeTaskManager := func() *TaskManager {
		client, _ := NewTaskManager(
			context.Background(),
		)
		return client.(*TaskManager)
	}

	t.Run("register one cron job", func(t *testing.T) {
		tm := makeTaskManager()

		desiredExecutions := 2

		type mockTarget struct {
			times chan bool
		}
		target := &mockTarget{make(chan bool, desiredExecutions)}

		err := tm.CronJobsInit(CronJobs{
			"test": {
				Period: 1 * time.Second,
				Handler: func(ctx context.Context) error {
					target.times <- true
					return nil
				},
			},
		})

		require.NoError(t, err)

		// wait for the task to run 'desiredExecutions' times
		executions := 0
		for i := 0; i < desiredExecutions; i++ {
			<-target.times
			executions++
		}

		// check number of executions
		require.Equal(t, desiredExecutions, executions)
	})

	t.Run("register two cron jobs", func(t *testing.T) {
		tm := makeTaskManager()

		desiredExecutions := 6

		type mockTarget struct {
			times chan int
		}
		target := &mockTarget{make(chan int, desiredExecutions)}

		err := tm.CronJobsInit(CronJobs{
			"test1": {
				Period: 1 * time.Second,
				Handler: func(ctx context.Context) error {
					target.times <- 1
					return nil
				},
			},
			"test2": {
				Period: 1 * time.Second,
				Handler: func(ctx context.Context) error {
					target.times <- 2
					return nil
				},
			},
		})

		require.NoError(t, err)

		executions1 := 0
		executions2 := 0
		for i := 0; i < desiredExecutions; i++ {
			jobID := <-target.times
			if jobID == 1 {
				executions1++
			} else {
				executions2++
			}
		}

		// check number of executions, both jobs should run at least once
		require.Greater(t, executions1, 0)
		require.Greater(t, executions2, 0)
	})
}
