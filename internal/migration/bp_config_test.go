package migration

import (
	"testing"

	"github.com/observiq/observiq-collector/internal/logging"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestBPLogConfigToLogConfig(t *testing.T) {
	defaultConfig, err := logging.DefaultConfig()
	require.NoError(t, err)

	testCases := []struct {
		name      string
		input     *bpLogConfig
		output    *logging.Config
		shouldErr bool
	}{
		{
			name:      "empty bp logging.yaml gives collector default",
			input:     &bpLogConfig{},
			output:    defaultConfig,
			shouldErr: false,
		},
		{
			name: "Full logging.yaml",
			input: &bpLogConfig{
				MaxBackups:   intAsPtr(9),
				MaxMegabytes: intAsPtr(10),
				MaxDays:      intAsPtr(11),
				Level:        levelAsPtr(zapcore.ErrorLevel),
			},
			output: &logging.Config{
				Collector: logging.LoggerConfig{
					MaxBackups:   9,
					MaxMegabytes: 10,
					MaxDays:      11,
					Level:        zapcore.ErrorLevel,
					File:         defaultConfig.Collector.File,
				},
				Manager: logging.LoggerConfig{
					MaxBackups:   9,
					MaxMegabytes: 10,
					MaxDays:      11,
					Level:        zapcore.ErrorLevel,
					File:         defaultConfig.Manager.File,
				},
			},
			shouldErr: false,
		},
		{
			name: "Partial logging.yaml",
			input: &bpLogConfig{
				MaxMegabytes: intAsPtr(10),
				Level:        levelAsPtr(zapcore.ErrorLevel),
			},
			output: &logging.Config{
				Collector: logging.LoggerConfig{
					MaxBackups:   defaultConfig.Collector.MaxBackups,
					MaxMegabytes: 10,
					MaxDays:      defaultConfig.Collector.MaxDays,
					Level:        zapcore.ErrorLevel,
					File:         defaultConfig.Collector.File,
				},
				Manager: logging.LoggerConfig{
					MaxBackups:   defaultConfig.Manager.MaxBackups,
					MaxMegabytes: 10,
					MaxDays:      defaultConfig.Manager.MaxDays,
					Level:        zapcore.ErrorLevel,
					File:         defaultConfig.Manager.File,
				},
			},
			shouldErr: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			conf, err := BPLogConfigToLogConfig(*testCase.input)
			if testCase.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, testCase.output, conf)
			}
		})
	}
}

func levelAsPtr(level zapcore.Level) *zapcore.Level {
	return &level
}

func intAsPtr(val int) *int {
	return &val
}
