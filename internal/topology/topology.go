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
	"slices"
	"sync"
	"time"
)

// TopologyStateRegistry represents a registry for the topology processor to register their TopologyState.
type TopologyStateRegistry interface {
	// RegisterTopologyState registers the topology state for the given processor.
	// It should return an error if the processor has already been registered.
	RegisterTopologyState(processorID string, data *TopologyState) error
	SetIntervalChan() chan time.Duration
	Reset()
}

// GatewayConfigInfo reprents the unique identifiable information about a bindplane gateway's configuration
type GatewayConfigInfo struct {
	ConfigName string
	AccountID  string
	OrgID      string
}

// TopologyState represents the data captured through topology processors.
type TopologyState struct {
	destConfig GatewayConfigInfo
	routeTable map[GatewayConfigInfo]time.Time
}

// TopologyInfo represents topology relationships between configs.
type TopologyInfo struct {
	ConfigName    string       `json:"configName"`
	AccountID     string       `json:"accountID"`
	OrgID         string       `json:"orgID"`
	SourceConfigs []ConfigInfo `json:"sourceConfigs"`
}

// ConfigInfo represents a source config and the time it was last detected
type ConfigInfo struct {
	ConfigName  string    `json:"configName"`
	AccountID   string    `json:"accountID"`
	OrgID       string    `json:"orgID"`
	LastUpdated time.Time `json:"lastUpdated"`
}

// NewTopologyState initializes a new TopologyState
func NewTopologyState(destGateway GatewayConfigInfo, interval time.Duration) (*TopologyState, error) {
	return &TopologyState{
		destConfig: destGateway,
		routeTable: make(map[GatewayConfigInfo]time.Time),
	}, nil
}

// UpsertRoute upserts given route.
func (ts *TopologyState) UpsertRoute(ctx context.Context, gw GatewayConfigInfo) {
	fmt.Println("\033[34m UPSERT ROUTE \033[0m", gw)
	ts.routeTable[gw] = time.Now()
}

// ResettableTopologyStateRegistry is a concrete version of TopologyDataRegistry that is able to be reset.
type ResettableTopologyStateRegistry struct {
	topology        *sync.Map
	setIntervalChan chan time.Duration
}

// NewResettableTopologyStateRegistry creates a new ResettableTopologyStateRegistry
func NewResettableTopologyStateRegistry() *ResettableTopologyStateRegistry {
	return &ResettableTopologyStateRegistry{
		topology:        &sync.Map{},
		setIntervalChan: make(chan time.Duration, 1),
	}
}

// RegisterTopologyState registers the TopologyState with the registry.
func (rtsr *ResettableTopologyStateRegistry) RegisterTopologyState(processorID string, topologyState *TopologyState) error {
	_, alreadyExists := rtsr.topology.LoadOrStore(processorID, topologyState)
	if alreadyExists {
		return fmt.Errorf("topology for processor %q was already registered", processorID)
	}

	return nil
}

// Reset unregisters all topology states in this registry
func (rtsr *ResettableTopologyStateRegistry) Reset() {
	rtsr.topology = &sync.Map{}
}

// SetIntervalChan returns the setIntervalChan
func (rtsr *ResettableTopologyStateRegistry) SetIntervalChan() chan time.Duration {
	return rtsr.setIntervalChan
}

// TopologyInfos returns all the topology data in this registry.
func (rtsr *ResettableTopologyStateRegistry) TopologyInfos() []TopologyInfo {
	states := []TopologyState{}

	rtsr.topology.Range(func(_, value any) bool {
		ts := value.(*TopologyState)
		states = append(states, *ts)
		return true
	})

	ti := []TopologyInfo{}
	for _, ts := range states {
		curInfo := TopologyInfo{}
		curInfo.ConfigName = ts.destConfig.ConfigName
		curInfo.AccountID = ts.destConfig.AccountID
		curInfo.OrgID = ts.destConfig.OrgID
		for gw, updated := range ts.routeTable {
			curInfo.SourceConfigs = append(curInfo.SourceConfigs, ConfigInfo{
				ConfigName:  gw.ConfigName,
				AccountID:   gw.AccountID,
				OrgID:       gw.OrgID,
				LastUpdated: updated.UTC(),
			})
		}
		if len(curInfo.SourceConfigs) > 0 {
			slices.SortFunc(curInfo.SourceConfigs, func(a ConfigInfo, b ConfigInfo) int {
				if a.OrgID < b.OrgID {
					return -1
				}
				if a.OrgID > b.OrgID {
					return 1
				}
				if a.AccountID < b.AccountID {
					return -1
				}
				if a.AccountID > b.AccountID {
					return 1
				}
				if a.ConfigName < b.ConfigName {
					return -1
				}
				if a.ConfigName > b.ConfigName {
					return 1
				}
				return 0
			})
			ti = append(ti, curInfo)
		}
	}

	return ti
}
