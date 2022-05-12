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
	loggingConfigPath := pflag.String("logging", "./logging.yaml", "The configuration for logging")
	_ = pflag.String("log-level", "", "not implemented") // TEMP(jsirianni): Required for OTEL k8s operator
	pflag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	settings, err := collector.NewSettings(*configPaths, version.Version(), nil)
	if err != nil {
		defaultLogger.Fatal("Settings configuration failed", zap.Error(err))
	}

	if _, err := os.Stat(*managerConfigPath); err == nil {
		log.Println("Starting Management Path")
		// Management config exists
		configManger := observiq.NewAgentConfigManager(defaultLogger.Sugar())
		configManger.AddConfig(opamp.CollectorConfigName, (*configPaths)[0], opamp.NewYamlValidator(make(map[string]interface{})))
		configManger.AddConfig(opamp.ManagerConfigName, *managerConfigPath, opamp.NewYamlValidator(new(opamp.Config)))
		configManger.AddConfig(opamp.LoggingConfigName, *loggingConfigPath, opamp.NoopValidator)

		currCollector := collector.New((*configPaths)[0], version.Version(), []zap.Option{})

		if err := runRemoteManaged(ctx, currCollector, configManger, *managerConfigPath); err != nil {
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

func runRemoteManaged(ctx context.Context, currCollector collector.Collector, configManager opamp.ConfigManager, managerConfigPath string) error {
	opampConfig, err := opamp.ParseConfig(managerConfigPath)
	if err != nil {
		return err
	}

	shutdownChan := make(chan struct{})
	client, err := observiq.NewClient(defaultLogger.Sugar(), *opampConfig, configManager, shutdownChan)
	if err != nil {
		return err
	}

	if err := currCollector.Run(ctx); err != nil {
		return fmt.Errorf("collector failed to start: %w", err)
	}

	if err := client.Connect(ctx, *opampConfig); err != nil {
		return err
	}

	// Wait for one shutdown or the other
	select {
	case <-ctx.Done():
		log.Println("Signal received")
	case <-shutdownChan:
		log.Println("Received shutdown from client")
	}

	// Stop collector
	currCollector.Stop()

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
