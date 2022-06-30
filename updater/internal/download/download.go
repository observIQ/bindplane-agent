package download

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	archiver "github.com/mholt/archiver/v3"
)

const extractFolder = "latest"

// Downloads the file into the outPath, truncating the file if it already exists
func downloadFile(downloadUrl string, outPath string) error {
	resp, err := http.Get(downloadUrl)
	if err != nil {
		return fmt.Errorf("could not GET url: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("got non-200 status code (%d)", resp.StatusCode)
	}

	f, err := os.OpenFile(outPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0640)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	if _, err = io.Copy(f, resp.Body); err != nil {
		return fmt.Errorf("failed to copy request body to file: %w", err)
	}

	return nil
}

// getOutputFilePath gets the output path relative to the base dir for the archive from the given URL.
func getOutputFilePath(basePath, downloadUrl string) (string, error) {
	url, err := url.Parse(downloadUrl)
	if err != nil {
		return "", fmt.Errorf("cannot parse url: %w", err)
	}

	if url.Path == "" {
		return "", errors.New("input url must have path")
	}

	return filepath.Join(basePath, filepath.Base(url.Path)), nil
}

func verifyContentHash(contentPath, hexExpectedContentHash string) error {
	expectedContentHash, err := hex.DecodeString(hexExpectedContentHash)
	if err != nil {
		return fmt.Errorf("failed to decode content hash: %w", err)
	}

	// Hash file at contentPath using sha256
	fileHash := sha256.New()

	f, err := os.Open(contentPath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	if _, err = io.Copy(fileHash, f); err != nil {
		return fmt.Errorf("failed to calculate file hash: %w", err)
	}

	actualContentHash := fileHash.Sum(nil)
	if subtle.ConstantTimeCompare(expectedContentHash, actualContentHash) == 0 {
		return errors.New("content hashes were not equal")
	}

	return nil
}

func DownloadAndVerify(url, dir, expectedHash string) error {
	archiveFilePath, err := getOutputFilePath(dir, url)
	if err != nil {
		return fmt.Errorf("failed to determine archive download path: %w", err)
	}

	if err := downloadFile(url, archiveFilePath); err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}

	extractPath := filepath.Join(dir, extractFolder)

	if err := verifyContentHash(archiveFilePath, expectedHash); err != nil {
		return fmt.Errorf("content hash could not be verified: %w", err)
	}

	// Clean the "latest" dir before extraction
	if err := os.RemoveAll(extractPath); err != nil {
		return fmt.Errorf("error cleaning archive extraction target path: %w", err)
	}

	if err := archiver.Unarchive(archiveFilePath, extractPath); err != nil {
		return fmt.Errorf("failed to extract file: %w", err)
	}

	return nil
}
