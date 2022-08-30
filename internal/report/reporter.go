package report

import "context"

// ReporterKind represents the set of reporters available
type ReporterKind string

// Reporter represents a a structure to collector and report specific structures
type Reporter interface {
	// Type returns the type of this reporter
	Type() ReporterKind

	// Report starts reporting with the passed in configuration.
	Report(config any) error

	// Stop stops the reporter
	Stop(context.Context) error
}
