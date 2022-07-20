package file

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopyFile(t *testing.T) {
	t.Run("Copies file when output does not exist", func(t *testing.T) {
		tmpDir := t.TempDir()

		inFile := filepath.Join("testdata", "test.txt")
		outFile := filepath.Join(tmpDir, "test.txt")

		err := CopyFile(inFile, outFile)
		require.NoError(t, err)
		require.FileExists(t, outFile)

		contentsIn, err := os.ReadFile(inFile)
		require.NoError(t, err)

		contentsOut, err := os.ReadFile(outFile)
		require.NoError(t, err)

		require.Equal(t, contentsIn, contentsOut)

		fi, err := os.Stat(outFile)
		require.NoError(t, err)
		// file mode on windows acts unlike unix, we'll only check for this on linux/darwin
		if runtime.GOOS != "windows" {
			require.Equal(t, fs.FileMode(0600), fi.Mode())
		}
	})

	t.Run("Copies file when output already exists", func(t *testing.T) {
		tmpDir := t.TempDir()

		inFile := filepath.Join("testdata", "test.txt")
		outFile := filepath.Join(tmpDir, "test.txt")

		contentsIn, err := os.ReadFile(inFile)
		require.NoError(t, err)

		err = os.WriteFile(outFile, []byte("This is a file that already exists"), 0640)
		require.NoError(t, err)

		err = CopyFile(inFile, outFile)
		require.NoError(t, err)
		require.FileExists(t, outFile)

		contentsOut, err := os.ReadFile(outFile)
		require.NoError(t, err)
		require.Equal(t, contentsIn, contentsOut)

		fi, err := os.Stat(outFile)
		require.NoError(t, err)
		// file mode on windows acts unlike unix, we'll only check for this on linux/darwin
		if runtime.GOOS != "windows" {
			require.Equal(t, fs.FileMode(0640), fi.Mode())
		}
	})

	t.Run("Fails when input file does not exist", func(t *testing.T) {
		tmpDir := t.TempDir()

		inFile := filepath.Join("testdata", "does-not-exist.txt")
		outFile := filepath.Join(tmpDir, "test.txt")

		err := CopyFile(inFile, outFile)
		require.ErrorContains(t, err, "failed to open input file")
		require.NoFileExists(t, outFile)
	})

	t.Run("Does not truncate if input file does not exist", func(t *testing.T) {
		tmpDir := t.TempDir()

		inFile := filepath.Join("testdata", "does-not-exist.txt")
		outFile := filepath.Join(tmpDir, "test.txt")

		err := os.WriteFile(outFile, []byte("This is a file that already exists"), 0600)
		require.NoError(t, err)

		err = CopyFile(inFile, outFile)
		require.ErrorContains(t, err, "failed to open input file")
		require.FileExists(t, outFile)

		contentsOut, err := os.ReadFile(outFile)
		require.NoError(t, err)
		require.Equal(t, []byte("This is a file that already exists"), contentsOut)
	})
}
