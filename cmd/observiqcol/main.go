package main

import (
	"context"
	"log"

	"github.com/observiq/observiq-collector/collector"
	"github.com/observiq/observiq-collector/internal/env"
	"github.com/observiq/observiq-collector/internal/logging"
	"github.com/observiq/observiq-collector/manager"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

func main() {
	// TODO: Revist default values for flags
	var managerConfigPath = pflag.String("manager-config", "./remote.yaml", "the manager config path")
	var collectorConfigPath = pflag.String("collector-config", "./config.yaml", "the collector config path")
	pflag.Parse()

	loggingOpts := getLoggingOpts()
	logger := getLogger()

	settings, err := collector.NewSettings(*collectorConfigPath, loggingOpts)
	if err != nil {
		log.Fatalf("Failed to get collector settings: %s", err)
	}
	collector := collector.New(settings)

	config, err := manager.ConfigFromFile(*managerConfigPath)
	if err != nil {
		log.Fatalf("Failed to get manager config: %s", err)
	}
	manager := manager.New(config, collector, logger)

	// TODO: Look into handling interupt signals with context
	if err := manager.Run(context.Background()); err != nil {
		log.Fatalf("Manager failed: %s", err)
	}
}

// TODO: Revisit logging to determine appropriate configuration and panic behavior
func getLoggingOpts() []zap.Option {
	var loggingOpts []zap.Option
	if env.IsFileLoggingEnabled() {
		if fp, ok := env.GetLoggingPath(); ok {
			loggingOpts = []zap.Option{logging.FileLoggingCoreOption(fp)}
		} else {
			panic("Failed to find file path for logs, is OIQ_COLLECTOR_HOME set?")
		}
	}
	return loggingOpts
}

// TODO: Determine our logging strategy
func getLogger() *zap.Logger {
	zapConfig := zap.NewProductionConfig()
	zapConfig.OutputPaths = []string{"stdout"}
	zapLogger, _ := zapConfig.Build()
	return zapLogger
}
