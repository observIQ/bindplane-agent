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

// Package opamp contains configurations and protocol implementations to handle OpAmp communication.
package opamp

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var (
	// errPrefixReadFile for error when reading config file
	errPrefixReadFile = "failed to read OpAmp config file"

	// errPrefixParse for error when parsing config
	errPrefixParse = "failed to parse OpAmp config"

	// errMissingTLSFiles for missing a required tls file
	errMissingTLSFiles = "must specify both a key and cert file for TLS"

	// errInvalidKeyFile for key file that is not readable
	errInvalidKeyFile = "failed to read TLS key file"

	// errInvalidCertFile for cert file that is not readable
	errInvalidCertFile = "failed to read TLS cert file"

	// errInvalidCAFile for ca file that is not readable
	errInvalidCAFile = "failed to read TLS CA file"
)

// Config contains the configuration for the collector to communicate with an OpAmp enabled platform.
type Config struct {
	Endpoint  string     `yaml:"endpoint"`
	SecretKey *string    `yaml:"secret_key,omitempty"`
	AgentID   string     `yaml:"agent_id"`
	TLS       *TLSConfig `yaml:"tls_config,omitempty"`

	// Updatable fields
	Labels    *string `yaml:"labels,omitempty"`
	AgentName *string `yaml:"agent_name,omitempty"`
}

// TLSConfig represents the TLS config to connect to OpAmp server
type TLSConfig struct {
	insecure bool
	KeyFile  *string `yaml:"key_file"`
	CertFile *string `yaml:"cert_file"`
	CAFile   *string `yaml:"ca_file"`
}

// ToTLS converts the config to a tls.Config
func (c Config) ToTLS() (*tls.Config, error) {
	if c.TLS == nil {
		return nil, nil
	}

	if c.TLS.insecure {
		return &tls.Config{
			InsecureSkipVerify: true,
		}, nil
	}

	cert, err := tls.LoadX509KeyPair(*c.TLS.CertFile, *c.TLS.KeyFile)
	if err != nil {
		return nil, errors.New(errMissingTLSFiles)
	}

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	// Load CA cert
	if c.TLS.CAFile != nil {
		caCert, err := ioutil.ReadFile(*c.TLS.CAFile)
		if err != nil {
			return nil, errors.New(errInvalidCAFile)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig.RootCAs = caCertPool
	}

	tlsConfig.BuildNameToCertificate()

	return tlsConfig, nil
}

// ParseConfig given a configuration file location will parse the config
func ParseConfig(configLocation string) (*Config, error) {
	configPath := filepath.Clean(configLocation)

	// Read in config file contents
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errPrefixReadFile, err)
	}

	// Parse the config
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("%s: %w", errPrefixParse, err)
	}

	if config.TLS != nil && config.TLS.insecure == false {
		if config.TLS.CertFile == nil || config.TLS.KeyFile == nil {
			return nil, errors.New(errMissingTLSFiles)
		}

		if _, err := os.Stat(*config.TLS.KeyFile); errors.Is(err, os.ErrNotExist) {
			return nil, errors.New(errInvalidKeyFile)
		}

		if _, err := os.Stat(*config.TLS.CertFile); errors.Is(err, os.ErrNotExist) {
			return nil, errors.New(errInvalidCertFile)
		}

		if config.TLS.CAFile != nil {
			if _, err := os.Stat(*config.TLS.CAFile); errors.Is(err, os.ErrNotExist) {
				return nil, errors.New(errInvalidCAFile)
			}
		}
	}
	return &config, nil
}

// Copy creates a deep copy of this config
func (c Config) Copy() *Config {

	cfgCopy := &Config{
		Endpoint: c.Endpoint,
		AgentID:  c.AgentID,
	}

	if c.SecretKey != nil {
		cfgCopy.SecretKey = new(string)
		*cfgCopy.SecretKey = *c.SecretKey
	}
	if c.Labels != nil {
		cfgCopy.Labels = new(string)
		*cfgCopy.Labels = *c.Labels
	}
	if c.AgentName != nil {
		cfgCopy.AgentName = new(string)
		*cfgCopy.AgentName = *c.AgentName
	}
	if c.TLS != nil {
		cfgCopy.TLS = c.TLS.copy()
	}

	return cfgCopy
}

func (t TLSConfig) copy() *TLSConfig {
	tlsCopy := TLSConfig{
		insecure: t.insecure,
	}

	if t.CertFile != nil {
		tlsCopy.CertFile = new(string)
		*tlsCopy.CertFile = *t.CertFile
	}
	if t.KeyFile != nil {
		tlsCopy.KeyFile = new(string)
		*tlsCopy.KeyFile = *t.KeyFile
	}
	if t.CAFile != nil {
		tlsCopy.CAFile = new(string)
		*tlsCopy.CAFile = *t.CAFile
	}

	return &tlsCopy
}

// GetSecretKey returns secret key if set else returns empty string
func (c Config) GetSecretKey() string {
	if c.SecretKey == nil {
		return ""
	}

	return *c.SecretKey
}

// CmpUpdatableFields compares updatable fields for equality
func (c Config) CmpUpdatableFields(o Config) (equal bool) {
	if !cmpStringPtr(c.AgentName, o.AgentName) {
		return false
	}

	return cmpStringPtr(c.Labels, o.Labels)
}

func cmpStringPtr(p1, p2 *string) bool {
	switch {
	case p1 == nil && p2 == nil:
		return true
	case p1 == nil && p2 != nil:
		fallthrough
	case p1 != nil && p2 == nil:
		fallthrough
	case *p1 != *p2:
		return false
	default:
		return true
	}
}
