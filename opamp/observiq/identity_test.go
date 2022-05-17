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
	"runtime"
	"testing"

	"github.com/observiq/observiq-otel-collector/internal/version"
	"github.com/observiq/observiq-otel-collector/opamp"
	"github.com/open-telemetry/opamp-go/protobufs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func Test_newIdentity(t *testing.T) {
	secretKeyContents := "b92222ee-a1fc-4bb1-98db-26de3448541b"
	labelsContents := "one=foo,two=bar"
	agentNameContents := "My Agent"

	cfg := opamp.Config{
		Endpoint:  "ws://localhost:1234",
		SecretKey: &secretKeyContents,
		AgentID:   "8321f735-a52c-4f49-aca9-66f9266c5fe5",
		Labels:    &labelsContents,
		AgentName: &agentNameContents,
	}

	got := newIdentity(zap.NewNop().Sugar(), cfg)

	// Check all fields from config
	require.Equal(t, cfg.AgentID, got.agentID)
	require.Equal(t, cfg.AgentName, got.agentName)
	require.Equal(t, cfg.Labels, got.labels)

	// Check fields that must not be empty
	require.NotEmpty(t, got.oSDetails)
	require.NotEmpty(t, got.hostname)
	require.NotEmpty(t, got.mac)

	// Check hardcoded/fields from runtime and other packages
	require.Equal(t, got.serviceName, "com.observiq.collector")
	require.Equal(t, got.version, version.Version())
	require.Equal(t, got.oSArch, runtime.GOARCH)
	require.Equal(t, got.oSFamily, runtime.GOOS)
}

func TestToAgentDescription(t *testing.T) {
	labelsContents := "one=foo,two=bar"
	agentNameContents := "My Agent"
	testCases := []struct {
		desc     string
		ident    *identity
		expected *protobufs.AgentDescription
	}{
		{
			desc: "Missing Agent Name and labels",
			ident: &identity{
				agentID:     "4322d8d1-f3e0-46db-b68d-b01a4689ef19",
				agentName:   nil,
				serviceName: "com.observiq.collector",
				version:     "v1.2.3",
				labels:      nil,
				oSArch:      "amd64",
				oSDetails:   "os details",
				oSFamily:    "linux",
				hostname:    "my-linux-box",
				mac:         "68-C7-B4-EB-A8-D2",
			},
			expected: &protobufs.AgentDescription{
				IdentifyingAttributes: []*protobufs.KeyValue{
					opamp.StringKeyValue("service.instance.id", "4322d8d1-f3e0-46db-b68d-b01a4689ef19"),
					opamp.StringKeyValue("service.name", "com.observiq.collector"),
					opamp.StringKeyValue("service.version", "v1.2.3"),
					opamp.StringKeyValue("service.instance.name", "my-linux-box"),
				},
				NonIdentifyingAttributes: []*protobufs.KeyValue{
					opamp.StringKeyValue("os.arch", "amd64"),
					opamp.StringKeyValue("os.details", "os details"),
					opamp.StringKeyValue("os.family", "linux"),
					opamp.StringKeyValue("host.name", "my-linux-box"),
					opamp.StringKeyValue("host.mac_address", "68-C7-B4-EB-A8-D2"),
				},
			},
		},
		{
			desc: "With Agent Name and labels",
			ident: &identity{
				agentID:     "4322d8d1-f3e0-46db-b68d-b01a4689ef19",
				agentName:   &agentNameContents,
				serviceName: "com.observiq.collector",
				version:     "v1.2.3",
				labels:      &labelsContents,
				oSArch:      "amd64",
				oSDetails:   "os details",
				oSFamily:    "linux",
				hostname:    "my-linux-box",
				mac:         "68-C7-B4-EB-A8-D2",
			},
			expected: &protobufs.AgentDescription{
				IdentifyingAttributes: []*protobufs.KeyValue{
					opamp.StringKeyValue("service.instance.id", "4322d8d1-f3e0-46db-b68d-b01a4689ef19"),
					opamp.StringKeyValue("service.name", "com.observiq.collector"),
					opamp.StringKeyValue("service.version", "v1.2.3"),
					opamp.StringKeyValue("service.instance.name", agentNameContents),
				},
				NonIdentifyingAttributes: []*protobufs.KeyValue{
					opamp.StringKeyValue("os.arch", "amd64"),
					opamp.StringKeyValue("os.details", "os details"),
					opamp.StringKeyValue("os.family", "linux"),
					opamp.StringKeyValue("host.name", "my-linux-box"),
					opamp.StringKeyValue("host.mac_address", "68-C7-B4-EB-A8-D2"),
					opamp.StringKeyValue("service.labels", labelsContents),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actual := tc.ident.ToAgentDescription()
			assert.Equal(t, tc.expected, actual)
		})
	}
}
