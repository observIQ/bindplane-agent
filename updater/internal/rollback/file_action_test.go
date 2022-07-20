package rollback

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCopyFileAction(t *testing.T) {
	t.Run("out file does not exist", func(t *testing.T) {
		scratchDir := t.TempDir()
		testTempDir := filepath.Join("testdata", "copyfileaction")
		outFile := filepath.Join(scratchDir, "test.txt")
		inFile := filepath.Join(testTempDir, "latest", "test.txt")

		a, err := NewCopyFileAction(inFile, outFile, testTempDir)
		require.NoError(t, err)

		require.Equal(t, &CopyFileAction{
			FromPath:    inFile,
			ToPath:      outFile,
			FileCreated: true,
			rollbackDir: filepath.Join(testTempDir, "rollback"),
			latestDir:   filepath.Join(testTempDir, "latest"),
		}, a)
	})

	t.Run("out file exists", func(t *testing.T) {
		scratchDir := t.TempDir()
		testTempDir := filepath.Join("testdata", "copyfileaction")
		outFile := filepath.Join(scratchDir, "test.txt")
		inFile := filepath.Join(testTempDir, "latest", "test.txt")

		f, err := os.Create(outFile)
		require.NoError(t, err)
		require.NoError(t, f.Close())

		a, err := NewCopyFileAction(inFile, outFile, testTempDir)
		require.NoError(t, err)

		require.Equal(t, &CopyFileAction{
			FromPath:    inFile,
			ToPath:      outFile,
			FileCreated: false,
			rollbackDir: filepath.Join(testTempDir, "rollback"),
			latestDir:   filepath.Join(testTempDir, "latest"),
		}, a)
	})
}

func TestCopyFileActionRollback(t *testing.T) {
	t.Run("deletes out file if it does not exist", func(t *testing.T) {
		scratchDir := t.TempDir()
		testTempDir := filepath.Join("testdata", "copyfileaction")
		outFile := filepath.Join(scratchDir, "test.txt")
		inFile := filepath.Join(testTempDir, "latest", "test.txt")

		a, err := NewCopyFileAction(inFile, outFile, testTempDir)
		require.NoError(t, err)

		inBytes, err := os.ReadFile(inFile)
		require.NoError(t, err)

		err = os.WriteFile(outFile, inBytes, 0600)
		require.NoError(t, err)

		err = a.Rollback()
		require.NoError(t, err)

		require.NoFileExists(t, outFile)
	})

	t.Run("Rolls back out file when it exists", func(t *testing.T) {
		scratchDir := t.TempDir()
		testTempDir := filepath.Join("testdata", "copyfileaction")
		outFile := filepath.Join(scratchDir, "test.txt")
		inFile := filepath.Join(testTempDir, "latest", "test.txt")
		originalFile := filepath.Join(testTempDir, "rollback", "test.txt")

		originalBytes, err := os.ReadFile(originalFile)
		require.NoError(t, err)

		err = os.WriteFile(outFile, originalBytes, 0600)
		require.NoError(t, err)

		a, err := NewCopyFileAction(inFile, outFile, testTempDir)
		require.NoError(t, err)

		// Overwrite original file with latest file
		inBytes, err := os.ReadFile(inFile)
		require.NoError(t, err)

		err = os.WriteFile(outFile, inBytes, 0600)
		require.NoError(t, err)

		err = a.Rollback()
		require.NoError(t, err)

		require.FileExists(t, outFile)

		rolledBackBytes, err := os.ReadFile(outFile)
		require.NoError(t, err)

		require.Equal(t, originalBytes, rolledBackBytes)
	})

	t.Run("Fails if backup file doesn't exist", func(t *testing.T) {
		scratchDir := t.TempDir()
		testTempDir := filepath.Join("testdata", "copyfileaction")
		outFile := filepath.Join(scratchDir, "test.txt")
		inFile := filepath.Join(testTempDir, "latest", "not_in_backup.txt")
		originalFile := filepath.Join(testTempDir, "rollback", "test.txt")

		// The latest file exists in the directory already, but for some reason is not copied to backup
		originalBytes, err := os.ReadFile(originalFile)
		require.NoError(t, err)

		err = os.WriteFile(outFile, originalBytes, 0600)
		require.NoError(t, err)

		a, err := NewCopyFileAction(inFile, outFile, testTempDir)
		require.NoError(t, err)

		// Overwrite original file with latest file
		latestBytes, err := os.ReadFile(inFile)
		require.NoError(t, err)

		err = os.WriteFile(outFile, latestBytes, 0600)
		require.NoError(t, err)

		err = a.Rollback()
		require.ErrorContains(t, err, "failed to copy file")
		require.FileExists(t, outFile)

		finalBytes, err := os.ReadFile(outFile)
		require.NoError(t, err)
		require.Equal(t, latestBytes, finalBytes)
	})

}
