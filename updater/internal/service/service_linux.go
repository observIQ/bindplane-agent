// Copyright  observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build linux

package service

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/observiq/observiq-otel-collector/updater/internal/path"
)

const linuxServiceName = "observiq-otel-collector"
const linuxServiceFilePath = "/usr/lib/systemd/system/observiq-otel-collector.service"

type ServiceOption func(linuxSvc *linuxService)

func WithServiceFile(svcFilePath string) ServiceOption {
	return func(linuxSvc *linuxService) {
		linuxSvc.newServiceFilePath = svcFilePath
	}
}

// NewService returns an instance of the Service interface for managing the observiq-otel-collector service on the current OS.
func NewService(latestPath string, opts ...ServiceOption) Service {
	linuxSvc := &linuxService{
		newServiceFilePath:       filepath.Join(path.ServiceFileDir(latestPath), "observiq-otel-collector.service"),
		serviceName:              linuxServiceName,
		installedServiceFilePath: linuxServiceFilePath,
	}

	for _, opt := range opts {
		opt(linuxSvc)
	}

	return linuxSvc
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
	//#nosec G204 -- serviceName is not determined by user input
	cmd := exec.Command("systemctl", "start", l.serviceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running systemctl failed: %w", err)
	}
	return nil
}

// Stop the service
func (l linuxService) Stop() error {
	//#nosec G204 -- serviceName is not determined by user input
	cmd := exec.Command("systemctl", "stop", l.serviceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running systemctl failed: %w", err)
	}
	return nil
}

// installs the service
func (l linuxService) install() error {
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

	//#nosec G204 -- serviceName is not determined by user input
	cmd = exec.Command("systemctl", "enable", l.serviceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("enabling unit file failed: %w", err)
	}

	return nil
}

// uninstalls the service
func (l linuxService) uninstall() error {
	//#nosec G204 -- serviceName is not determined by user input
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

func (l linuxService) Update() error {
	if err := l.uninstall(); err != nil {
		return fmt.Errorf("failed to uninstall old service: %w", err)
	}

	if err := l.install(); err != nil {
		return fmt.Errorf("failed to install new service: %w", err)
	}

	return nil
}

func (l linuxService) Backup(outDir string) error {
	oldFile, err := os.Open(l.installedServiceFilePath)
	if err != nil {
		return fmt.Errorf("failed to open old service file: %w", err)
	}
	defer func() {
		err := oldFile.Close()
		if err != nil {
			log.Default().Printf("linuxService.Backup: failed to close out file: %s", err)
		}
	}()

	// Create the file in the specified location; If the file already exists, an error will be returned,
	// since we don't want to overwrite the file
	backupFile, err := os.OpenFile(path.BackupServiceFile(outDir), os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0600)
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}
	defer func() {
		err := backupFile.Close()
		if err != nil {
			log.Default().Printf("linuxService.Backup: failed to close out file: %s", err)
		}
	}()

	if _, err := io.Copy(backupFile, oldFile); err != nil {
		return fmt.Errorf("failed to copy service file to backup: %w", err)
	}

	return nil
}
