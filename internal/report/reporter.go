package report

import "context"

// ReporterType represents the set of reporters available
type ReporterType string

// Reporter represents a a structure to collector and report specific structures
type Reporter interface {
	// Type returns the type of this reporter
	Type() ReporterType

	// IsEnabled signals if this reporter is enabled
	IsEnabled() bool

	// Start kicks off the reporter sending to the client
	Start(Client) error

	// Stop stops the reporter
	Stop(context.Context) error
}
