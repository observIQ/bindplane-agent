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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/observiq/observiq-otel-collector/opamp"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestGetDefaultCollectorConfigPathENV(t *testing.T) {
	fakeConfigPath := "./fake/path/config.yaml"

	t.Setenv(configPathENV, fakeConfigPath)

	expected := []string{fakeConfigPath}
	actual := getDefaultCollectorConfigPaths()

	require.Equal(t, expected, actual)
}

func TestGetDefaultCollectorConfigPathNone(t *testing.T) {
	expected := []string{"./config.yaml"}
	actual := getDefaultCollectorConfigPaths()

	require.Equal(t, expected, actual)
}

func TestGetDefaultLoggingConfigPathENV(t *testing.T) {
	fakeLoggingPath := "./fake/path/logging.yaml"

	t.Setenv(loggingPathENV, fakeLoggingPath)

	expected := fakeLoggingPath
	actual := getDefaultLoggingConfigPath()

	require.Equal(t, expected, actual)
}

func TestGetDefaultLoggingConfigPathNone(t *testing.T) {
	expected := "./logging.yaml"
	actual := getDefaultLoggingConfigPath()

	require.Equal(t, expected, actual)
}

func TestGetDefaultManagerConfigPathENV(t *testing.T) {
	fakeManagerPath := "./fake/path/Manager.yaml"

	t.Setenv(managerPathENV, fakeManagerPath)

	expected := fakeManagerPath
	actual := getDefaultManagerConfigPath()

	require.Equal(t, expected, actual)
}

func TestGetDefaultManagerConfigPathNone(t *testing.T) {
	expected := "./manager.yaml"
	actual := getDefaultManagerConfigPath()

	require.Equal(t, expected, actual)
}

func TestCheckManagerNoConfig(t *testing.T) {
	manager := "./manager.yaml"
	err := checkManagerConfig(&manager)
	require.ErrorIs(t, err, os.ErrNotExist)

	tmp := "\000"
	err = checkManagerConfig(&tmp)
	require.Error(t, err)
}

func TestCheckManagerConfigNoFile(t *testing.T) {
	t.Setenv(endpointENV, "0.0.0.0")

	t.Setenv(agentNameENV, "agent name")

	t.Setenv(agentIDENV, "agent ID")

	t.Setenv(secretkeyENV, "secretKey")

	t.Setenv(labelsENV, "this is a label")

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

func TestManagerConfigNoAgentIDWillSet(t *testing.T) {
	tmpDir := t.TempDir()
	manager := filepath.Join(tmpDir, "manager.yaml")
	data := []byte("")
	require.NoError(t, os.WriteFile(manager, data, 0600))
	err := checkManagerConfig(&manager)
	require.NoError(t, err)

	cfgBytes, err := ioutil.ReadFile(manager)
	require.NoError(t, err)

	var config opamp.Config
	require.NoError(t, yaml.Unmarshal(cfgBytes, &config))
	require.NotEmpty(t, config.AgentID)
	uuidv4, err := uuid.Parse(config.AgentID)
	require.NoError(t, err)
	require.NotEmpty(t, uuidv4)
}

func TestManagerConfigWillNotOverwriteCurrentAgentID(t *testing.T) {
	tmpDir := t.TempDir()
	manager := filepath.Join(tmpDir, "manager.yaml")

	id := uuid.NewString()
	data := []byte(fmt.Sprintf(`
---
agent_id: %s
`, id))
	require.NoError(t, os.WriteFile(manager, data, 0600))
	err := checkManagerConfig(&manager)
	require.NoError(t, err)

	cfgBytes, err := ioutil.ReadFile(manager)
	require.NoError(t, err)

	var config opamp.Config
	require.NoError(t, yaml.Unmarshal(cfgBytes, &config))
	require.Equal(t, config.AgentID, id)
}

func TestManagerConfigWillErrorOnInvalidOpAmpConfig(t *testing.T) {
	tmpDir := t.TempDir()
	manager := filepath.Join(tmpDir, "manager.yaml")
	data := []byte(`
---
agent_id:
  - some-kind-of-array
  - should-blow-up
`)
	require.NoError(t, os.WriteFile(manager, data, 0600))
	err := checkManagerConfig(&manager)
	require.Error(t, err)
	require.ErrorContains(t, err, "unable to interpret config file")
}

func TestManagerConfigCheckFileModes(t *testing.T) {
	tmpDir := t.TempDir()

	testCases := []struct {
		name        string
		fileMode    os.FileMode
		expectedErr error
	}{
		{
			name:        "read_only",
			fileMode:    0400,
			expectedErr: errors.New("failed to rewrite manager config with identifying fields"),
		},
		{
			name:     "valid_read_write",
			fileMode: 0600,
		},
	}

	for idx, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			manager := filepath.Join(tmpDir, fmt.Sprintf("manager-%d.yaml", idx))
			require.NoError(t, os.WriteFile(manager, []byte(""), tc.fileMode))
			err := checkManagerConfig(&manager)
			if tc.expectedErr != nil {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
