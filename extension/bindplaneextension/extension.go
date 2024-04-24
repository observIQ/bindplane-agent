package bindplaneextension

import (
	"context"

	"go.opentelemetry.io/collector/component"
)

type bindplaneExtension struct{}

func (bindplaneExtension) Start(ctx context.Context, host component.Host) error {
	return nil
}

func (bindplaneExtension) Shutdown(ctx context.Context) error {
	return nil
}
