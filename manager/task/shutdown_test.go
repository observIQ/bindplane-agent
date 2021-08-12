package task

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExecuteShutdown(t *testing.T) {
	testCases := []struct {
		name            string
		task            *Task
		expectedMessage string
		expectedStatus  Status
		expectedExit    bool
	}{
		{
			name: "invalid type",
			task: &Task{
				Type: "invalid",
			},
			expectedMessage: "task is not a shutdown",
			expectedStatus:  Failure,
		},
		{
			name: "with success",
			task: &Task{
				Type: Shutdown,
			},
			expectedStatus: Success,
			expectedExit:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			exit := make(chan int, 1)
			response := ExecuteShutdown(tc.task, exit)
			require.Equal(t, tc.expectedStatus, response.Status)
			require.Equal(t, tc.expectedMessage, response.Message)

			if tc.expectedExit {
				exitCode := <-exit
				require.Equal(t, 216, exitCode)
			}
		})
	}
}
