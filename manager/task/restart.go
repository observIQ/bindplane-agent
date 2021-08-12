package task

import (
	"errors"

	"github.com/observiq/observiq-collector/collector"
)

// Restart is the task type assigned to restarting the collector.
// Previously, this task was used to restart the entire manager.
// For the sake of backwards compatibility, the key remains the same.
const Restart Type = "restartManager"

// ExecuteRestart executes a restart task.
func ExecuteRestart(task *Task, collector *collector.Collector) Response {
	if task.Type != Restart {
		err := errors.New("invalid type")
		return task.Failure("task is not a restart", err)
	}

	err := collector.Restart()
	if err != nil {
		return task.Failure("failed to restart", err)
	}

	return task.Success()
}
