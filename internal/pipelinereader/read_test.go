package pipelinereader

import (
	"path/filepath"
	"testing"

	"github.com/observiq/observiq-collector/collector"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestRead(t *testing.T) {
	validConfigPath := filepath.Join(".", "testdata", "valid.yaml")
	invalidConfigPath := filepath.Join(".", "testdata", "invalid.yaml")

	testCases := []struct {
		name           string
		configPath     string
		isValid        bool
		pipelineLength int
	}{
		{
			name:           "valid config",
			configPath:     validConfigPath,
			isValid:        true,
			pipelineLength: 1,
		},
		{
			name:       "invalid config",
			configPath: invalidConfigPath,
			isValid:    false,
		},
	}

	for _, tc := range testCases {
		col := setupCollector(tc.configPath)
		pipeline, err := Read(col)
		if tc.isValid {
			require.NoError(t, err)
			require.Equal(t, len(pipeline), tc.pipelineLength)
		} else {
			require.Error(t, err)
		}
	}
}

func setupCollector(configPath string) *collector.Collector {
	col := collector.New(configPath, []zap.Option{})
	return col
}
