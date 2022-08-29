package report

import (
	"net/http"
)

// Client represents a client that can report information to a platform
type Client interface {
	// Do makes the specified request and return the body contents.
	// An error will be returned if an error occurred or non-200 code was returned
	Do(req *http.Request) (*http.Response, error)
}
