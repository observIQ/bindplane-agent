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
	"os/signal"
	"syscall"
	"time"

	"github.com/observiq/observiq-otel-collector/collector"
	"github.com/observiq/observiq-otel-collector/internal/version"
	"github.com/observiq/observiq-otel-collector/opamp"
	"github.com/observiq/observiq-otel-collector/opamp/observiq"
	"github.com/spf13/pflag"
	"go.opentelemetry.io/collector/service"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultLogger *zap.Logger

func main() {
	setupDefaultLogger()

	configPaths := pflag.StringSlice("config", []string{"./config.yaml"}, "the collector config path")
	managerConfigPath := pflag.String("manager", "./manager.yaml", "The configuration for remote management")
	_ = pflag.String("log-level", "", "not implemented") // TEMP(jsirianni): Required for OTEL k8s operator
	var showVersion = pflag.BoolP("version", "v", false, "prints the version of the collector")
	pflag.Parse()

	if *showVersion {
		fmt.Println("observiq-otel-collector version", version.Version())
		fmt.Println("commit:", version.GitHash())
		fmt.Println("built at:", version.Date())
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	settings, err := collector.NewSettings(*configPaths, version.Version(), nil)
	if err != nil {
		defaultLogger.Fatal("Settings configuration failed", zap.Error(err))
	}

	// See if manager config file exists. If so run in remote managed mode otherwise standalone mode
	// TODO(cpheps) clean this up in follow up work
	if _, err := os.Stat(*managerConfigPath); err == nil {
		log.Println("Starting Management Path")
		opampConfig, err := opamp.ParseConfig(*managerConfigPath)
		if err != nil {
			defaultLogger.Fatal("Failed to parse manager config", zap.Error(err))
		}

		// Create collector
		currCollector := collector.New((*configPaths)[0], version.Version(), []zap.Option{})

		// Create a blank log file as iris needs this currently
		logFilePath := "./logging.yaml"
		_, err = os.Create(logFilePath)
		if err != nil {
			defaultLogger.Fatal("Failed to create logging.yaml. Please create an empty logging.yaml file next to the collector if you wish to proceed in managed mode", zap.Error(err))
		}

		// Create client Args
		clientArgs := &observiq.NewClientArgs{
			DefaultLogger:       defaultLogger.Sugar(),
			Config:              *opampConfig,
			Collector:           currCollector,
			ManagerConfigPath:   *managerConfigPath,
			CollectorConfigPath: (*configPaths)[0],
			LoggerConfigPath:    logFilePath, // temporary as iris needs a logging file
		}

		if err := runRemoteManaged(ctx, clientArgs); err != nil {
			defaultLogger.Fatal("Remote Managment failed", zap.Error(err))
		}
	} else if errors.Is(err, os.ErrNotExist) {
		// Run standalone
		if err := run(ctx, *settings); err != nil {
			defaultLogger.Fatal("Collector error", zap.Error(err))
		}
	} else {
		defaultLogger.Fatal("Error while searching for management config", zap.Error(err))
	}

}

func runInteractive(ctx context.Context, params service.CollectorSettings) error {
	svc, err := service.New(params)
	if err != nil {
		return fmt.Errorf("failed to create new service: %w", err)
	}

	if err := svc.Run(ctx); err != nil {
		return fmt.Errorf("collector server run finished with error: %w", err)
	}

	return nil
}

func runRemoteManaged(ctx context.Context, clientArgs *observiq.NewClientArgs) error {
	// Create new client
	client, err := observiq.NewClient(clientArgs)
	if err != nil {
		return err
	}

	// Connect to manager platform
	if err := client.Connect(ctx, clientArgs.Config.Endpoint, clientArgs.Config.GetSecretKey()); err != nil {
		return err
	}

	// Wait for close signal
	<-ctx.Done()
	defaultLogger.Info("Signal received")

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
