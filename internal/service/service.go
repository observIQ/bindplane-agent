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

package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/observiq/observiq-otel-collector/collector"
	"go.uber.org/zap"
)

const (
	startTimeout = 10 * time.Second
	stopTimeout  = 10 * time.Second
)

// RunnableService may be run as a service.
type RunnableService interface {
	// Start asynchronously starts the underlying service. The service may not necessarily be "ready"
	// once this returns, but could be asynchronously starting up.
	Start(ctx context.Context) error
	// Stop synchronously shuts down the service. After this function returns, the underlying service should be completely stopped.
	Stop(ctx context.Context) error
	// Error returns an error channel that should emit an error when the service must unexpectedly quit.
	Error() <-chan error
}

// StandaloneCollectorService is a RunnableService that runs the collector in standalone mode.
type StandaloneCollectorService struct {
	col      collector.Collector
	doneChan chan struct{}
	errChan  chan error
	wg       *sync.WaitGroup
}

// NewStandaloneCollectorService creates a new StandaloneCollectorService
func NewStandaloneCollectorService(c collector.Collector) StandaloneCollectorService {
	return StandaloneCollectorService{
		col:      c,
		doneChan: make(chan struct{}, 1),
		errChan:  make(chan error, 1),
		wg:       &sync.WaitGroup{},
	}
}

// Start starts the collector
func (s StandaloneCollectorService) Start(ctx context.Context) error {
	collectorStartedChan := make(chan struct{})
	go func() {
		// Start's context is only valid for the lifetime of "Start",
		// but the collector expects a context that is valid for the lifetime of the service.
		err := s.col.Run(context.Background())
		if err != nil {
			s.errChan <- err
		}
		close(collectorStartedChan)
	}()

	select {
	case <-collectorStartedChan: // OK
	case <-ctx.Done():
		return fmt.Errorf("failed while waiting for service startup: %w", ctx.Err())
	}

	// monitor status for errors, so we don't zombie the service
	s.wg.Add(1)
	go s.monitorStatus(s.wg)
	return nil
}

// monitorStatus monitors the collector's status for errors, and reports them
// to the error channel to trigger a shutdown.
func (s StandaloneCollectorService) monitorStatus(wg *sync.WaitGroup) {
	defer wg.Done()
	statusChan := s.col.Status()
	for {
		select {
		case status := <-statusChan:
			if status.Err != nil {
				s.errChan <- status.Err
			} else if !status.Running {
				// If we aren't running, bail out. Otherwise the collector is effectively a "zombie" process.
				s.errChan <- errors.New("collector unexpectedly stopped running")
			}
		case <-s.doneChan:
			return
		}
	}
}

// Error returns a channel that can emit asynchronous, unrecoverable errors
func (s StandaloneCollectorService) Error() <-chan error {
	return s.errChan
}

// Stop shuts down the underlying collector
func (s StandaloneCollectorService) Stop(ctx context.Context) error {
	close(s.doneChan)

	collectorStoppedChan := make(chan struct{})
	go func() {
		s.col.Stop()
		s.wg.Wait()
		close(collectorStoppedChan)
	}()

	select {
	case <-collectorStoppedChan:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("failed while waiting for service shutdown: %w", ctx.Err())
	}
}

// runServiceInteractive runs the service in an "interactive" mode (responds to SIGINT and SIGTERM).
// This mode is always used in linux, and is used in Windows when the collector
// is not running as a service.
func runServiceInteractive(ctx context.Context, logger *zap.Logger, svc RunnableService) error {
	startupTimeoutCtx, startupCancel := context.WithTimeout(ctx, startTimeout)
	defer startupCancel()

	if err := svc.Start(startupTimeoutCtx); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	var svcErr error
	// Service is started; Wait for a stop signal.
	select {
	case <-ctx.Done():
	case svcErr = <-svc.Error():
		logger.Error("Unexpected error while running service", zap.Error(svcErr))
	}

	stopTimeoutCtx, stopCancel := context.WithTimeout(context.Background(), stopTimeout)
	defer stopCancel()

	if err := svc.Stop(stopTimeoutCtx); err != nil {
		return fmt.Errorf("failed to stop service: %w", err)
	}

	return svcErr
}
