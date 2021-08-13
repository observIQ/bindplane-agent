package main

import (
	"log"
	"os"

	"github.com/observiq/observiq-collector/collector"
	"github.com/observiq/observiq-collector/internal/context"
	"github.com/observiq/observiq-collector/internal/env"
	"github.com/observiq/observiq-collector/internal/logging"
	"github.com/observiq/observiq-collector/internal/version"
	"github.com/observiq/observiq-collector/manager"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

func main() {
	// TODO: Revist default values for flags
	var managerConfigPath = pflag.String("manager-config", "./remote.yaml", "the manager config path")
	var collectorConfigPath = pflag.String("collector-config", "./config.yaml", "the collector config path")
	var loggingConfigPath = pflag.String("logging-config", "", "the logging config path")
	pflag.Parse()

	loggingConfig := getLoggingConfig(*loggingConfigPath)
	loggingOpts := logging.GetCollectorLoggingOpts(loggingConfig)
	logger := logging.GetManagerLogger(loggingConfig)

	logger.Info("Starting observIQ Agent", zap.String("version", version.Version),
		zap.String("git-hash", version.GitHash), zap.String("date-compiled", version.Date))

	collector := collector.New(*collectorConfigPath, loggingOpts)
	managerConfig, err := manager.ReadConfig(*managerConfigPath)
	if err != nil {
		log.Fatalf("Failed to read manager config: %s", err)
	}

	ctx := context.EmptyContext()
	if ppid := env.GetLauncherID(); ppid != 0 {
		ctx = context.WithParent(ppid)
	}

	managerCtx, cancel := context.WithInterrupt(ctx)
	defer cancel()

	manager := manager.New(managerConfig, collector, logger)
	exitCode := manager.Run(managerCtx)
	// nolint
	os.Exit(exitCode)
}

// getLoggingConfig loads the configuration from the file path given.
//  If the config cannot be loaded, a default config is used instead.
//  If an empty string is passed as the path (no logging config specified),
//  then a "nil" config is returned.
func getLoggingConfig(logConfigPath string) *logging.Config {
	if logConfigPath == "" {
		return nil
	}

	c, err := logging.LoadConfig(logConfigPath)
	if err != nil {
		c = logging.DefaultConfig()
	}

	return c
}
