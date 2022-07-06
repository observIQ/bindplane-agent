package install

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func InstallArtifacts(latestDirPath, installDirPath string, svc Service) error {
	// Stop service
	err := svc.Stop()
	if err != nil {
		return fmt.Errorf("failed to stop service: %w", err)
	}

	// install files that go to installDirPath to their correct location,
	// excluding any config files (logging.yaml, config.yaml, manager.yaml)
	if err := installFiles(latestDirPath, installDirPath); err != nil {
		return fmt.Errorf("failed to install new files: %w", err)
	}

	// Uninstall previous service
	if err := svc.Uninstall(); err != nil {
		return fmt.Errorf("failed to uninstall service: %w", err)
	}

	// Install new service
	if err := svc.Install(); err != nil {
		return fmt.Errorf("failed to install service: %w", err)
	}

	// Start service
	if err := svc.Start(); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	return nil
}

func installFiles(latestDirPath, installDirPath string) error {
	err := filepath.WalkDir(latestDirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if skipFile(path) {
			return nil
		}

		relPath, err := filepath.Rel(latestDirPath, path)
		if err != nil {
			return err
		}

		outPath := filepath.Join(installDirPath, relPath)
		outDir := filepath.Dir(outPath)

		if err := os.MkdirAll(outDir, 0750); err != nil {
			return fmt.Errorf("failed to create dir: %w", err)
		}

		outFile, err := os.OpenFile(outPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
		if err != nil {
			return fmt.Errorf("failed to open output file: %w", err)
		}
		defer func() {
			err := outFile.Close()
			if err != nil {
				log.Default().Printf("installFiles: Failed to close output file: %s", err)
			}
		}()

		inFile, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open input file: %w", err)
		}
		defer func() {
			err := inFile.Close()
			if err != nil {
				log.Default().Printf("installFiles: Failed to close input file: %s", err)
			}
		}()

		if _, err := io.Copy(outFile, inFile); err != nil {
			return fmt.Errorf("failed to copy file: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk latest dir: %w", err)
	}

	return nil
}

// skipFile returns true if the given path is a special config file.
// These files should not be overwritten.
func skipFile(path string) bool {
	var configFiles = []string{
		"config.yaml",
		"logging.yaml",
		"manager.yaml",
	}

	fileName := filepath.Base(path)

	for _, f := range configFiles {
		if fileName == f {
			return true
		}
	}

	return false
}
