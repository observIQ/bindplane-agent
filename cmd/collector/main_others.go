//go:build !windows
// +build !windows

package main

import (
	"context"

	"go.opentelemetry.io/collector/service"
)

func run(ctx context.Context, params service.CollectorSettings) error {
	return runInteractive(ctx, params)
}
