package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// CreateFileCore creates a new core for file logging based on the provided config
func CreateFileCore(config *LoggerConfig) zapcore.Core {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   config.File,
		MaxBackups: config.MaxBackups,
		MaxSize:    config.MaxMegabytes,
		MaxAge:     config.MaxDays,
	})

	zapCfg := zap.NewProductionEncoderConfig()
	zapCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	wc := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zapCfg),
		w,
		config.Level,
	)

	return wc
}
