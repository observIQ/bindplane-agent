package task

import (
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/observiq/observiq-collector/collector"
	"github.com/stretchr/testify/require"
)

func TestReconfigureParams(t *testing.T) {
	noopOperator := map[string]interface{}{"type": "noop"}
	cabinOperator := map[string]interface{}{"type": "cabin_output"}

	emptyParams := ReconfigureParams{
		Config: StanzaConfig{
			Pipeline: StanzaPipeline{},
		},
	}

	multipleParams := ReconfigureParams{
		Config: StanzaConfig{
			Pipeline: StanzaPipeline{cabinOperator, noopOperator},
		},
	}

	emptyPipeline := emptyParams.getStanzaPipeline()
	require.Equal(t, 1, len(emptyPipeline))
	require.Equal(t, noopOperator, emptyPipeline[0])

	multiplePipeline := multipleParams.getStanzaPipeline()
	require.Equal(t, 2, len(multiplePipeline))
	require.Equal(t, cabinOperator, multiplePipeline[1])
}

func TestReconfigure(t *testing.T) {
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
			name: "invalid task type",
			task: &Task{
				Type: "invalid",
			},
			collector:       collector.New(configPath, nil),
			permissions:     0777,
			expectedStatus:  Failure,
			expectedMessage: "task is not a reconfigure",
		},
		{
			name: "invalid task parameters",
			task: &Task{
				Type: Reconfigure,
				Parameters: map[string]interface{}{
					"config": 1,
				},
			},
			collector:       collector.New(configPath, nil),
			permissions:     0777,
			expectedStatus:  Failure,
			expectedMessage: "unable to decode parameters",
		},
		{
			name: "missing collector config",
			task: &Task{
				Type:       Reconfigure,
				Parameters: createPipeline("noop"),
			},
			collector:       collector.New("invalid.yaml", nil),
			permissions:     0777,
			expectedStatus:  Failure,
			expectedMessage: "failed to read existing config",
		},
		{
			name: "invalid existing config",
			task: &Task{
				Type:       Reconfigure,
				Parameters: createPipeline("noop"),
			},
			collector:       collector.New(configPath, nil),
			permissions:     0777,
			existingConfig:  "invalid yaml contents",
			expectedStatus:  Failure,
			expectedMessage: "failed to decode existing config",
		},
		{
			name: "invalid new config",
			task: &Task{
				Type:       Reconfigure,
				Parameters: createPipeline(&badYAML{}),
			},
			collector:       collector.New(configPath, nil),
			permissions:     0777,
			existingConfig:  validConfig,
			expectedStatus:  Failure,
			expectedMessage: "failed to convert new config to yaml",
		},
		{
			name: "missing write permissions",
			task: &Task{
				Type:       Reconfigure,
				Parameters: createPipeline("noop"),
			},
			collector:       collector.New(configPath, nil),
			permissions:     0444,
			existingConfig:  validConfig,
			expectedStatus:  Failure,
			expectedMessage: "failed to write new config",
		},
		{
			name: "config validation failure",
			task: &Task{
				Type:       Reconfigure,
				Parameters: createPipeline("noop"),
			},
			collector:       collector.New(configPath, nil),
			permissions:     0777,
			existingConfig:  "receivers: null",
			expectedStatus:  Failure,
			expectedMessage: "new config failed validation",
		},
		{
			name: "failed restart",
			task: &Task{
				Type:       Reconfigure,
				Parameters: createPipeline(nil),
			},
			collector:       collector.New(configPath, nil),
			permissions:     0777,
			existingConfig:  validConfig,
			expectedStatus:  Failure,
			expectedMessage: "failed to restart collector",
		},
		{
			name: "success with noop",
			task: &Task{
				Type:       Reconfigure,
				Parameters: createPipeline("noop"),
			},
			collector:       collector.New(configPath, nil),
			permissions:     0777,
			existingConfig:  validConfig,
			expectedRunning: true,
			expectedStatus:  Success,
			expectedMessage: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ioutil.WriteFile(configPath, []byte(tc.existingConfig), tc.permissions)
			require.NoError(t, err)
			defer os.Remove(configPath)
			defer tc.collector.Stop()

			response := ExecuteReconfigure(tc.task, tc.collector)
			require.Equal(t, tc.expectedStatus, response.Status)
			require.Equal(t, tc.expectedMessage, response.Message)
			require.Equal(t, tc.expectedRunning, tc.collector.Status().Running)
		})
	}
}

// createPipeline creates a stanza pipeline for testing
func createPipeline(operatorType interface{}) map[string]interface{} {
	operator := map[string]interface{}{"type": operatorType}
	operators := []map[string]interface{}{operator}
	return map[string]interface{}{
		"config": map[string]interface{}{
			"pipeline": operators,
		},
	}
}

// badYAML is a struct used to mock yaml marshalling errors
type badYAML struct{}

// MarshalYAML implements the yaml Marshaler interface
func (b *badYAML) MarshalYAML() (interface{}, error) {
	return nil, errors.New("bad yaml")
}

var validConfig = `
receivers:
  stanza:
    pipeline:
      - type: noop
  
exporters:
  logging:
  
service:
  pipelines:
    logs:
      receivers: [stanza]
      exporters: [logging]
`
