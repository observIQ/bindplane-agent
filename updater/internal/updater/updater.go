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

// Package updater handles all aspects of updating the collector from a provided archive
package updater

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/observiq/bindplane-agent/packagestate"
	"github.com/observiq/bindplane-agent/updater/internal/action"
	"github.com/observiq/bindplane-agent/updater/internal/install"
	"github.com/observiq/bindplane-agent/updater/internal/path"
	"github.com/observiq/bindplane-agent/updater/internal/rollback"
	"github.com/observiq/bindplane-agent/updater/internal/service"
	"github.com/observiq/bindplane-agent/updater/internal/state"
	"github.com/open-telemetry/opamp-go/protobufs"
	"go.uber.org/zap"
)

// Updater is a struct that can be used to perform a collector update
type Updater struct {
	installDir string
	installer  install.Installer
	svc        service.Service
	rollbacker rollback.Rollbacker
	monitor    state.Monitor
	logger     *zap.Logger
}

// NewUpdater creates a new updater which can be used to update the installation based at
// installDir
func NewUpdater(logger *zap.Logger, installDir string) (*Updater, error) {
	hcePort, err := findRandomPort()
	if err != nil {
		logger.Error("failed to get random port for collector health check extension, continuing with port 12345", zap.Error(err))
		hcePort = 12345
	}

	monitor, err := state.NewCollectorMonitor(logger, installDir, hcePort)
	if err != nil {
		return nil, fmt.Errorf("failed to create monitor: %w", err)
	}

	svc := service.NewService(logger, installDir)
	return &Updater{
		installDir: installDir,
		installer:  install.NewInstaller(logger, installDir, hcePort, svc),
		svc:        svc,
		rollbacker: rollback.NewRollbacker(logger, installDir),
		monitor:    monitor,
		logger:     logger,
	}, nil
}

// Update performs the update of the collector binary
func (u *Updater) Update() error {
	// Stop the service before backing up the install directory;
	// We want to stop as early as possible so that we don't hit the collector's timeout
	// while it waits to be shutdown.
	if err := u.svc.Stop(); err != nil {
		return fmt.Errorf("failed to stop service: %w", err)
	}
	// Record that we stopped the service
	u.rollbacker.AppendAction(action.NewServiceStopAction(u.svc))

	// Now that we stopped the service, it will be our responsibility to cleanup the tmp dir.
	// We will do this regardless of whether we succeed or fail after this point.
	defer u.removeTmpDir()

	u.logger.Debug("Stopped the service")

	// Create the backup
	if err := u.rollbacker.Backup(); err != nil {
		u.logger.Error("Failed to backup", zap.Error(err))

		// Set the state to failed before rollback so collector knows it failed
		if setErr := u.monitor.SetState(packagestate.CollectorPackageName, protobufs.PackageStatusEnum_PackageStatusEnum_InstallFailed, err); setErr != nil {
			u.logger.Error("Failed to set state on backup failure", zap.Error(setErr))
		}

		u.rollbacker.Rollback()

		u.logger.Error("Rollback complete")
		return fmt.Errorf("failed to backup: %w", err)
	}

	// Install artifacts
	if err := u.installer.Install(u.rollbacker); err != nil {
		u.logger.Error("Failed to install", zap.Error(err))

		// Set the state to failed before rollback so collector knows it failed
		if setErr := u.monitor.SetState(packagestate.CollectorPackageName, protobufs.PackageStatusEnum_PackageStatusEnum_InstallFailed, err); setErr != nil {
			u.logger.Error("Failed to set state on install failure", zap.Error(setErr))
		}

		u.rollbacker.Rollback()

		u.logger.Error("Rollback complete")
		return fmt.Errorf("failed to install: %w", err)
	}

	// Create a context with timeout to wait for a success or failed status
	checkCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	u.logger.Debug("Installation successful, begin monitor for success")

	// Monitor the install state
	if err := u.monitor.MonitorForSuccess(checkCtx); err != nil {
		u.logger.Error("Failed to install", zap.Error(err))

		// Set the state to failed before rollback so collector knows it failed
		if setErr := u.monitor.SetState(packagestate.CollectorPackageName, protobufs.PackageStatusEnum_PackageStatusEnum_InstallFailed, err); setErr != nil {
			u.logger.Error("Failed to set state on install failure", zap.Error(setErr))
		}
		u.rollbacker.Rollback()

		u.logger.Error("Rollback complete")
		return fmt.Errorf("failed while monitoring for success: %w", err)
	}

	// Remove excess files & folders
	u.cleanup()

	// Successful update
	u.logger.Info("Update Complete")
	return nil
}

// removeTmpDir removes the temporary directory that holds the update artifacts.
func (u *Updater) removeTmpDir() {
	err := os.RemoveAll(path.TempDir(u.installDir))
	if err != nil {
		u.logger.Error("failed to remove temporary directory", zap.Error(err))
	}
}

// cleanup removes unnecessary files now that agent has updated from v1 to v2.
// Update is already successfully completed at this point, so we don't want to return errors.
func (u *Updater) cleanup() {
	files := []string{filepath.Join("log", "collector.log"), "config.yaml", "manager.yaml", "updater", "config.bak.yaml", "logging.yaml", "package_statuses.json"}
	for _, f := range files {
		err := os.Remove(filepath.Join(u.installDir, f))
		if err != nil {
			u.logger.Info("failed to cleanup a directory entry", zap.String("file", f), zap.Error(err))
		}
	}
}

func findRandomPort() (int, error) {
	l, err := net.Listen("tcp", "localhost:0")

	if err != nil {
		return 0, err
	}

	port := l.Addr().(*net.TCPAddr).Port

	err = l.Close()

	if err != nil {
		return 0, err
	}

	return port, nil
}
