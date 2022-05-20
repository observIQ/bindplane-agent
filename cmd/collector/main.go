// Copyright  observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/observiq/observiq-otel-collector/collector"
	"github.com/observiq/observiq-otel-collector/internal/logging"
	"github.com/observiq/observiq-otel-collector/internal/service"
	"github.com/observiq/observiq-otel-collector/internal/version"
	"github.com/observiq/observiq-otel-collector/opamp"
	"github.com/observiq/observiq-otel-collector/opamp/observiq"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultLogger *zap.Logger

func main() {
	setupDefaultLogger()

	configPaths := pflag.StringSlice("config", []string{"./config.yaml"}, "the collector config path")
	managerConfigPath := pflag.String("manager", "./manager.yaml", "The configuration for remote management")
	loggingPath := pflag.String("logging", "./logging.yaml", "the collector logging config path")

	_ = pflag.String("log-level", "", "not implemented") // TEMP(jsirianni): Required for OTEL k8s operator
	var showVersion = pflag.BoolP("version", "v", false, "prints the version of the collector")
	pflag.Parse()

	if *showVersion {
		fmt.Println("observiq-otel-collector version", version.Version())
		fmt.Println("commit:", version.GitHash())
		fmt.Println("built at:", version.Date())
		return
	}

	logOpts, err := logOptions(loggingPath)
	if err != nil {
		log.Fatalf("Failed to get log options: %v", err)
	}

	// logOpts will override options here
	logger, err := zap.NewProduction(logOpts...)
	if err != nil {
		defaultLogger.Fatal("Settings configuration failed", zap.Error(err))
	}

	var runnableService service.RunnableService

	col := collector.New(*configPaths, version.Version(), logOpts)

	// See if manager config file exists. If so run in remote managed mode otherwise standalone mode
	// TODO(cpheps) clean this up in follow up work
	if _, err := os.Stat(*managerConfigPath); err == nil {
		log.Println("Starting Management Path")
		opampConfig, err := opamp.ParseConfig(*managerConfigPath)
		if err != nil {
			defaultLogger.Fatal("Failed to parse manager config", zap.Error(err))
		}

		// Create client Args
		clientArgs := &observiq.NewClientArgs{
			DefaultLogger:       defaultLogger.Sugar(),
			Config:              *opampConfig,
			Collector:           col,
			ManagerConfigPath:   *managerConfigPath,
			CollectorConfigPath: (*configPaths)[0],
			LoggerConfigPath:    *loggingPath, // temporary as iris needs a logging file
		}

		if err := runRemoteManaged(context.Background(), clientArgs); err != nil {
			defaultLogger.Fatal("Remote Management failed", zap.Error(err))
		}
	} else if errors.Is(err, os.ErrNotExist) {

		runnableService = service.NewStandaloneCollectorService(col)
	} else {
		defaultLogger.Fatal("Error while searching for management config", zap.Error(err))
		log.Fatalf("Failed to set up logger: %v", err)
	}

	err = service.RunService(logger, runnableService)
	if err != nil {
		logger.Fatal("RunService returned error", zap.Error(err))
	}

}

func logOptions(loggingConfigPath *string) ([]zap.Option, error) {
	if loggingConfigPath == nil {
		return nil, nil
	}

	l, err := logging.NewLoggerConfig(*loggingConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger config: %w", err)
	}

	return l.Options()
}

func runRemoteManaged(ctx context.Context, clientArgs *observiq.NewClientArgs) error {
	// Create new client
	client, err := observiq.NewClient(clientArgs)
	if err != nil {
		return err
	}

	// Connect to manager platform
	if err := client.Connect(ctx); err != nil {
		return err
	}

	// Wait for close signal
	<-ctx.Done()
	defaultLogger.Info("Exit signal received shutting down collector")

	// Disconnect from opamp
	waitCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Disconnect(waitCtx); err != nil {
		return fmt.Errorf("error when client disconnect: %w", err)
	}

	return nil
}

func setupDefaultLogger() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	syncer := zapcore.Lock(os.Stdout)
	core := zapcore.NewCore(encoder, syncer, zapcore.DebugLevel)
	defaultLogger = zap.New(core)
}
