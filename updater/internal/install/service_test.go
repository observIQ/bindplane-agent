package install

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReplaceInstallDir(t *testing.T) {
	testCases := []struct {
		input      []byte
		installDir string
		output     []byte
	}{
		{
			input:      []byte("[INSTALLDIR]"),
			installDir: "some/install/directory",
			output:     []byte(filepath.Join("some", "install", "directory") + string(os.PathSeparator)),
		},
		{
			input:      []byte("no install dir"),
			installDir: "some/install/directory",
			output:     []byte("no install dir"),
		},
		{
			input:      []byte("[INSTALLDIR]observiq-otel-collector"),
			installDir: "some/install/directory",
			output:     []byte(filepath.Join("some", "install", "directory", "observiq-otel-collector")),
		},
	}

	for _, tc := range testCases {
		t.Run(string(tc.input), func(t *testing.T) {
			out := replaceInstallDir(tc.input, tc.installDir)
			require.Equal(t, tc.output, out)
		})
	}
}
