package report

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var (
	// errNoActiveReporter is the error returned when getting a specific reporter and it isn't initialized
	errNoActiveReporter = errors.New("no currently active reporter")
)

// Manager represents a structure that manages all of the different reporters
type Manager struct {
	client    Client
	reporters map[ReporterType]Reporter
}

// ResetConfig resets the current config
func (m *Manager) ResetConfig() error {
	return nil
}

// Shutdown shuts down and cleans up all managed reporters
func (m *Manager) Shutdown(ctx context.Context) error {
	for _, reporter := range m.reporters {
		if err := reporter.Stop(ctx); err != nil {
			return fmt.Errorf("failed to shutdown reporter %s: %w", reporter.Type(), err)
		}
	}

	return nil
}

// variables to manage singleton
var (
	managerOnce sync.Once
	manager     *Manager
)

// GetManager returns the global Manager
func GetManager() *Manager {
	managerOnce.Do(func() {
		manager = &Manager{}
	})

	return manager
}

// GetSnapshotReporter returns the
func GetSnapshotReporter() (*SnapshotReporter, error) {
	reporter, ok := manager.reporters[snapShotType]
	if !ok {
		return nil, errNoActiveReporter
	}

	// Do a cast to expected typ to double check
	snapshotReporter, ok := reporter.(*SnapshotReporter)
	if !ok {
		// This should not ever happen but making this as a failsafe
		return nil, errors.New("invalid reporter for type")
	}

	return snapshotReporter, nil
}
