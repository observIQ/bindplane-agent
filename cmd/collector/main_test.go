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
	"os/exec"
	"testing"

	"github.com/observiq/observiq-otel-collector/opamp"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestCheckManagerConfigNoFile(t *testing.T) {
	exec.Command("rm", "-r", "./manager.yaml").Run()
	manager := "./manager.yaml"
	err := checkManagerConfig(&manager)
	require.Error(t, err)

	os.Setenv(endpoint, "0.0.0.0")
	defer os.Unsetenv(endpoint)

	os.Setenv(secretKey, "secretKey")
	defer os.Unsetenv(secretKey)

	os.Setenv(labels, "this is a label")
	defer os.Unsetenv(labels)
	defer os.Unsetenv(agentID)

	manager = "./manager.yaml"
	err = checkManagerConfig(&manager)
	require.NoError(t, err)

	dat, err := os.ReadFile("./manager.yaml")
	out := &opamp.Config{}
	err = yaml.Unmarshal(dat, out)
	require.Equal(t,
		&opamp.Config{
			Endpoint: "0.0.0.0",
		},
		&opamp.Config{
			Endpoint: out.Endpoint,
		})
}

func TestCheckManagerConfig(t *testing.T) {
	manager := "./manager.yaml"
	err := checkManagerConfig(&manager)
	require.NoError(t, err)
}
