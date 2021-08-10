package main

import (
	"log"
	"os"

	"github.com/observiq/observiq-collector/collector"
	"github.com/observiq/observiq-collector/internal/context"
	"github.com/observiq/observiq-collector/internal/env"
	"github.com/observiq/observiq-collector/internal/logging"
	"github.com/observiq/observiq-collector/manager"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// TODO: Revist default values for flags
	var managerConfigPath = pflag.String("manager-config", "./remote.yaml", "the manager config path")
	var collectorConfigPath = pflag.String("collector-config", "./config.yaml", "the collector config path")
	var loggingConfigPath = pflag.String("logging-config", "", "the logging config path")
	pflag.Parse()

	loggingConfig := getLoggingConfig(*loggingConfigPath)
	loggingOpts := getCollectorLoggingOpts(loggingConfig)
	logger := getManagerLogger(loggingConfig)

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

	c, err := logging.LoadConfig(logConfigPath)
	if err != nil {
		// TODO: Should we panic here instead of using defaults?
		c = logging.DefaultConfig()
	}

	return c
}

// getCollectorLoggingOpts gets the logging options passed to the collector
func getCollectorLoggingOpts(config *logging.Config) []zap.Option {
	var loggingOpts []zap.Option
	if config != nil {
		loggingOpts = append(loggingOpts, zap.WrapCore(func(c zapcore.Core) zapcore.Core {
			return logging.CreateFileCore(&config.Collector)
		}))
	}
	return loggingOpts
}

// getManagerLogger creates the logger for the manager.
//  this may be a file logger, or a console logger, depending of env.FileLoggingEnabled()
func getManagerLogger(c *logging.Config) *zap.Logger {
	var zapLogger *zap.Logger

	if c != nil {
		zapLogger = zap.New(logging.CreateFileCore(&c.Manager))
	} else {
		zapConfig := zap.NewProductionConfig()
		zapConfig.OutputPaths = []string{"stdout"}

		var err error
		zapLogger, err = zapConfig.Build()
		if err != nil {
			// TODO: Evaluate panic-ing here
			panic("Failed to create stdout logger")
		}

	}

	return zapLogger
}
