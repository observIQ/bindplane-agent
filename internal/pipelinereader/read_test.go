package pipelinereader

import (
	"path/filepath"
	"testing"

	"github.com/observiq/observiq-collector/collector"
	"github.com/observiq/observiq-collector/receiver/logsreceiver"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestRead(t *testing.T) {
	validConfigPath := filepath.Join(".", "testdata", "valid.yaml")
	invalidConfigPath := filepath.Join(".", "testdata", "invalid.yaml")

	testCases := []struct {
		name          string
		configPath    string
		isValid       bool
		pipelineTypes []string
	}{
		{
			name:          "valid config",
			configPath:    validConfigPath,
			isValid:       true,
			pipelineTypes: []string{"noop", "kubernetes_container", "cabin_output"},
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
			require.Equal(t, len(pipeline), len(tc.pipelineTypes))
			validatePipeline(t, tc.pipelineTypes, pipeline)
		} else {
			require.Error(t, err)
		}
	}
}

func setupCollector(configPath string) *collector.Collector {
	col := collector.New(configPath, []zap.Option{})
	return col
}

func validatePipeline(t *testing.T, pipelineTypes []string, pipeline logsreceiver.OperatorConfigs) {
	for idx1, pt := range pipelineTypes {
		found := false
		sameOrder := false
		for idx2, o := range pipeline {
			if o["type"] == pt {
				found = true
			}
			if idx1 == idx2 {
				sameOrder = true
			}
		}
		require.True(t, found && sameOrder)
	}
}
