package config

import (
	"github.com/BuxOrg/bux/taskmanager"
)

const (
	TaskManagerEngine    = taskmanager.TaskQ
	TaskManagerQueueName = "bux_queue"
)
