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

// Package observiq contains OpAmp structures compatible with the observiq client
package observiq

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/observiq/observiq-otel-collector/internal/version"
	"github.com/open-telemetry/opamp-go/client/types"
	"github.com/open-telemetry/opamp-go/protobufs"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// Ensure interface is satisfied
var _ types.PackagesStateProvider = (*packagesStateProvider)(nil)

// packagesStateProvider represents a PackagesStateProvider which uses a yaml file to persist PackageStatuses
type packagesStateProvider struct {
	yamlPath string
	logger   *zap.Logger
}

// NewPackagesStateProvider creates a new OpAmp PackagesStateProvider
func newPackagesStateProvider(logger *zap.Logger, yamlPath string) types.PackagesStateProvider {
	return &packagesStateProvider{
		yamlPath: filepath.Clean(yamlPath),
		logger:   logger,
	}
}

// AllPackagesHash not implemented so returns an error with this info
func (p *packagesStateProvider) AllPackagesHash() ([]byte, error) {
	p.logger.Debug("Retrieve all packages hash")

	return nil, errors.New("method not implemented: PackageStateProvider AllPackagesHash")
}

// SetAllPackagesHash not implemented so returns an error with this info
func (p *packagesStateProvider) SetAllPackagesHash(_ []byte) error {
	p.logger.Debug("Set all packages hash")

	return errors.New("method not implemented: PackageStateProvider SetAllPackagesHash")
}

// Packages not implemented so returns an error with this info
func (p *packagesStateProvider) Packages() ([]string, error) {
	p.logger.Debug("Retrieve package names")

	return nil, errors.New("method not implemented: PackageStateProvider Packages")
}

// PackageState not implemented so returns an error with this info
func (p *packagesStateProvider) PackageState(_ string) (state types.PackageState, err error) {
	p.logger.Debug("Retrieve package state")

	packageState := types.PackageState{}

	return packageState, errors.New("method not implemented: PackageStateProvider PackageState")
}

// SetPackageState not implemented so returns an error with this info
func (p *packagesStateProvider) SetPackageState(_ string, _ types.PackageState) error {
	p.logger.Debug("Set package state")

	return errors.New("method not implemented: PackageStateProvider SetPackageState")
}

// CreatePackage not implemented so returns an error with this info
func (p *packagesStateProvider) CreatePackage(_ string, _ protobufs.PackageAvailable_PackageType) error {
	p.logger.Debug("Create package")

	return errors.New("method not implemented: PackageStateProvider CreatePackage")
}

// FileContentHash not implemented so returns an error with this info
func (p *packagesStateProvider) FileContentHash(_ string) ([]byte, error) {
	p.logger.Debug("Retrieve package content hash")

	return nil, errors.New("method not implemented: PackageStateProvider FileContentHash")
}

// UpdateContent not implemented so returns an error with this info
func (p *packagesStateProvider) UpdateContent(_ context.Context, _ string, _ io.Reader, _ []byte) error {
	p.logger.Debug("Update package content")

	return errors.New("method not implemented: PackageStateProvider UpdateContent")
}

// DeletePackage not implemented so returns an error with this info
func (p *packagesStateProvider) DeletePackage(_ string) error {
	p.logger.Debug("Delete package")

	return errors.New("method not implemented: PackageStateProvider DeletePackage")
}

// LastReportedStatuses retrieves the PackagesStatuses from a saved yaml file
func (p *packagesStateProvider) LastReportedStatuses() (*protobufs.PackageStatuses, error) {
	p.logger.Debug("Retrieve last reported package statuses")

	packages := map[string]*protobufs.PackageStatus{
		mainPackageName: {
			Name:            mainPackageName,
			AgentHasVersion: version.Version(),
			Status:          protobufs.PackageStatus_Installed,
		},
	}
	packageStatuses := &protobufs.PackageStatuses{
		Packages: packages,
	}

	// If there's a problem reading the file, we just return a barebones PackageStatuses.
	if _, err := os.Stat(p.yamlPath); errors.Is(err, os.ErrNotExist) {
		p.logger.Debug("Package statuses yaml doesn't exist")
		return packageStatuses, nil
	}

	statusesBytes, err := os.ReadFile(p.yamlPath)
	if err != nil {
		return packageStatuses, fmt.Errorf("failed to read package statuses yaml: %w", err)
	}

	var readPackageStatuses protobufs.PackageStatuses
	if err := yaml.Unmarshal(statusesBytes, &readPackageStatuses); err != nil {
		return packageStatuses, fmt.Errorf("failed to unmarshal package statuses: %w", err)
	}

	return &readPackageStatuses, nil
}

// SetLastReportedStatuses saves the given PackageStatuses into a yaml file
func (p *packagesStateProvider) SetLastReportedStatuses(statuses *protobufs.PackageStatuses) error {
	p.logger.Debug("Set last reported package statuses")

	// If there is any problem saving the new package statuses, make sure that we delete any existing file
	// in order to not have outdated data as its better to start fresh
	if err := os.Remove(p.yamlPath); err != nil {
		p.logger.Debug("Failed to delete package statuses yaml", zap.Error(err))
	}

	data, err := yaml.Marshal(statuses)
	if err != nil {
		return fmt.Errorf("failed to marshal package statuses: %w", err)
	}

	// Write data to a package_statuses.yaml file, with 0600 file permission
	if err := os.WriteFile(p.yamlPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write package statuses yaml: %w", err)
	}

	return nil
}
