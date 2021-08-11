package task

import "fmt"

func ExecuteUnknown(task *Task) Response {
	err := fmt.Errorf("unsupported type: %s", task.Type)
	return task.Failure("unsupported type", err)
}
