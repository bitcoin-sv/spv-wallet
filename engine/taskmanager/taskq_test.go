package taskmanager

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/stretchr/testify/require"
)

// NOTE: because of the taskq package has global state, the names of tasks must be unique
func TestNewTaskManager_Single(t *testing.T) {
	ctx := context.Background()
	c, err := NewTaskManager(ctx)
	require.NoError(t, err)
	require.NotNil(t, c)

	task1Chan := make(chan string, 1)
	task1Arg := "task a"

	err = c.RegisterTask(task1Arg, func(name string) error {
		task1Chan <- name
		return nil
	})
	require.NoError(t, err)

	require.Equal(t, 1, len(c.Tasks()))

	// Run single task
	err = c.RunTask(ctx, &TaskRunOptions{
		Arguments: []interface{}{task1Arg},
		TaskName:  task1Arg,
	})
	require.NoError(t, err)

	require.Equal(t, task1Arg, <-task1Chan)

	// Close the client
	err = c.Close(context.Background())
	require.NoError(t, err)
}

func TestNewTaskManager_Multiple(t *testing.T) {
	ctx := context.Background()
	c, _ := NewTaskManager(ctx)

	task1Chan := make(chan string, 1)
	task2Chan := make(chan string, 1)

	task1Arg := "task b"
	task2Arg := "task c"

	err := c.RegisterTask(task1Arg, func(name string) error {
		task1Chan <- name
		return nil
	})
	require.NoError(t, err)

	err = c.RegisterTask(task2Arg, func(name string) error {
		task2Chan <- name
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, 2, len(c.Tasks()))

	// Run tasks
	err = c.RunTask(ctx, &TaskRunOptions{
		Arguments: []interface{}{task1Arg},
		TaskName:  task1Arg,
	})
	require.NoError(t, err)

	err = c.RunTask(ctx, &TaskRunOptions{
		Arguments: []interface{}{task2Arg},
		TaskName:  task2Arg,
	})
	require.NoError(t, err)

	require.Equal(t, task1Arg, <-task1Chan)
	require.Equal(t, task2Arg, <-task2Chan)

	// Close the client
	err = c.Close(context.Background())
	require.NoError(t, err)
}

func TestNewTaskManager_RegisterTwice(t *testing.T) {
	ctx := context.Background()
	c, _ := NewTaskManager(ctx)

	task1Arg := "task d"
	resultChan := make(chan int, 1)

	err := c.RegisterTask(task1Arg, func(_ string) error {
		resultChan <- 1
		return nil
	})
	require.NoError(t, err)

	err = c.RegisterTask(task1Arg, func(_ string) error {
		resultChan <- 2
		return nil
	})
	require.NoError(t, err)

	err = c.RunTask(ctx, &TaskRunOptions{
		Arguments: []interface{}{task1Arg},
		TaskName:  task1Arg,
	})
	require.NoError(t, err)

	require.Equal(t, 1, <-resultChan)
}

func TestNewTaskManager_RunTwice(t *testing.T) {
	ctx := context.Background()
	c, _ := NewTaskManager(ctx)

	task1Arg := "task e"

	err := c.RegisterTask(task1Arg, func(_ string) error {
		return nil
	})
	require.NoError(t, err)

	err = c.RunTask(ctx, &TaskRunOptions{
		TaskName: task1Arg,
	})
	require.NoError(t, err)

	err = c.RunTask(ctx, &TaskRunOptions{
		TaskName: task1Arg,
	})
	require.NoError(t, err)
}

func TestNewTaskManager_NotRegistered(t *testing.T) {
	ctx := context.Background()
	c, _ := NewTaskManager(ctx)

	task1Arg := "task f"

	err := c.RunTask(ctx, &TaskRunOptions{
		TaskName: task1Arg,
	})
	require.Error(t, err)
}

func TestNewTaskManager_WithRedis(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live local redis tests")
	}

	queueName, _ := utils.RandomHex(8)
	ctx := context.Background()
	c, err := NewTaskManager(ctx, WithTaskqConfig(DefaultTaskQConfig(queueName, WithRedis("redis://localhost:6379"))))
	require.NoError(t, err)
	require.NotNil(t, c)

	task1Chan := make(chan string, 1)
	task1Arg := "task redis"

	err = c.RegisterTask(task1Arg, func(name string) error {
		task1Chan <- name
		return nil
	})
	require.NoError(t, err)

	require.Equal(t, 1, len(c.Tasks()))

	// Run single task
	err = c.RunTask(ctx, &TaskRunOptions{
		Arguments: []interface{}{task1Arg},
		TaskName:  task1Arg,
	})
	require.NoError(t, err)

	require.Equal(t, task1Arg, <-task1Chan)

	// Close the client
	err = c.Close(context.Background())
	require.NoError(t, err)
}
