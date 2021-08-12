package task

import (
	"errors"
)

// Shutdown is a task that schedules a shutdown.
const Shutdown Type = "shutdownAgent"

// ExecuteShutdown executes a shutdown task.
func ExecuteShutdown(task *Task, exit chan int) Response {
	if task.Type != Shutdown {
		err := errors.New("invalid type")
		return task.Failure("task is not a shutdown", err)
	}

	exit <- 216
	return task.Success()
}
