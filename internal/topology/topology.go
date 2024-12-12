// Copyright observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package topology provides code to help manage topology updates for BindPlane and the topology processor.
package topology

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ConfigTopologyRegistry represents a registry for the topology processor to register their ConfigTopology.
type ConfigTopologyRegistry interface {
	// RegisterConfigTopology registers the topology state for the given processor.
	// It should return an error if the processor has already been registered.
	RegisterConfigTopology(processorID string, data *ConfigTopologyState) error
	SetIntervalChan() chan time.Duration
	Reset()
}

// GatewayInfo represents the unique identifiable information about a bindplane gateway's configuration
type GatewayInfo struct {
	GatewayID  string
	ConfigName string
	AccountID  string
	OrgID      string
}

// ConfigTopologyState represents the data captured through topology processors.
type ConfigTopologyState struct {
	ConfigTopology configTopology
	mux            sync.Mutex
}

type configTopology struct {
	DestGateway GatewayInfo
	RouteTable  map[GatewayInfo]time.Time
}

// ConfigTopologyInfo represents topology relationships between configs.
type ConfigTopologyInfo struct {
	GatewayID     string         `json:"gatewayID"`
	ConfigName    string         `json:"configName"`
	AccountID     string         `json:"accountID"`
	OrgID         string         `json:"orgID"`
	SourceConfigs []ConfigRecord `json:"sourceConfigs"`
}

// ConfigRecord represents a gateway source and the time it was last detected
type ConfigRecord struct {
	ConfigName  string    `json:"configName"`
	AccountID   string    `json:"accountID"`
	OrgID       string    `json:"orgID"`
	LastUpdated time.Time `json:"lastUpdated"`
}

// NewConfigTopologyState initializes a new ConfigTopologyState
func NewConfigTopologyState(destGateway GatewayInfo) (*ConfigTopologyState, error) {
	return &ConfigTopologyState{
		ConfigTopology: configTopology{
			DestGateway: destGateway,
			RouteTable:  make(map[GatewayInfo]time.Time),
		},
		mux: sync.Mutex{},
	}, nil
}

// UpsertRoute upserts given route.
func (ts *ConfigTopologyState) UpsertRoute(_ context.Context, gw GatewayInfo) {
	ts.mux.Lock()
	defer ts.mux.Unlock()

	ts.ConfigTopology.RouteTable[gw] = time.Now()
}

// ResettableConfigTopologyRegistry is a concrete version of TopologyDataRegistry that is able to be reset.
type ResettableConfigTopologyRegistry struct {
	topology        *sync.Map
	setIntervalChan chan time.Duration
}

// NewResettableConfigTopologyRegistry creates a new ResettableConfigTopologyRegistry
func NewResettableConfigTopologyRegistry() *ResettableConfigTopologyRegistry {
	return &ResettableConfigTopologyRegistry{
		topology: &sync.Map{},
	}
}

// RegisterConfigTopology registers the ConfigTopology with the registry.
func (rtsr *ResettableConfigTopologyRegistry) RegisterConfigTopology(processorID string, configTopology *ConfigTopologyState) error {
	_, alreadyExists := rtsr.topology.LoadOrStore(processorID, configTopology)
	if alreadyExists {
		return fmt.Errorf("topology for processor %q was already registered", processorID)
	}

	return nil
}

// Reset unregisters all topology states in this registry
func (rtsr *ResettableConfigTopologyRegistry) Reset() {
	rtsr.topology = &sync.Map{}
}

// SetIntervalChan returns the setIntervalChan
func (rtsr *ResettableConfigTopologyRegistry) SetIntervalChan() chan time.Duration {
	return rtsr.setIntervalChan
}

// TopologyInfos returns all the topology data in this registry.
func (rtsr *ResettableConfigTopologyRegistry) TopologyInfos() []ConfigTopologyInfo {
	states := []configTopology{}

	rtsr.topology.Range(func(_, value any) bool {
		ts := value.(*ConfigTopologyState)
		states = append(states, ts.ConfigTopology)
		return true
	})

	ti := []ConfigTopologyInfo{}
	for _, ts := range states {
		curInfo := ConfigTopologyInfo{}
		curInfo.GatewayID = ts.DestGateway.GatewayID
		curInfo.ConfigName = ts.DestGateway.ConfigName
		curInfo.AccountID = ts.DestGateway.AccountID
		curInfo.OrgID = ts.DestGateway.OrgID
		for gw, updated := range ts.RouteTable {
			curInfo.SourceConfigs = append(curInfo.SourceConfigs, ConfigRecord{
				ConfigName:  gw.ConfigName,
				AccountID:   gw.AccountID,
				OrgID:       gw.OrgID,
				LastUpdated: updated.UTC(),
			})
		}
		if len(curInfo.SourceConfigs) > 0 {
			ti = append(ti, curInfo)
		}
	}

	return ti
}
