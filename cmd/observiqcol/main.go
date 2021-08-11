package main

import (
	"log"

	"github.com/observiq/observiq-collector/collector"
	"github.com/observiq/observiq-collector/internal/context"
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
	logger, err := newLogger(loggingOpts)
	if err != nil {
		log.Fatalf("Failed to create logger: %s", err)
	}

	collector := collector.New(*collectorConfigPath, loggingOpts)
	managerConfig, err := manager.ReadConfig(*managerConfigPath)
	if err != nil {
		log.Fatalf("Failed to read manager config: %s", err)
	}

	ctx, cancel := context.WithInterrupt()
	defer cancel()

	manager := manager.New(managerConfig, collector, logger)
	if err := manager.Run(ctx); err != nil {
		log.Fatalf("Manager exited with error: %s", err)
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

// newLogger creates a new logger for the manager.
func newLogger(opts []zap.Option) (*zap.Logger, error) {
	zapConfig := zap.NewProductionConfig()
	zapConfig.OutputPaths = []string{"stdout"}
	return zapConfig.Build(opts...)
}
