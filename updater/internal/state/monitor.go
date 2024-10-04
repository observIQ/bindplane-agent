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

// Package state contains structures to monitor and update the state of the collector in the package status
package state

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/observiq/bindplane-agent/packagestate"
	"github.com/open-telemetry/opamp-go/protobufs"
	"go.uber.org/zap"
)

var (
	// ErrFailedStatus is the error when the Package status indicates a failure
	ErrFailedStatus = errors.New("package status indicates failure")
)

// Monitor allows checking and setting state of active install
//
//go:generate mockery --name Monitor --filename mock_monitor.go --structname MockMonitor
type Monitor interface {
	// SetState sets the state for the package.
	// If passed in statusErr is not nil it will record the error as the message
	SetState(packageName string, status protobufs.PackageStatusEnum, statusErr error) error

	// MonitorForSuccess will periodically check the state of the package. It will keep checking until the context is canceled or a failed/success state is detected.
	// It will return an error if status is Failed or if the context times out.
	MonitorForSuccess(ctx context.Context, hcePort int) error
}

// CollectorMonitor implements Monitor interface for monitoring the Collector Package Status file
type CollectorMonitor struct {
	stateManager  packagestate.StateManager
	currentStatus *protobufs.PackageStatuses
	logger        *zap.Logger
}

// NewCollectorMonitor create a new Monitor specifically for the collector
func NewCollectorMonitor(logger *zap.Logger, installDir string) (Monitor, error) {
	namedLogger := logger.Named("collector-monitor")

	// Create a collector monitor
	packageStatusPath := filepath.Join(installDir, packagestate.DefaultFileName)
	collectorMonitor := &CollectorMonitor{
		stateManager: packagestate.NewFileStateManager(namedLogger, packageStatusPath),
		logger:       namedLogger,
	}

	// Load the current status to ensure the package status file exists
	var err error
	collectorMonitor.currentStatus, err = collectorMonitor.stateManager.LoadStatuses()
	if err != nil {
		return nil, fmt.Errorf("failed to load package statues: %w", err)
	}

	return collectorMonitor, nil

}

// SetState sets the status on the specified package and saves it to the package status file
func (c *CollectorMonitor) SetState(packageName string, status protobufs.PackageStatusEnum, statusErr error) error {
	// Verify we have package by that name
	targetPackage, ok := c.currentStatus.GetPackages()[packageName]
	if !ok {
		return fmt.Errorf("no package for name %s", packageName)
	}

	// Update the status
	targetPackage.Status = status

	// If that passed in error is not nil set it as the error message
	if statusErr != nil {
		targetPackage.ErrorMessage = statusErr.Error()
	}

	c.currentStatus.GetPackages()[packageName] = targetPackage

	// Save to updated status to disk
	return c.stateManager.SaveStatuses(c.currentStatus)
}

// MonitorForSuccess checks the collector health check extension to verify if it is healthy
// If an InstallFailed status is read this returns ErrFailedStatus error.
// If the context is canceled the context error will be returned.
// Uses a retry loop with 3 max tries and 3 second delay between each.
func (c *CollectorMonitor) MonitorForSuccess(ctx context.Context, hcePort int) error {
	endpoint := fmt.Sprintf("http://127.0.0.1:%d", hcePort)
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return err
	}
	client := http.Client{
		Timeout: time.Second * 10,
	}

	var resp *http.Response
	for r := 0; r < 3; r++ {
		resp, err = client.Do(req)
		if err == nil {
			break
		}

		c.logger.Info("request failed, retrying...", zap.Error(err))
		time.Sleep(time.Second * 3)
	}
	if err != nil {
		return fmt.Errorf("failed to reach agent after 3 attempts: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check on %s returned %d", hcePort, resp.StatusCode)
	}

	return nil
}
