package opamp

import "context"

// Client implements a connection with OpAmp enabled server
type Client interface {

	// Connect initiates a connection to the OpAmp server based on the supplied configuration
	Connect(config Config) error

	// Disconnect disconnects from the server
	Disconnect(ctx context.Context) error
}
