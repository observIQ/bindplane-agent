// Copyright  observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logging

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func TestNewLoggerConfig(t *testing.T) {
	t.Setenv("MYVAR", "/some/path")

	cases := []struct {
		name       string
		configPath string
		expect     *LoggerConfig
	}{
		{
			"file config",
			"testdata/info.yaml",
			&LoggerConfig{
				Output: fileOutput,
				Level:  zapcore.InfoLevel,
				File: &lumberjack.Logger{
					Filename:   "log/collector.log",
					MaxBackups: 5,
					MaxSize:    1,
					MaxAge:     7,
				},
			},
		},
		{
			"stdout config",
			"testdata/stdout.yaml",
			&LoggerConfig{
				Output: stdOutput,
				Level:  zapcore.DebugLevel,
			},
		},
		{
			"config with environment variables in filename",
			"testdata/expand-env.yaml",
			&LoggerConfig{
				Output: fileOutput,
				Level:  zapcore.InfoLevel,
				File: &lumberjack.Logger{
					Filename:   "/some/path/collector.log",
					MaxBackups: 5,
					MaxSize:    1,
					MaxAge:     7,
				},
			},
		},
		{
			"config does not exist",
			"testdata/does-not-exist.yaml",
			&LoggerConfig{
				Output: stdOutput,
				Level:  zapcore.InfoLevel,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			conf, err := NewLoggerConfig(tc.configPath)
			require.NoError(t, err)
			require.Equal(t, tc.expect, conf)

			opts, err := conf.Options()
			require.NoError(t, err)
			require.NotNil(t, opts)
			require.Len(t, opts, 1)
		})
	}
}

func TestNewLoggerConfigNotSpecified(t *testing.T) {
	t.Run("config does not exist in default location", func(t *testing.T) {
		tempDir := t.TempDir()
		chDir(t, tempDir)

		require.NoFileExists(t, defaultConfigPath)

		conf, err := NewLoggerConfig("")
		require.NoError(t, err)
		require.Equal(t, defaultConfig(), conf)

		require.FileExists(t, defaultConfigPath)

		// Calling again with the existing config should give the same result
		conf, err = NewLoggerConfig("")
		require.NoError(t, err)
		require.Equal(t, defaultConfig(), conf)
	})

	t.Run("config exists in the default location", func(t *testing.T) {
		tempDir := t.TempDir()

		testYaml, err := filepath.Abs(filepath.Join("testdata", "info.yaml"))
		require.NoError(t, err)

		testYamlBytes, err := os.ReadFile(testYaml)
		require.NoError(t, err)

		chDir(t, tempDir)

		err = os.WriteFile(defaultConfigPath, testYamlBytes, 0600)
		require.NoError(t, err)

		conf, err := NewLoggerConfig("")
		require.NoError(t, err)
		require.Equal(t, &LoggerConfig{
			Output: fileOutput,
			Level:  zapcore.InfoLevel,
			File: &lumberjack.Logger{
				Filename:   "log/collector.log",
				MaxBackups: 5,
				MaxSize:    1,
				MaxAge:     7,
			},
		}, conf)
	})

}

func chDir(t *testing.T, dir string) {
	t.Helper()

	oldWd, err := os.Getwd()
	require.NoError(t, err)

	err = os.Chdir(dir)
	require.NoError(t, err)

	t.Cleanup(func() {
		err = os.Chdir(oldWd)
		require.NoError(t, err)
	})
}
