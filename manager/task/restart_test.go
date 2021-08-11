package task

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/observiq/observiq-collector/collector"
	"github.com/stretchr/testify/require"
)

func TestExecuteRestart(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "temp")
	require.NoError(t, err)
	configPath := filepath.Join(tempDir, "config.yaml")

	testCases := []struct {
		name            string
		task            *Task
		collector       *collector.Collector
		permissions     fs.FileMode
		existingConfig  string
		expectedMessage string
		expectedStatus  Status
		expectedRunning bool
	}{
		{
			name: "invalid type",
			task: &Task{
				Type: "invalid",
			},
			collector:       collector.New(configPath, nil),
			expectedMessage: "task is not a restart",
			expectedStatus:  Failure,
			expectedRunning: false,
		},
		{
			name: "restart failure",
			task: &Task{
				Type: Restart,
			},
			collector:       collector.New(configPath, nil),
			existingConfig:  "invalid config",
			expectedMessage: "failed to restart",
			expectedStatus:  Failure,
			expectedRunning: false,
		},
		{
			name: "valid restart",
			task: &Task{
				Type: Restart,
			},
			collector:       collector.New(configPath, nil),
			existingConfig:  validConfig,
			expectedStatus:  Success,
			expectedRunning: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ioutil.WriteFile(configPath, []byte(tc.existingConfig), 0777)
			require.NoError(t, err)
			defer os.Remove(configPath)
			defer tc.collector.Stop()

			response := ExecuteRestart(tc.task, tc.collector)
			require.Equal(t, tc.expectedStatus, response.Status)
			require.Equal(t, tc.expectedMessage, response.Message)
			require.Equal(t, tc.expectedRunning, tc.collector.Status().Running)
		})
	}
}
