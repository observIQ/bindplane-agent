package main

import (
	"log"
	"os"

	"github.com/observiq/observiq-collector/collector"
	"github.com/observiq/observiq-collector/internal/context"
	"github.com/observiq/observiq-collector/internal/env"
	"github.com/observiq/observiq-collector/internal/logging"
	"github.com/observiq/observiq-collector/internal/migration"
	"github.com/observiq/observiq-collector/manager"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

func main() {
	// TODO: Revist default values for flags
	var managerConfigPath = pflag.String("manager-config", env.DefaultRemoteConfigFile(), "the manager config path")
	var collectorConfigPath = pflag.String("collector-config", "./config.yaml", "the collector config path")
	var loggingConfigPath = pflag.String("logging-config", env.DefaultLoggingConfigFile(), "the logging config path")
	var tryMigration = pflag.BoolP("migrate", "m", false, "try migrating configs from BPAgent on startup")
	pflag.Parse()

	loggingConfig := getLoggingConfig(*loggingConfigPath)
	loggingOpts := logging.GetCollectorLoggingOpts(loggingConfig)
	logger := logging.GetManagerLogger(loggingConfig)

	if *tryMigration {
		migrationLogger := logger.Named("migration")
		shouldMigrate, err := migration.ShouldMigrate()
		switch {
		case err != nil:
			migrationLogger.Error("Skipping config migration, encountered error when looking for BPAgent install", zap.Error(err))
		case shouldMigrate:
			migrationLogger.Info("Detected BPAgent install, attempting config migration")
			err := migration.Migrate()
			if err != nil {
				migrationLogger.Panic("Failed to migrate!", zap.Error(err))
			}
		default:
			migrationLogger.Info("Skipping config migration, no install of BPAgent detected")
		}
	}

	collector := collector.New(*collectorConfigPath, loggingOpts)
	managerConfig, err := manager.ReadConfig(*managerConfigPath)
	if err != nil {
		log.Fatalf("Failed to read manager config: %s", err)
	}

	ctx := context.EmptyContext()
	if ppid := env.GetLauncherPPID(); ppid != 0 {
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

	c, loadErr := logging.LoadConfig(logConfigPath)
	if loadErr != nil {
		c = logging.DefaultConfig()
	}

	return c
}
