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

package report

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	// errNoActiveReporter is the error returned when getting a specific reporter and it isn't initialized
	errNoActiveReporter = errors.New("no currently active reporter")
)

// Manager represents a structure that manages all of the different reporters
type Manager struct {
	client    Client
	reporters map[string]Reporter
	mutex     sync.Mutex
}

// SetClient sets the client of the manager to the passed in client
func (m *Manager) SetClient(client Client) error {
	if client == nil {
		return errors.New("client must not be nil")
	}

	m.client = client

	return nil
}

// ResetConfig resets the current config
func (m *Manager) ResetConfig(configData []byte) error {
	// Create a basic map so we can unmarshal on reporter specific configs
	cfg := make(map[string]any)

	if err := yaml.Unmarshal(configData, cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Lock so we don't have multiple configs being processed at concurrently
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Iterate through all reporter configs and marshal
	for kind, rawCfg := range cfg {
		switch kind {
		case snapShotKind:
			var ssCfg snapshotConfig
			if err := unmarshalReporterConfig(rawCfg, &ssCfg); err != nil {
				return fmt.Errorf("failed to unmarshal Snapshot config: %w", err)
			}

			// Verify we have a snapshot reporter initialized
			reporter, ok := m.reporters[kind]
			if !ok {
				reporter = NewSnapshotReporter(m.client)
				m.reporters[kind] = reporter
			}

			// Reconfigure reporter
			if err := m.reconfigureReporter(reporter, &ssCfg); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unrecognized reporter kind %s", kind)
		}
	}

	return nil
}

func (m *Manager) reconfigureReporter(reporter Reporter, cfg any) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Stop the reporter
	if err := reporter.Stop(ctx); err != nil {
		return fmt.Errorf("reporter %s failed to stop: %w", reporter.Kind(), err)
	}

	// Apply the new config
	if err := reporter.Report(cfg); err != nil {
		return fmt.Errorf("reporter %s failed to report with new config: %w", reporter.Kind(), err)
	}

	return nil
}

// Shutdown shuts down and cleans up all managed reporters
func (m *Manager) Shutdown(ctx context.Context) error {
	for _, reporter := range m.reporters {
		if err := reporter.Stop(ctx); err != nil {
			return fmt.Errorf("failed to shutdown reporter %s: %w", reporter.Kind(), err)
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
		manager = &Manager{
			client:    http.DefaultClient,
			reporters: make(map[string]Reporter),
		}
	})

	return manager
}

// GetSnapshotReporter returns the
func GetSnapshotReporter() *SnapshotReporter {
	currentManager := GetManager()

	// Look if we have a snapshot reporter if not create one
	currentManager.mutex.Lock()
	reporter, ok := currentManager.reporters[snapShotKind]
	if !ok {
		// Create new snapshot reporter
		reporter = NewSnapshotReporter(currentManager.client)
		currentManager.reporters[snapShotKind] = reporter
	}
	currentManager.mutex.Unlock()

	// should be safe as we only set the reporter in this function
	snapshotReporter := reporter.(*SnapshotReporter)

	return snapshotReporter
}

// unmarshalReporterConfig unmarshals a raw yaml interface into a reporter specific config structure
func unmarshalReporterConfig(inCfg, outCfg any) error {
	data, err := yaml.Marshal(inCfg)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, outCfg)
}
