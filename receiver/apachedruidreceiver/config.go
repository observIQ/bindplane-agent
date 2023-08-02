// Copyright The OpenTelemetry Authors
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

// Package apachedruidreceiver provides a receiver that receives telemetry from an Apache Druid instance.
package apachedruidreceiver

import (
	"errors"
	"fmt"
	"net"

	"github.com/observiq/bindplane-agent/receiver/apachedruidreceiver/internal/metadata"
	"go.opentelemetry.io/collector/config/configtls"
	"go.uber.org/multierr"
)

// Config is the configuration of an Apache Druid receiver
type Config struct {
	Metrics MetricsConfig `mapstructure:"metrics"`
}

// MetricsConfig is the metrics portion of the configuration of an Apache Druid receiver
type MetricsConfig struct {
	BasicAuth            *BasicAuth                    `mapstructure:"basic_auth"`
	Endpoint             string                        `mapstructure:"endpoint"`
	TLS                  *configtls.TLSServerSetting   `mapstructure:"tls"`
	MetricsBuilderConfig metadata.MetricsBuilderConfig `mapstructure:",squash"`
}

// BasicAuth is basic username-password authentication
type BasicAuth struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

var (
	errNoEndpoint = errors.New("an endpoint must be specified")
	errNoCert     = errors.New("tls was configured, but no cert file was specified")
	errNoKey      = errors.New("tls was configured, but no key file was specified")
	errNoUsername = errors.New("basic_auth was configured, but no username was specified")
	errNoPass     = errors.New("basic_auth was configured, but no password was specified")
)

// Validate validates missing and invalid configuration fields.
func (c *Config) Validate() error {
	if c.Metrics.Endpoint == "" {
		return errNoEndpoint
	}

	var errs error
	if c.Metrics.TLS != nil {
		if c.Metrics.TLS.CertFile == "" {
			errs = multierr.Append(errs, errNoCert)
		}

		if c.Metrics.TLS.KeyFile == "" {
			errs = multierr.Append(errs, errNoKey)
		}
	}

	if c.Metrics.BasicAuth != nil {
		if c.Metrics.BasicAuth.Username == "" {
			errs = multierr.Append(errs, errNoUsername)
		}

		if c.Metrics.BasicAuth.Password == "" {
			errs = multierr.Append(errs, errNoPass)
		}
	}

	_, _, err := net.SplitHostPort(c.Metrics.Endpoint)
	if err != nil {
		errs = multierr.Append(errs, fmt.Errorf("failed to split endpoint into 'host:port' pair: %w", err))
	}

	return errs
}
