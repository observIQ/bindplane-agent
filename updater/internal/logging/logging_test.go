package logging

import (
	"os"
	"path/filepath"
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
	tmpDir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(tmpDir, "log"), 0775))

	logger, err := NewLogger(tmpDir, zapcore.DebugLevel)
	require.NoError(t, err)

	logger.Info("This is a log message")
	require.NoError(t, logger.Sync())

	require.FileExists(t, filepath.Join(tmpDir, "log", "updater.log"))
}
