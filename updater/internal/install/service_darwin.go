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

// NewService returns an instance of the Service interface for managing the observiq-otel-collector service on the current OS.
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
	//#nosec G204 -- serviceName is not determined by user input
	cmd := exec.Command("launchctl", "start", d.serviceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running launchctl failed: %w", err)
	}
	return nil
}

// Stop the service
func (d darwinService) Stop() error {
	//#nosec G204 -- serviceName is not determined by user input
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

	//#nosec G204 -- installedServiceFilePath is not determined by user input
	cmd := exec.Command("launchctl", "load", d.installedServiceFilePath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running launchctl failed: %w", err)
	}

	return nil
}

// Uninstalls the service
func (d darwinService) Uninstall() error {
	//#nosec G204 -- installedServiceFilePath is not determined by user input
	cmd := exec.Command("launchctl", "unload", d.installedServiceFilePath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running launchctl failed: %w", err)
	}

	if err := os.Remove(d.installedServiceFilePath); err != nil {
		return fmt.Errorf("failed to remove service file: %w", err)
	}

	return nil
}

// InstallDir returns the filepath to the install directory
//revive:disable-next-line:exported it stutters but is an apt name
func InstallDir() (string, error) {
	return "/opt/observiq-otel-collector", nil
}
