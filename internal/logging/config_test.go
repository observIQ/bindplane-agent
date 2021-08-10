package logging

import (
	"path"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLoadConfig(t *testing.T) {
	defaultConfig := DefaultConfig()
	testCases := []struct {
		name   string
		path   string
		config *Config
		errStr string
	}{
		{
			name: "valid_config.yaml",
			path: path.Join(".", "testdata", "valid_config.yaml"),
			config: &Config{
				Collector: LoggerConfig{
					Level:        zap.DebugLevel,
					MaxBackups:   1,
					MaxMegabytes: 1,
					MaxDays:      1,
					File:         "./local/collector.log",
				},
				Manager: LoggerConfig{
					Level:        zap.InfoLevel,
					MaxBackups:   5,
					MaxMegabytes: 6,
					MaxDays:      7,
					File:         "./local/manager.log",
				},
			},
		},
		{
			name:   "invalid_log_level.yaml",
			path:   path.Join(".", "testdata", "invalid_log_level.yaml"),
			errStr: "unrecognized level: \"inof\"",
		},
		{
			name:   "invalid_collector_mb.yaml",
			path:   path.Join(".", "testdata", "invalid_collector_mb.yaml"),
			errStr: "error validating Collector logger config: LoggerConfig.MaxMegabytes must be >= 0",
		},
		{
			name:   "empty_config.yaml",
			path:   path.Join(".", "testdata", "empty_config.yaml"),
			config: defaultConfig,
		},
		{
			name: "partial_config.yaml",
			path: path.Join(".", "testdata", "partial_config.yaml"),
			config: &Config{
				Collector: LoggerConfig{
					Level:        zap.InfoLevel,
					MaxBackups:   2,
					MaxMegabytes: 2,
					MaxDays:      2,
					File:         defaultConfig.Collector.File,
				},
				Manager: defaultConfig.Manager,
			},
		},
		{
			name:   "non-existant file",
			path:   path.Join(".", "testdata", "does_not_exist.yaml"),
			errStr: "open testdata/does_not_exist.yaml: no such file or directory",
		},
		{
			name:   "empty path",
			path:   "",
			errStr: "open : no such file or directory",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c, err := LoadConfig(testCase.path)
			if testCase.errStr == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Equal(t, testCase.errStr, err.Error())
			}

			require.Equal(t, testCase.config, c)
		})
	}
}

func TestValidateConfig(t *testing.T) {
	testCases := []struct {
		name   string
		input  *Config
		errStr string
	}{
		{
			name:  "Default config is valid",
			input: DefaultConfig(),
		},
		{
			name: "Config with collector logging max size == 0 is invalid",
			input: &Config{
				Collector: LoggerConfig{
					Level:        zapcore.ErrorLevel,
					MaxBackups:   1,
					MaxMegabytes: -1,
					MaxDays:      1,
					File:         "somefile.log",
				},
				Manager: LoggerConfig{
					Level:        zapcore.ErrorLevel,
					MaxBackups:   1,
					MaxMegabytes: 1,
					MaxDays:      1,
					File:         "somefile2.log",
				},
			},
			errStr: "error validating Collector logger config: LoggerConfig.MaxMegabytes must be >= 0",
		},
		{
			name: "Config with manager max backups < 0 is invalid",
			input: &Config{
				Collector: LoggerConfig{
					Level:        zapcore.ErrorLevel,
					MaxBackups:   1,
					MaxMegabytes: 1,
					MaxDays:      1,
					File:         "somefile.log",
				},
				Manager: LoggerConfig{
					Level:        zapcore.ErrorLevel,
					MaxBackups:   -1,
					MaxMegabytes: 1,
					MaxDays:      1,
					File:         "somefile2.log",
				},
			},
			errStr: "error validating Manager logger config: LoggerConfig.MaxBackups must be >= 0",
		},
		{
			name: "Config with collector logging max days <= 0 is invalid",
			input: &Config{
				Collector: LoggerConfig{
					Level:        zapcore.ErrorLevel,
					MaxBackups:   1,
					MaxMegabytes: 1,
					MaxDays:      1,
					File:         "somefile.log",
				},
				Manager: LoggerConfig{
					Level:        zapcore.ErrorLevel,
					MaxBackups:   1,
					MaxMegabytes: 1,
					MaxDays:      -1,
					File:         "somefile2.log",
				},
			},
			errStr: "error validating Manager logger config: LoggerConfig.MaxDays must be >= 0",
		},
		{
			name: "Config with File == \"\" is invalid",
			input: &Config{
				Collector: LoggerConfig{
					Level:        zapcore.ErrorLevel,
					MaxBackups:   1,
					MaxMegabytes: 1,
					MaxDays:      1,
					File:         "",
				},
				Manager: LoggerConfig{
					Level:        zapcore.ErrorLevel,
					MaxBackups:   1,
					MaxMegabytes: 1,
					MaxDays:      1,
					File:         "somefile2.log",
				},
			},
			errStr: "error validating Collector logger config: LoggerConfig.File must not be empty",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := ValidateConfig(testCase.input)
			if testCase.errStr == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Equal(t, testCase.errStr, err.Error())
			}
		})
	}
}
