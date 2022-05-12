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

package observiq

import (
	"fmt"
	"runtime"

	ios "github.com/observiq/observiq-otel-collector/internal/os"
	"github.com/observiq/observiq-otel-collector/internal/version"
	"github.com/observiq/observiq-otel-collector/opamp"
	"github.com/open-telemetry/opamp-go/client"
	"github.com/open-telemetry/opamp-go/protobufs"
	"go.uber.org/zap"
)

// identity contains identifying information about the Collector
type identity struct {
	agentID     string
	agentName   *string
	serviceName string
	version     string
	labels      *string
	oSArch      string
	oSDetails   string
	oSFamily    string
	hostname    string
	mac         string
}

// newIdentity constructs a new identity for this collector
func newIdentity(logger *zap.SugaredLogger, config opamp.Config) *identity {
	// Grab various fields from OS
	hostname, err := ios.Hostname()
	if err != nil {
		logger.Warn("Failed to retrieve hostname for collector. Creating partial identity", zap.Error(err))
	}

	name, err := ios.Name()
	if err != nil {
		logger.Warn("Failed to retrieve host details on collector. Creating partial identity", zap.Error(err))
	}

	return &identity{
		agentID:     config.AgentID,
		agentName:   config.AgentName,
		serviceName: "com.observiq.collector", // TODO figure out if this should be hardcoded like so or read from system.
		version:     version.Version(),
		labels:      config.Labels,
		oSArch:      runtime.GOARCH,
		oSDetails:   name,
		oSFamily:    runtime.GOOS,
		hostname:    hostname,
		mac:         ios.MACAddress(),
	}
}

func (i *identity) ToAgentDescription() (*protobufs.AgentDescription, error) {
	identifyingAttributes := []*protobufs.KeyValue{
		opamp.StringKeyValue("service.instance.id", i.agentID),
		opamp.StringKeyValue("service.name", i.serviceName),
		opamp.StringKeyValue("service.version", i.version),
	}

	if i.agentName != nil {
		identifyingAttributes = append(identifyingAttributes, opamp.StringKeyValue("service.instance.name", *i.agentName))
	} else {
		identifyingAttributes = append(identifyingAttributes, opamp.StringKeyValue("service.instance.name", i.hostname))
	}

	nonIdentifyingAttributes := []*protobufs.KeyValue{
		opamp.StringKeyValue("os.arch", i.oSArch),
		opamp.StringKeyValue("os.details", i.oSDetails),
		opamp.StringKeyValue("os.family", i.oSFamily),
		opamp.StringKeyValue("host.name", i.hostname),
		opamp.StringKeyValue("host.mac_address", i.mac),
	}

	if i.labels != nil {
		nonIdentifyingAttributes = append(nonIdentifyingAttributes, opamp.StringKeyValue("service.labels", *i.labels))
	}

	agentDesc := &protobufs.AgentDescription{
		IdentifyingAttributes:    identifyingAttributes,
		NonIdentifyingAttributes: nonIdentifyingAttributes,
	}

	// Compute hash
	if err := client.CalcHashAgentDescription(agentDesc); err != nil {
		// Still return agentDesc it will just be missing a hash
		return agentDesc, fmt.Errorf("error while computing agent description hash: %w", err)
	}

	return agentDesc, nil
}
