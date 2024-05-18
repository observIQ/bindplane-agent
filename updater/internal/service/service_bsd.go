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

//go:build freebsd || openbsd || netbsd

package service

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/observiq/bindplane-agent/updater/internal/file"
	"github.com/observiq/bindplane-agent/updater/internal/path"
	"go.uber.org/zap"
)

const bsdUnixServiceName = "observiq-otel-collector"
const bsdUnixServiceFilePath = "/usr/lib/systemd/system/observiq-otel-collector.service"

// Option is an extra option for creating a Service
type Option func(bsdUnixSvc *bsdUnixService)

// WithServiceFile returns an option setting the service file to use when updating using the service
func WithServiceFile(svcFilePath string) Option {
	return func(bsdUnixSvc *bsdUnixService) {
		bsdUnixSvc.newServiceFilePath = svcFilePath
	}
}

// NewService returns an instance of the Service interface for managing the observiq-otel-collector service on the current OS.
func NewService(logger *zap.Logger, installDir string, opts ...Option) Service {
	bsdUnixSvc := &bsdUnixService{
		newServiceFilePath:       filepath.Join(path.ServiceFileDir(installDir), "observiq-otel-collector.rc"),
		serviceName:              bsdUnixServiceName,
		installedServiceFilePath: bsdUnixServiceFilePath,
		installDir:               installDir,
		logger:                   logger.Named("bsdUnix-service"),
	}

	for _, opt := range opts {
		opt(bsdUnixSvc)
	}

	return bsdUnixSvc
}

type bsdUnixService struct {
	// newServiceFilePath is the file path to the new unit file
	newServiceFilePath string
	// serviceName is the name of the service
	serviceName string
	// installedServiceFilePath is the file path to the installed unit file
	installedServiceFilePath string
	installDir               string
	logger                   *zap.Logger
}

// Start the service
func (l bsdUnixService) Start() error {
	//#nosec G204 -- serviceName is not determined by user input
	cmd := exec.Command("service", l.serviceName, "start")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running service failed: %w", err)
	}
	return nil
}

// Stop the service
func (l bsdUnixService) Stop() error {
	//#nosec G204 -- serviceName is not determined by user input
	cmd := exec.Command("service", l.serviceName, "stop")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running service failed: %w", err)
	}
	return nil
}

// installs the service
func (l bsdUnixService) install() error {
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

	//#nosec G204 -- serviceName is not determined by user input
	cmd := exec.Command("service", l.serviceName, "enable")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("enabling service file failed: %w", err)
	}

	return nil
}

// uninstalls the service
func (l bsdUnixService) uninstall() error {
	//#nosec G204 -- serviceName is not determined by user input
	cmd := exec.Command("service", l.serviceName, "disable")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to disable service: %w", err)
	}

	if err := os.Remove(l.installedServiceFilePath); err != nil {
		return fmt.Errorf("failed to remove service file: %w", err)
	}

	cmd = exec.Command("service", "daemon-reload")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("reloading service failed: %w", err)
	}

	return nil
}

func (l bsdUnixService) Update() error {
	if err := l.uninstall(); err != nil {
		return fmt.Errorf("failed to uninstall old service: %w", err)
	}

	if err := l.install(); err != nil {
		return fmt.Errorf("failed to install new service: %w", err)
	}

	return nil
}

func (l bsdUnixService) Backup() error {
	if err := file.CopyFileNoOverwrite(l.logger.Named("copy-file"), l.installedServiceFilePath, path.BackupServiceFile(l.installDir)); err != nil {
		return fmt.Errorf("failed to copy service file: %w", err)
	}

	return nil
}
