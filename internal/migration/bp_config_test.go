package migration

import (
	"fmt"
	"path"
	"testing"

	"github.com/observiq/observiq-collector/internal/logging"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var loggingConfigDir = path.Join(".", "testdata", "bpagentlogging")

func TestConvertBPLogConfig(t *testing.T) {
	defaultConfig := logging.DefaultConfig()

	testCases := []struct {
		name         string
		fileIn       string
		readError    string
		convertError string
		out          *logging.Config
	}{
		{
			name:   "Default bpagent config",
			fileIn: path.Join(loggingConfigDir, "test_log.yaml"),
			out: &logging.Config{
				Collector: logging.LoggerConfig{
					Level:        zap.InfoLevel,
					MaxBackups:   5,
					MaxMegabytes: 1,
					MaxDays:      7,
					File:         defaultConfig.Collector.File,
				},
				Manager: logging.LoggerConfig{
					Level:        zap.InfoLevel,
					MaxBackups:   5,
					MaxMegabytes: 1,
					MaxDays:      7,
					File:         defaultConfig.Manager.File,
				},
			},
		},
		{
			name:      "File does not exist",
			fileIn:    path.Join(loggingConfigDir, "does_not_exist.yaml"),
			readError: fmt.Sprintf("open %s: no such file or directory", path.Join(loggingConfigDir, "does_not_exist.yaml")),
		},
		{
			name:   "Does not fail with extra keys",
			fileIn: path.Join(loggingConfigDir, "extra_keys.yaml"),
			out: &logging.Config{
				Collector: logging.LoggerConfig{
					Level:        zap.InfoLevel,
					MaxBackups:   5,
					MaxMegabytes: 1,
					MaxDays:      7,
					File:         defaultConfig.Collector.File,
				},
				Manager: logging.LoggerConfig{
					Level:        zap.InfoLevel,
					MaxBackups:   5,
					MaxMegabytes: 1,
					MaxDays:      7,
					File:         defaultConfig.Manager.File,
				},
			},
		},
		{
			name:      "Fails with non-yaml file",
			fileIn:    path.Join(loggingConfigDir, "not_yaml.txt"),
			readError: "yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `This is...` into migration.bpLogConfig",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			bpConf, err := LoadBPLogConfig(testCase.fileIn)
			if testCase.readError == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Equal(t, testCase.readError, err.Error())
				return
			}

			out, err := BPLogConfigToLogConfig(*bpConf)
			if testCase.convertError == "" {
				require.NoError(t, err)
				require.Equal(t, testCase.out, out)
			} else {
				require.Error(t, err)
				require.Equal(t, testCase.convertError, err.Error())
			}
		})
	}
}

func TestBPLogConfigToLogConfig(t *testing.T) {
	defaultConfig := logging.DefaultConfig()

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
