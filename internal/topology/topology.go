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

// TopologyStateRegistry represents a registry for the topology processor to register their TopologyState.
type TopologyStateRegistry interface {
	// RegisterTopologyState registers the topology state for the given processor.
	// It should return an error if the processor has already been registered.
	RegisterTopologyState(processorID string, data *TopologyState) error
}

// GatewayConfigInfo reprents the unique identifiable information about a bindplane gateway's configuration
type GatewayConfigInfo struct {
	ConfigName string
	AccountID  string
	OrgID      string
}

// TopologyState represents the data captured through topology processors.
type TopologyState struct {
	destGateway GatewayConfigInfo
	routeTable  map[GatewayConfigInfo]time.Time
}

type TopologyMessage struct {
	DestGateway    GatewayConfigInfo `json:"destGateway"`
	SourceGateways []GatewayState    `json:"sourceGateways"`
}

type GatewayState struct {
	Gateway GatewayConfigInfo `json:"gateway"`
	Updated time.Time         `json:"updated"`
}

// NewTopologyState initializes a new TopologyState
func NewTopologyState(destGateway GatewayConfigInfo) (*TopologyState, error) {
	return &TopologyState{
		destGateway: destGateway,
		routeTable:  make(map[GatewayConfigInfo]time.Time),
	}, nil
}

// UpsertRoute upserts given route.
func (ts *TopologyState) UpsertRoute(ctx context.Context, gw GatewayConfigInfo) {
	fmt.Println("\033[34m UPSERT ROUTE \033[0m", gw)
	ts.routeTable[gw] = time.Now()
}

// ResettableTopologyStateRegistry is a concrete version of TopologyDataRegistry that is able to be reset.
type ResettableTopologyStateRegistry struct {
	topology *sync.Map
}

// NewResettableTopologyStateRegistry creates a new ResettableTopologyStateRegistry
func NewResettableTopologyStateRegistry() *ResettableTopologyStateRegistry {
	return &ResettableTopologyStateRegistry{
		topology: &sync.Map{},
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

// TopologyMessages returns all the topology states in this registry.
func (rtsr *ResettableTopologyStateRegistry) TopologyMessages() []TopologyMessage {
	states := []TopologyState{}

	rtsr.topology.Range(func(_, value any) bool {
		ts := value.(*TopologyState)
		states = append(states, *ts)
		return true
	})

	messages := []TopologyMessage{}
	for _, ts := range states {
		curMessage := TopologyMessage{}
		curMessage.DestGateway = ts.destGateway
		for gw, updated := range ts.routeTable {
			curMessage.SourceGateways = append(curMessage.SourceGateways, GatewayState{
				Gateway: gw,
				Updated: updated,
			})
		}
		messages = append(messages, curMessage)
	}

	return messages
}
