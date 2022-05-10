// Package observiq contains OpAmp structures compatible with the observiq client
package observiq

import (
	"context"
	"errors"
	"net/url"

	"github.com/observiq/observiq-otel-collector/opamp"
	"github.com/open-telemetry/opamp-go/client"
	"github.com/open-telemetry/opamp-go/protobufs"
	"go.uber.org/zap"
)

var (
	// ErrUnsupportedURL is error returned when creating a client with an unsupported URL scheme
	ErrUnsupportedURL = errors.New("unsupported URL")
)

// Client represents a client that is connected to Iris via OpAmp
type Client struct {
	opampClient client.OpAMPClient
	logger      *zap.SugaredLogger
}

// NewClient creates a new OpAmp client
func NewClient(defaultLogger *zap.SugaredLogger, config opamp.Config) (opamp.Client, error) {
	clientLogger := defaultLogger.Named("opamp")

	// Parse URL to determin scheme
	opampURL, err := url.Parse(config.Endpoint)
	if err != nil {
		return nil, err
	}

	switch opampURL.Scheme {
	case "http", "https": // TODO might not be supported
		return &Client{
			opampClient: client.NewHTTP(clientLogger),
			logger:      clientLogger,
		}, nil
	case "wss":
		return &Client{
			opampClient: client.NewWebSocket(clientLogger),
			logger:      clientLogger,
		}, nil
	default:
		return nil, ErrUnsupportedURL
	}
}

// Connect initiates a connection to the OpAmp server based on the supplied configuration
func (o *Client) Connect(ctx context.Context, config opamp.Config) error {
	settings := client.StartSettings{
		OpAMPServerURL:                    config.Endpoint,
		AuthorizationHeader:               *config.SecretKey,
		TLSConfig:                         nil, // TODO add support for TLS
		InstanceUid:                       config.AgentID,
		AgentDescription:                  &protobufs.AgentDescription{},
		Callbacks:                         nil,
		LastRemoteConfigHash:              []byte{},
		LastEffectiveConfig:               &protobufs.EffectiveConfig{},
		LastConnectionSettingsHash:        []byte{},
		LastServerProvidedAllPackagesHash: []byte{},
	}

	return nil
}
