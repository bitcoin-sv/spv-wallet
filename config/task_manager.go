package config

import (
	"github.com/BuxOrg/bux/taskmanager"
)

// TaskManager defaults
const (
	TaskManagerEngine    = taskmanager.TaskQ
	TaskManagerQueueName = "bux_queue"
)
