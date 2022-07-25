package logging

import (
	"fmt"

	"github.com/observiq/observiq-otel-collector/updater/internal/path"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger(installDir string) (*zap.Logger, error) {
	prodConf := zap.NewProductionConfig()

	prodLogger, err := prodConf.Build(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return core(installDir)
	}))

	if err != nil {
		return nil, fmt.Errorf("failed to make logger: %w", err)
	}

	return prodLogger, nil
}

func core(installDir string) zapcore.Core {
	logger := &lumberjack.Logger{
		Filename:   path.LogFile(installDir),
		MaxSize:    10,
		MaxBackups: 3,
	}

	return zapcore.NewCore(encoder(), zapcore.AddSync(logger), zapcore.DebugLevel)
}

func encoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}
