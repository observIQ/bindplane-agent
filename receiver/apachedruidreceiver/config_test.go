//// Copyright  observIQ, Inc.
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

package apachedruidreceiver

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/observiq/bindplane-agent/receiver/apachedruidreceiver/internal/metadata"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/collector/confmap/confmaptest"
)

const (
	validEndpoint = "0.0.0.0:12345"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		name         string
		config       *Config
		expectedErrs string
	}{
		{
			name: "Valid config without TLS or basic auth",
			config: &Config{
				Metrics: MetricsConfig{
					Endpoint: validEndpoint,
				},
			},
		},
		{
			name: "Valid config with TLS and basic auth",
			config: &Config{
				Metrics: MetricsConfig{
					Endpoint: validEndpoint,
					BasicAuth: &BasicAuth{
						Username: validUsername,
						Password: validPassword,
					},
					TLS: &configtls.TLSServerSetting{
						TLSSetting: configtls.TLSSetting{
							CertFile: "some_cert_file",
							KeyFile:  "some_key_file",
						},
					},
				},
			},
		},
		{
			name: "No endpoint config",
			config: &Config{
				Metrics: MetricsConfig{},
			},
			expectedErrs: errNoEndpoint.Error(),
		},
		{
			name: "Missing port config",
			config: &Config{
				Metrics: MetricsConfig{
					Endpoint: "www.google.com",
				},
			},
			expectedErrs: "failed to split endpoint into 'host:port' pair: address www.google.com: missing port in address",
		},
		{
			name: "TLS config missing cert",
			config: &Config{
				Metrics: MetricsConfig{
					Endpoint: validEndpoint,
					TLS: &configtls.TLSServerSetting{
						TLSSetting: configtls.TLSSetting{
							KeyFile: "some_key_file",
						},
					},
				},
			},
			expectedErrs: errNoCert.Error(),
		},
		{
			name: "TLS config missing key",
			config: &Config{
				Metrics: MetricsConfig{
					Endpoint: validEndpoint,
					TLS: &configtls.TLSServerSetting{
						TLSSetting: configtls.TLSSetting{
							CertFile: "some_cert_file",
						},
					},
				},
			},
			expectedErrs: errNoKey.Error(),
		},
		{
			name: "Empty TLS config",
			config: &Config{
				Metrics: MetricsConfig{
					Endpoint: validEndpoint,
					TLS: &configtls.TLSServerSetting{
						TLSSetting: configtls.TLSSetting{},
					},
				},
			},
			expectedErrs: createMultiErr(errNoCert.Error(), errNoKey.Error()),
		},
		{
			name: "Basic Auth missing username",
			config: &Config{
				Metrics: MetricsConfig{
					Endpoint: validEndpoint,
					BasicAuth: &BasicAuth{
						Password: validPassword,
					},
				},
			},
			expectedErrs: errNoUsername.Error(),
		},
		{
			name: "Basic Auth missing password",
			config: &Config{
				Metrics: MetricsConfig{
					Endpoint: validEndpoint,
					BasicAuth: &BasicAuth{
						Username: validUsername,
					},
				},
			},
			expectedErrs: errNoPass.Error(),
		},
		{
			name: "Empty basic auth",
			config: &Config{
				Metrics: MetricsConfig{
					Endpoint:  validEndpoint,
					BasicAuth: &BasicAuth{},
				},
			},
			expectedErrs: createMultiErr(errNoUsername.Error(), errNoPass.Error()),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			errs := tc.config.Validate()
			if tc.expectedErrs == "" {
				require.NoError(t, errs)
			} else {
				require.ErrorContains(t, errs, tc.expectedErrs)
			}
		})
	}
}

func TestLoadTLSConfig(t *testing.T) {
	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "configTLS.yaml"))
	require.NoError(t, err)

	cases := []struct {
		name           string
		expectedConfig component.Config
	}{
		{
			name: "",
			expectedConfig: &Config{
				Metrics: MetricsConfig{
					MetricsBuilderConfig: metadata.DefaultMetricsBuilderConfig(),
					Endpoint:             validEndpoint,
					TLS: &configtls.TLSServerSetting{
						TLSSetting: configtls.TLSSetting{
							CertFile: "some_cert_file",
							KeyFile:  "some_key_file",
						},
					},
					BasicAuth: &BasicAuth{
						Username: validUsername,
						Password: validPassword,
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			factory := NewFactory()
			cfg := factory.CreateDefaultConfig()

			loaded, err := cm.Sub(component.NewIDWithName(typeStr, tc.name).String())
			require.NoError(t, err)
			require.NoError(t, component.UnmarshalConfig(loaded, cfg))
			require.Equal(t, tc.expectedConfig, cfg)
			require.NoError(t, component.ValidateConfig(cfg))
		})
	}
}

func TestLoadConfig(t *testing.T) {
	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)

	cases := []struct {
		name           string
		expectedConfig component.Config
	}{
		{
			name: "",
			expectedConfig: &Config{
				Metrics: MetricsConfig{
					MetricsBuilderConfig: metadata.DefaultMetricsBuilderConfig(),
					Endpoint:             validEndpoint,
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			factory := NewFactory()
			cfg := factory.CreateDefaultConfig()

			loaded, err := cm.Sub(component.NewIDWithName(typeStr, tc.name).String())
			require.NoError(t, err)
			require.NoError(t, component.UnmarshalConfig(loaded, cfg))
			require.Equal(t, tc.expectedConfig, cfg)
			require.NoError(t, component.ValidateConfig(cfg))
		})
	}
}

func createMultiErr(errors ...string) string {
	return strings.Join(errors, "; ")
}
