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

package opamp

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManagedConfigGetCurrentConfigHash(t *testing.T) {
	managedConfig := &ManagedConfig{
		ConfigPath:        "./path.yaml",
		Reload:            NoopReloadFunc,
		currentConfigHash: []byte("hello world"),
	}

	actual := managedConfig.GetCurrentConfigHash()
	require.Equal(t, managedConfig.currentConfigHash, actual)
}

func TestManagedConfigComputeConfigHash(t *testing.T) {
	testCases := []struct {
		desc     string
		testFunc func(*testing.T)
	}{
		{
			desc: "Missing config file",
			testFunc: func(t *testing.T) {
				managedConfig := &ManagedConfig{
					ConfigPath: "./path.yaml",
				}

				err := managedConfig.ComputeConfigHash()
				assert.Error(t, err)
			},
		},
		{
			desc: "Successful hash",
			testFunc: func(t *testing.T) {
				tmpDir := t.TempDir()
				cfgPath := filepath.Join(tmpDir, "config.yaml")
				err := os.WriteFile(cfgPath, []byte("hello world"), 0600)
				assert.NoError(t, err)

				managedConfig := &ManagedConfig{
					ConfigPath: cfgPath,
				}

				expected := []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0xe3, 0xb0, 0xc4, 0x42, 0x98, 0xfc, 0x1c, 0x14, 0x9a, 0xfb, 0xf4, 0xc8, 0x99, 0x6f, 0xb9, 0x24, 0x27, 0xae, 0x41, 0xe4, 0x64, 0x9b, 0x93, 0x4c, 0xa4, 0x95, 0x99, 0x1b, 0x78, 0x52, 0xb8, 0x55}

				err = managedConfig.ComputeConfigHash()
				assert.NoError(t, err)
				assert.Equal(t, expected, managedConfig.currentConfigHash)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, tc.testFunc)
	}
}

func TestNewManagedConfig(t *testing.T) {
	testCases := []struct {
		desc     string
		testFunc func(*testing.T)
	}{
		{
			desc: "Missing config file",
			testFunc: func(t *testing.T) {
				managedConfig, err := NewManagedConfig("./path.yml", NoopReloadFunc)
				assert.Error(t, err)
				assert.Nil(t, managedConfig)
			},
		},
		{
			desc: "Success",
			testFunc: func(t *testing.T) {
				tmpDir := t.TempDir()
				cfgPath := filepath.Join(tmpDir, "config.yaml")
				err := os.WriteFile(cfgPath, []byte("hello world"), 0600)
				assert.NoError(t, err)

				expected := &ManagedConfig{
					ConfigPath:        cfgPath,
					Reload:            NoopReloadFunc,
					currentConfigHash: []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0xe3, 0xb0, 0xc4, 0x42, 0x98, 0xfc, 0x1c, 0x14, 0x9a, 0xfb, 0xf4, 0xc8, 0x99, 0x6f, 0xb9, 0x24, 0x27, 0xae, 0x41, 0xe4, 0x64, 0x9b, 0x93, 0x4c, 0xa4, 0x95, 0x99, 0x1b, 0x78, 0x52, 0xb8, 0x55},
				}

				managedConfig, err := NewManagedConfig(cfgPath, NoopReloadFunc)
				assert.NoError(t, err)
				assert.Equal(t, expected.ConfigPath, managedConfig.ConfigPath)
				assert.Equal(t, expected.currentConfigHash, managedConfig.currentConfigHash)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, tc.testFunc)
	}
}
