//go:build tools

package tools

// This file is to keep tool dependencies, so they can be versioned and easily installed.
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
// nolint
import (
	_ "github.com/client9/misspell/cmd/misspell"
	_ "github.com/goreleaser/goreleaser/cmd"
	_ "github.com/sigstore/cosign/cmd/cosign"
	_ "github.com/uw-labs/lichen"
	_ "golang.org/x/tools/cmd/goimports"
)
