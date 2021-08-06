package main

import (
	"context"
	"log"
	"os"

	"github.com/observiq/observiq-collector/collector"
	"github.com/observiq/observiq-collector/internal/env"
	"github.com/observiq/observiq-collector/internal/logging"
	"github.com/observiq/observiq-collector/manager"
	"go.uber.org/zap"
)

func main() {
	managerConfigPath, ok := os.LookupEnv("OBSERVIQ_MANAGER_CONFIG")
	if !ok {
		log.Fatalf("OBSERVIQ_MANAGER_CONFIG environment variable is not set")
	}

	collectorConfigPath, ok := os.LookupEnv("OBSERVIQ_COLLECTOR_CONFIG")
	if !ok {
		log.Fatalf("OBSERVIQ_COLLECTOR_CONFIG environment variable is not set")
	}

	loggingOpts := getLoggingOpts()
	logger := getLogger()

	settings, err := collector.NewSettings(collectorConfigPath, loggingOpts)
	if err != nil {
		log.Fatalf("Failed to get collector settings: %s", err)
	}
	collector := collector.New(settings)

	config, err := manager.ConfigFromFile(managerConfigPath)
	if err != nil {
		log.Fatalf("Failed to get manager config: %s", err)
	}
	manager := manager.New(config, collector, logger)

	if err := run(manager); err != nil {
		log.Fatal(err)
	}
}

func runInteractive(manager *manager.Manager) error {
	if err := manager.Run(context.Background()); err != nil {
		return err
	}

	return nil
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
