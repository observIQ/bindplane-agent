package rollback

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/observiq/observiq-otel-collector/updater/internal/path"
)

type CopyFileAction struct {
	fromPath string
	toPath   string
	// fileCreated is a bool that records whether this action had to create a new file or not
	fileCreated bool
	rollbackDir string
	latestDir   string
}

// NewCopyFileAction creates a new CopyFileAction that indicates a file was copied from
// fromPath into toPath. tmpDir is specified for rollback purposes.
// NOTE: This action MUST be created BEFORE the action actually takes place; This allows
// for previous existence of the file to be recorded.
func NewCopyFileAction(fromPath, toPath, tmpDir string) (*CopyFileAction, error) {
	fileExists := false
	_, err := os.Stat(toPath)
	switch {
	case errors.Is(err, os.ErrNotExist):
		fileExists = true
	case err != nil:
		return nil, fmt.Errorf("unexpected error stat-ing file: %w", err)
	}

	return &CopyFileAction{
		fromPath: fromPath,
		toPath:   toPath,
		// The file will be created if it doesn't already exist
		fileCreated: !fileExists,
		rollbackDir: path.BackupDirFromTempDir(tmpDir),
		latestDir:   path.LatestDirFromTempDir(tmpDir),
	}, nil
}

// Rollback will undo the file copy, by either deleting the file if the file did not originally exist,
// or it will copy the old file in the rollback dir if it already exists.
func (c CopyFileAction) Rollback() error {
	if c.fileCreated {
		// File did not exist before this action.
		// We just need to delete this file.
		return os.RemoveAll(c.toPath)
	}

	// Copy from rollback dir over the current file
	// the backup file should have the same relative path from
	// rollback dir as the fromPath does from the latest dir
	rel, err := filepath.Rel(c.latestDir, c.fromPath)
	if err != nil {
		return fmt.Errorf("could not determine relative path between latestDir (%s) and fromPath (%s): %w", c.latestDir, c.fromPath, err)
	}

	backupFilePath := filepath.Join(c.rollbackDir, rel)
	backupFile, err := os.Open(backupFilePath)
	if err != nil {
		return fmt.Errorf("failed to open to output file: %w", err)
	}
	defer func() {
		err := backupFile.Close()
		if err != nil {
			log.Default().Printf("failed to close backup file: %s", err)
		}
	}()

	toFile, err := os.Open(c.toPath)
	if err != nil {
		return fmt.Errorf("failed to open to output file: %w", err)
	}
	defer func() {
		err := toFile.Close()
		if err != nil {
			log.Default().Printf("failed to close output file: %s", err)
		}
	}()

	// Copy backup file over the in-place file in order to restore the original
	if _, err := io.Copy(toFile, backupFile); err != nil {
		return fmt.Errorf("failed to copy backup file into output file: %w", err)
	}

	return nil
}
