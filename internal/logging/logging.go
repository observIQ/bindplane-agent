package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// GetCollectorLoggingOpts gets the logging options passed to the collector
func GetCollectorLoggingOpts(config *Config) []zap.Option {
	var loggingOpts []zap.Option
	if config != nil {
		loggingOpts = append(loggingOpts, zap.WrapCore(func(c zapcore.Core) zapcore.Core {
			return CreateFileCore(&config.Collector)
		}))
	}
	return loggingOpts
}

// GetManagerLogger creates the logger for the manager.
//  this may be a file logger, or a console logger, depending of env.FileLoggingEnabled()
func GetManagerLogger(c *Config) *zap.Logger {
	var zapLogger *zap.Logger

	if c != nil {
		zapLogger = zap.New(CreateFileCore(&c.Manager))
	} else {
		zapConfig := zap.NewProductionConfig()
		zapConfig.OutputPaths = []string{"stdout"}

		var err error
		zapLogger, err = zapConfig.Build()
		if err != nil {
			panic("Failed to create stdout logger")
		}

	}

	return zapLogger
}
