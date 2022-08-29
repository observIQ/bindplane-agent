package report

import "context"

// ReporterKind represents the set of reporters available
type ReporterKind string

// Reporter represents a a structure to collector and report specific structures
type Reporter interface {
	// Type returns the type of this reporter
	Type() ReporterKind

	// ApplyConfig applies a new configuration for the reporter
	ApplyConfig(any) error

	// Start kicks off the reporter.
	// If this starts a goroutine it should be terminated by calling Stop
	Start() error

	// Stop stops the reporter
	Stop(context.Context) error
}
