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

package googlemanagedprometheusexporter

import (
	"testing"

	gmp "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/googlemanagedprometheusexporter"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
)

func TestCreateDefaultConfig(t *testing.T) {
	cfg := createDefaultConfig()
	googleCfg, ok := cfg.(*Config)
	require.True(t, ok)

	require.Equal(t, defaultUserAgent, googleCfg.GCPConfig.UserAgent)
	require.Nil(t, googleCfg.Validate())
}

func TestSetClientOptionsWithCredentials(t *testing.T) {
	testCases := []struct {
		name   string
		config *Config
		opts   []option.ClientOption
	}{
		{
			name: "With no credentials",
			config: &Config{
				GCPConfig: &gmp.Config{},
			},
			opts: []option.ClientOption{},
		},
		{
			name: "With credentials json",
			config: &Config{
				Credentials: "testjson",
				GCPConfig:   &gmp.Config{},
			},
			opts: []option.ClientOption{
				option.WithCredentialsJSON([]byte("testjson")),
			},
		},
		{
			name: "With credentials file",
			config: &Config{
				CredentialsFile: "testfile",
				GCPConfig:       &gmp.Config{},
			},
			opts: []option.ClientOption{
				option.WithCredentialsFile("testfile"),
			},
		},
		{
			name: "With both credentials json and credentials file",
			config: &Config{
				Credentials:     "testjson",
				CredentialsFile: "testfile",
				GCPConfig:       &gmp.Config{},
			},
			opts: []option.ClientOption{
				option.WithCredentialsJSON([]byte("testjson")),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Nil(t, tc.config.GCPConfig.MetricConfig.ClientConfig.GetClientOptions)

			tc.config.setClientOptions()
			require.NotNil(t, tc.config.GCPConfig.MetricConfig.ClientConfig.GetClientOptions)

			opts := tc.config.GCPConfig.MetricConfig.ClientConfig.GetClientOptions()
			require.Equal(t, tc.opts, opts)
		})
	}
}

func TestSetProject(t *testing.T) {
	testCases := []struct {
		name            string
		config          *Config
		expectedProject string
		expectedErr     string
	}{
		{
			name: "With project already set",
			config: &Config{
				GCPConfig: &gmp.Config{
					GMPConfig: gmp.GMPConfig{
						ProjectID: "test",
					},
				},
			},
			expectedProject: "test",
		},
		{
			name: "With project in json credentials",
			config: &Config{
				Credentials: `{"project_id":"test"}`,
				GCPConfig:   &gmp.Config{},
			},
			expectedProject: "test",
		},
		{
			name: "With missing json key",
			config: &Config{
				Credentials: `{"test":"value"}`,
				GCPConfig:   &gmp.Config{},
			},
			expectedErr: "project id does not exist",
		},
		{
			name: "With invalid json",
			config: &Config{
				Credentials: `{`,
				GCPConfig:   &gmp.Config{},
			},
			expectedErr: "failed to unmarshal credentials",
		},
		{
			name: "With invalid string",
			config: &Config{
				Credentials: `{"project_id":100}`,
				GCPConfig:   &gmp.Config{},
			},
			expectedErr: "project id is not a string",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.setProject()
			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			}

			require.Equal(t, tc.expectedProject, tc.config.GCPConfig.ProjectID)
		})
	}
}
