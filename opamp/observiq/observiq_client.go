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
}

// NewClient creates a new OpAmp client
func NewClient(defaultLogger *zap.SugaredLogger, configManager opamp.ConfigManager, config opamp.Config) (opamp.Client, error) {
	clientLogger := defaultLogger.Named("opamp")

	observiqClient := &Client{
		logger:        clientLogger,
		ident:         NewIdentity(clientLogger, config),
		configManager: configManager,
	}

	// Parse URL to determin scheme
	opampURL, err := url.Parse(config.Endpoint)
	if err != nil {
		return nil, err
	}

	// Create collect client based on URL scheme
	switch opampURL.Scheme {
	case "wss":
		observiqClient.opampClient = client.NewWebSocket(clientLogger)
	default:
		return nil, ErrUnsupportedURL
	}

	return observiqClient, nil
}

// Connect initiates a connection to the OpAmp server based on the supplied configuration
func (c *Client) Connect(config opamp.Config) error {
	effectiveConfig, err := c.configManager.ComposeEffectiveConfig()
	if err != nil {
		c.logger.Errorf("Error while composing effective config", zap.Error(err))
		return fmt.Errorf("failed to compose effective config: %w", err)
	}

	settings := client.StartSettings{
		OpAMPServerURL:      config.Endpoint,
		AuthorizationHeader: *config.SecretKey,
		TLSConfig:           nil, // TODO add support for TLS
		InstanceUid:         config.AgentID,
		AgentDescription:    c.ident.ToAgentDescription(),
		Callbacks: types.CallbacksStruct{
			OnConnectFunc:       c.onConnectHandler,
			OnConnectFailedFunc: c.onConnectFailedHandler,
			OnErrorFunc:         c.onErrorHandler,
			OnRemoteConfigFunc:  c.onRemoteConfigHandler,
			// The below handlers are not currently implemented
			// OnOpampConnectionSettingsFunc
			// OnOpampConnectionSettingsAcceptedFunc
			// OnOwnTelemetryConnectionSettingsFunc
			// OnOtherConnectionSettingsFunc
			// OnPackagesAvailableFunc
			// OnAgentIdentificationFunc
			// OnCommandFunc
		},
		LastRemoteConfigHash: effectiveConfig.GetHash(),
		LastEffectiveConfig:  effectiveConfig,
	}

	return c.opampClient.Start(settings)
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
		// TODO restart and remove return
	}

	return effectiveConfig, false, nil
}
