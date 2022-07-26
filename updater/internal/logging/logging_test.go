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
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestLevelFromString(t *testing.T) {
	testCases := []struct {
		name        string
		in          string
		out         zapcore.Level
		expectedErr string
	}{
		{
			name: "debug level",
			in:   "debug",
			out:  zapcore.DebugLevel,
		},
		{
			name: "info level",
			in:   "info",
			out:  zapcore.InfoLevel,
		},
		{
			name: "error level",
			in:   "error",
			out:  zapcore.ErrorLevel,
		},
		{
			name:        "unrecognized level",
			in:          "unrecognized",
			expectedErr: "unrecognized level",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := LevelFromString(tc.in)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.Equal(t, tc.out, out)
			}
		})
	}
}

func TestNewLogger(t *testing.T) {
	if runtime.GOOS == "windows" {
		// Skip on windows, because the log file will still be open
		// when the test attempts to remove the temp dir, which ends up making
		// the test fail.
		t.SkipNow()
	}
	tmpDir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(tmpDir, "log"), 0775))

	logger, err := NewLogger(tmpDir, zapcore.DebugLevel)
	require.NoError(t, err)

	logger.Info("This is a log message")
	require.NoError(t, logger.Sync())

	require.FileExists(t, filepath.Join(tmpDir, "log", "updater.log"))
}
