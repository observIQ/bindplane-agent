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

// Package observiq contains OpAmp structures compatible with the observiq client
package observiq

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/observiq/observiq-otel-collector/opamp"
	"github.com/open-telemetry/opamp-go/client"
	"github.com/open-telemetry/opamp-go/client/types"
	"github.com/open-telemetry/opamp-go/protobufs"
	"go.uber.org/zap"
)

var (
	// ErrUnsupportedURL is error returned when creating a client with an unsupported URL scheme
	ErrUnsupportedURL = errors.New("unsupported URL")
)

// Ensure interface is satisfied
var _ opamp.Client = (*Client)(nil)

// Client represents a client that is connected to Iris via OpAmp
type Client struct {
	opampClient   client.OpAMPClient
	logger        *zap.SugaredLogger
	ident         *identity
	configManager opamp.ConfigManager

	// Channel that when closed signals to the listeners that a shutdown is necessary
	// NOTE: This currently is a stopgap until we can hot reload configurations
	shutdownChan chan struct{}
}

// NewClient creates a new OpAmp client
// The passed in configmanager should be preloaded with all known configs.
// The shutdown channel will be closed when a reconfiguration occurs an a process shutdown is required.
func NewClient(defaultLogger *zap.SugaredLogger, config opamp.Config, configManager opamp.ConfigManager, shutdownChan chan struct{}) (opamp.Client, error) {
	clientLogger := defaultLogger.Named("opamp")

	observiqClient := &Client{
		logger:        clientLogger,
		ident:         newIdentity(clientLogger, config),
		configManager: configManager,
		shutdownChan:  shutdownChan,
	}

	// Parse URL to determin scheme
	opampURL, err := url.Parse(config.Endpoint)
	if err != nil {
		return nil, err
	}

	// Create collect client based on URL scheme
	switch opampURL.Scheme {
	case "ws", "wss":
		observiqClient.opampClient = client.NewWebSocket(clientLogger)
	default:
		return nil, ErrUnsupportedURL
	}

	return observiqClient, nil
}

// Connect initiates a connection to the OpAmp server based on the supplied configuration
func (c *Client) Connect(ctx context.Context, config opamp.Config) error {
	// Compose and set the agent description
	if err := c.opampClient.SetAgentDescription(c.ident.ToAgentDescription()); err != nil {
		c.logger.Error("Error while setting agent description", zap.Error(err))
		return err
	}

	settings := types.StartSettings{
		OpAMPServerURL:      config.Endpoint,
		AuthorizationHeader: *config.SecretKey,
		TLSConfig:           nil, // TODO add support for TLS
		InstanceUid:         config.AgentID,
		Callbacks: types.CallbacksStruct{
			OnConnectFunc:          c.onConnectHandler,
			OnConnectFailedFunc:    c.onConnectFailedHandler,
			OnErrorFunc:            c.onErrorHandler,
			OnRemoteConfigFunc:     c.onRemoteConfigHandler,
			GetEffectiveConfigFunc: c.onGetEffectiveConfigHandler,
			// Unimplemented Handles
			// SaveRemoteConfigStatusFunc:
			// OnOpampConnectionSettingsFunc
			// OnOpampConnectionSettingsAcceptedFunc
			// OnOwnTelemetryConnectionSettingsFunc
			// OnOtherConnectionSettingsFunc
			// OnPackagesAvailableFunc
			// OnAgentIdentificationFunc
			// OnCommandFunc
		},
	}

	return c.opampClient.Start(ctx, settings)
}

// Disconnect disconnects from the server
func (c *Client) Disconnect(ctx context.Context) error {
	return c.opampClient.Stop(ctx)
}

// client callbacks

func (c *Client) onConnectHandler() {
	c.logger.Info("Successfully connected to server")

	if err := c.opampClient.SetAgentDescription(c.ident.ToAgentDescription()); err != nil {
		c.logger.Error("Failed to set agent description", zap.Error(err))
	}
}

func (c *Client) onConnectFailedHandler(err error) {
	c.logger.Error("Failed to connect to server", zap.Error(err))
}

func (c *Client) onErrorHandler(errResp *protobufs.ServerErrorResponse) {
	c.logger.Error("Server returned an error response", zap.String("Error", errResp.GetErrorMessage()))
}

func (c *Client) onRemoteConfigHandler(_ context.Context, remoteConfig *protobufs.AgentRemoteConfig) (*protobufs.EffectiveConfig, bool, error) {
	c.logger.Debug("Remote config handler")

	effectiveConfig, changed, err := c.configManager.ApplyConfigChanges(remoteConfig)
	if err != nil {
		c.logger.Error("Failed applying remote config", zap.Error(err))
		return nil, changed, fmt.Errorf("Failed to apply config changes: %w", err)
	}

	// Since we can't hot reload configs we need to restart in order to take advantage of the new configurations.
	// We will exit and trust the service manager wrapping us to restart
	if changed {
		// Close shutdown channel to signal restart
		close(c.shutdownChan)
	}

	return effectiveConfig, false, nil
}

func (c *Client) onGetEffectiveConfigHandler(_ context.Context) (*protobufs.EffectiveConfig, error) {
	c.logger.Debug("Remote Compose Effective config handler")
	return c.configManager.ComposeEffectiveConfig()
}
