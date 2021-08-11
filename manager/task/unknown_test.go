package task

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExecuteUnknown(t *testing.T) {
	task := &Task{
		Type: "test-type",
		ID:   "test-id",
	}
	response := ExecuteUnknown(task)
	require.Equal(t, Failure, response.Status)
	require.Equal(t, "unsupported type", response.Message)
}
