//go:build darwin && !linux && !windows

package install

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const darwinServiceName = "com.observiq.collector"
const darwinServiceFilePath = "/Library/LaunchDaemons/com.observiq.collector.plist"

func NewService(latestPath string) Service {
	return &darwinService{
		newServiceFilePath:       filepath.Join(latestPath, "install", "com.observiq.collector.plist"),
		serviceName:              darwinServiceName,
		installedServiceFilePath: darwinServiceFilePath,
	}
}

type darwinService struct {
	// newServiceFilePath is the file path to the new plist file
	newServiceFilePath string
	// serviceName is the name of the service
	serviceName string
	// installedServiceFilePath is the file path to the installed plist file
	installedServiceFilePath string
}

// Start the service
func (d darwinService) Start() error {
	cmd := exec.Command("launchctl", "start", d.serviceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running launchctl failed: %w", err)
	}
	return nil
}

// Stop the service
func (d darwinService) Stop() error {
	cmd := exec.Command("launchctl", "stop", d.serviceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running launchctl failed: %w", err)
	}
	return nil
}

// Installs the service
func (d darwinService) Install() error {
	inFile, err := os.Open(d.newServiceFilePath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer func() {
		err := inFile.Close()
		if err != nil {
			log.Default().Printf("Service Install: Failed to close input file: %s", err)
		}
	}()

	outFile, err := os.OpenFile(d.installedServiceFilePath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}
	defer func() {
		err := outFile.Close()
		if err != nil {
			log.Default().Printf("Service Install: Failed to close output file: %s", err)
		}
	}()

	if _, err := io.Copy(outFile, inFile); err != nil {
		return fmt.Errorf("failed to copy service file: %w", err)
	}

	cmd := exec.Command("launchctl", "load", d.installedServiceFilePath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running launchctl failed: %w", err)
	}

	return nil
}

// Uninstalls the service
func (d darwinService) Uninstall() error {
	cmd := exec.Command("launchctl", "unload", d.installedServiceFilePath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running launchctl failed: %w", err)
	}

	if err := os.Remove(d.installedServiceFilePath); err != nil {
		return fmt.Errorf("failed to remove service file: %w", err)
	}

	return nil
}
