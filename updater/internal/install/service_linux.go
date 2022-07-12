//go:build linux

package install

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const linuxServiceName = "observiq-otel-collector"
const linuxServiceFilePath = "/usr/lib/systemd/system/observiq-otel-collector.service"

func NewService(latestPath string) Service {
	return linuxService{
		newServiceFilePath:       filepath.Join(latestPath, "install", "observiq-otel-collector.service"),
		serviceName:              linuxServiceName,
		installedServiceFilePath: linuxServiceFilePath,
	}
}

type linuxService struct {
	// newServiceFilePath is the file path to the new unit file
	newServiceFilePath string
	// serviceName is the name of the service
	serviceName string
	// installedServiceFilePath is the file path to the installed unit file
	installedServiceFilePath string
}

// Start the service
func (l linuxService) Start() error {
	cmd := exec.Command("systemctl", "start", l.serviceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running systemctl failed: %w", err)
	}
	return nil
}

// Stop the service
func (l linuxService) Stop() error {
	cmd := exec.Command("systemctl", "stop", l.serviceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running systemctl failed: %w", err)
	}
	return nil
}

// Installs the service
func (l linuxService) Install() error {
	inFile, err := os.Open(l.newServiceFilePath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer func() {
		err := inFile.Close()
		if err != nil {
			log.Default().Printf("Service Install: Failed to close input file: %s", err)
		}
	}()

	outFile, err := os.OpenFile(l.installedServiceFilePath, os.O_CREATE|os.O_WRONLY, 0600)
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

	cmd := exec.Command("systemctl", "daemon-reload")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("reloading systemctl failed: %w", err)
	}

	cmd = exec.Command("systemctl", "enable", l.serviceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("enabling unit file failed: %w", err)
	}

	return nil
}

// Uninstalls the service
func (l linuxService) Uninstall() error {
	cmd := exec.Command("systemctl", "disable", l.serviceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to disable unit: %w", err)
	}

	if err := os.Remove(l.installedServiceFilePath); err != nil {
		return fmt.Errorf("failed to remove service file: %w", err)
	}

	cmd = exec.Command("systemctl", "daemon-reload")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("reloading systemctl failed: %w", err)
	}

	return nil
}

func InstallDir() (string, error) {
	return "/opt/observiq-otel-collector", nil
}
