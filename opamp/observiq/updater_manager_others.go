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

//go:build !windows

package observiq

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

const updaterFolder = "latest"

// Ensure interface is satisfied
var _ UpdaterManager = (*OthersUpdaterManager)(nil)

// OthersUpdaterManager handles starting a Updater binary and watching it for failure with a timeout
type OthersUpdaterManager struct {
	tmpPath string
	logger  *zap.Logger
}

// newUpdaterManager creates a new UpdaterManager
func newUpdaterManager(defaultLogger *zap.Logger, tmpPath string) UpdaterManager {
	return &OthersUpdaterManager{
		tmpPath: filepath.Clean(tmpPath),
		logger:  defaultLogger.Named("updater manager"),
	}
}

// StartAndMonitorUpdater will start the Updater binary and wait to see if it finishes unexpectedly.
// While waiting for Updater, it should kill the collector and we should never execute any code past running it
func (m OthersUpdaterManager) StartAndMonitorUpdater() error {
	updaterPath := filepath.Join(m.tmpPath, updaterFolder, "updater")
	absTmpPath, err := filepath.Abs(m.tmpPath)
	if err != nil {
		m.logger.Warn("Failed to get absolute path of tmp dir", zap.Error(err))
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()
	cmd := exec.CommandContext(ctx, updaterPath, "--tmpdir", absTmpPath)

	// Blocks until updater finishes
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("updater exited before shutting down collector: %w", err)
	}

	// Ideally we should not get here as we will be killed by the updater.
	// Updater should either exit before us with error or we die before it does.
	return fmt.Errorf("updater exited before shutting down collector")
}
