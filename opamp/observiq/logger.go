package observiq

import (
	"context"

	"github.com/open-telemetry/opamp-go/client/types"
	"go.uber.org/zap"
)

type zapOpAMPLoggerAdapter struct {
	logger *zap.SugaredLogger
}

var _ types.Logger = (*zapOpAMPLoggerAdapter)(nil)

func newZapOpAMPLoggerAdapter(logger *zap.Logger) *zapOpAMPLoggerAdapter {
	return &zapOpAMPLoggerAdapter{
		logger: logger.Sugar(),
	}
}

func (o zapOpAMPLoggerAdapter) Debugf(_ context.Context, format string, v ...any) {
	o.logger.Debugf(format, v...)
}

func (o zapOpAMPLoggerAdapter) Errorf(_ context.Context, format string, v ...any) {
	o.logger.Errorf(format, v...)
}
