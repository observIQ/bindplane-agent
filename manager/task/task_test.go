package task

import (
	"errors"
	"testing"

	"github.com/observiq/observiq-collector/manager/message"
	"github.com/stretchr/testify/require"
)

func TestTaskSuccess(t *testing.T) {
	task := Task{
		Type: "unknown",
		ID:   "test-id",
	}

	response := task.Success()
	require.Equal(t, task.Type, response.Type)
	require.Equal(t, task.ID, response.ID)
	require.Equal(t, Success, response.Status)
}

func TestTaskFailure(t *testing.T) {
	task := Task{
		Type: "unknown",
		ID:   "test-id",
	}
	err := errors.New("unknown failure")

	response := task.Failure("task failed", err)
	require.Equal(t, task.Type, response.Type)
	require.Equal(t, task.ID, response.ID)
	require.Equal(t, Failure, response.Status)
	require.Equal(t, "task failed", response.Message)
	require.Equal(t, err.Error(), response.Details["Error"])
}

func TestResponseToMessage(t *testing.T) {
	response := &Response{}
	msg := response.ToMessage()
	require.Equal(t, message.TaskResponse, msg.Type)
}

func TestFromMessageSuccess(t *testing.T) {
	expected := &Task{
		Type:       "test-task",
		ID:         "test-id",
		Parameters: map[string]interface{}{},
	}

	msg, err := message.New(message.TaskRequest, expected)
	require.NoError(t, err)

	task, err := FromMessage(msg)
	require.NoError(t, err)
	require.Equal(t, expected, task)
}

func TestFromMessageInvalidType(t *testing.T) {
	sampleTask := &Task{
		Type:       "test-task",
		ID:         "test-id",
		Parameters: map[string]interface{}{},
	}

	msg, err := message.New(message.TaskResponse, sampleTask)
	require.NoError(t, err)

	task, err := FromMessage(msg)
	require.Error(t, err)
	require.Nil(t, task)
	require.Contains(t, err.Error(), "invalid message type")
}

func TestFromMessageInvalidContent(t *testing.T) {
	msg := &message.Message{
		Type: message.TaskRequest,
		Content: map[string]interface{}{
			"id": make(chan int),
		},
	}

	task, err := FromMessage(msg)
	require.Error(t, err)
	require.Nil(t, task)
	require.Contains(t, err.Error(), "failed to decode task")
}
