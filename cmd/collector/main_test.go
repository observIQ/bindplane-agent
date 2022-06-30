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

package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/observiq/observiq-otel-collector/opamp"
	"github.com/stretchr/testify/require"
)

func TestCheckManagerNoConfig(t *testing.T) {
	manager := "./manager.yaml"
	err := checkManagerConfig(&manager)
	require.ErrorIs(t, err, os.ErrNotExist)

	tmp := "\000"
	err = checkManagerConfig(&tmp)
	require.Error(t, err)
}

func TestCheckManagerConfigNoFile(t *testing.T) {
	os.Setenv(endpointENV, "0.0.0.0")
	defer os.Unsetenv(endpointENV)

	os.Setenv(agentNameENV, "agent name")
	defer os.Unsetenv(agentNameENV)

	os.Setenv(agentIDENV, "agent ID")
	defer os.Unsetenv(agentIDENV)

	os.Setenv(secretkeyENV, "secretKey")
	defer os.Unsetenv(secretkeyENV)

	os.Setenv(labelsENV, "this is a label")
	defer os.Unsetenv(labelsENV)

	tmpdir := t.TempDir()
	manager := filepath.Join(tmpdir, "manager.yaml")
	err := checkManagerConfig(&manager)
	require.NoError(t, err)

	actual, _ := opamp.ParseConfig(manager)
	expected := &opamp.Config{
		Endpoint:  "0.0.0.0",
		AgentID:   "agent ID",
		AgentName: new(string),
		SecretKey: new(string),
		Labels:    new(string),
	}
	*expected.AgentName = "agent name"
	*expected.SecretKey = "secretKey"
	*expected.Labels = "this is a label"

	require.Equal(t, expected, actual)
}

func TestCheckManagerConfig(t *testing.T) {
	tmpdir := t.TempDir()
	manager := filepath.Join(tmpdir, "manager.yaml")

	data := []byte("temporary directory")
	os.WriteFile(manager, data, 0600)
	err := checkManagerConfig(&manager)
	require.NoError(t, err)
}
