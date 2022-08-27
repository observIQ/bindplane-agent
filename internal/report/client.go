package report

import "context"

// Client represents a client that can report information to a platform
type Client interface {

	// ReportSnapShot reports the payload associated with the component ID
	ReportSnapShot(ctx context.Context, componentID string, payload []byte) error
}
