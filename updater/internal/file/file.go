package file

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func CopyFile(pathIn, pathOut string) error {
	pathInClean := filepath.Clean(pathIn)
	pathOutClean := filepath.Clean(pathOut)

	// Open the output file, creating it if it does not exist and truncating it.
	outFile, err := os.OpenFile(pathOutClean, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}
	defer func() {
		err := outFile.Close()
		if err != nil {
			log.Default().Printf("CopyFile: Failed to close output file: %s", err)
		}
	}()

	// Open the input file for reading.
	inFile, err := os.Open(pathInClean)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer func() {
		err := inFile.Close()
		if err != nil {
			log.Default().Printf("CopyFile: Failed to close input file: %s", err)
		}
	}()

	// Copy the input file to the output file.
	if _, err := io.Copy(outFile, inFile); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}
