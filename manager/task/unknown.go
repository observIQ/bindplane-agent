package task

func ExecuteUnknown(task Task) Response {
	return Response{
		Type:    task.Type,
		ID:      task.ID,
		Status:  Exception,
		Message: "unknown task type",
	}
}
