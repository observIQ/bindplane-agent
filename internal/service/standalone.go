package service

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/observiq/observiq-otel-collector/collector"
)

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
	go s.monitorStatus()
	return nil
}

// monitorStatus monitors the collector's status for errors, and reports them
// to the error channel to trigger a shutdown.
func (s StandaloneCollectorService) monitorStatus() {
	defer s.wg.Done()
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
