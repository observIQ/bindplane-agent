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
	"errors"
	"fmt"
	"log"
	"os"
	_ "time/tzdata"

	"github.com/google/uuid"
	"github.com/observiq/observiq-otel-collector/collector"
	"github.com/observiq/observiq-otel-collector/internal/logging"
	"github.com/observiq/observiq-otel-collector/internal/service"
	"github.com/observiq/observiq-otel-collector/internal/version"
	"github.com/observiq/observiq-otel-collector/opamp"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

const (
	// env variable name constants
	endpoint  = "OPAMP_ENDPOINT"
	agentID   = "OPAMP_AGENT_ID"
	secretKey = "OPAMP_SECRET_KEY"
	labels    = "OPAMP_LABELS"
)

func main() {
	collectorConfigPaths := pflag.StringSlice("config", []string{"./config.yaml"}, "the collector config path")
	managerConfigPath := pflag.String("manager", "./manager.yaml", "The configuration for remote management")
	loggingConfigPath := pflag.String("logging", "./logging.yaml", "the collector logging config path")

	_ = pflag.String("log-level", "", "not implemented") // TEMP(jsirianni): Required for OTEL k8s operator
	var showVersion = pflag.BoolP("version", "v", false, "prints the version of the collector")
	pflag.Parse()

	if *showVersion {
		fmt.Println("observiq-otel-collector version", version.Version())
		fmt.Println("commit:", version.GitHash())
		fmt.Println("built at:", version.Date())
		return
	}

	logOpts, err := logOptions(loggingConfigPath)
	if err != nil {
		log.Fatalf("Failed to get log options: %v", err)
	}

	// logOpts will override options here
	logger, err := zap.NewProduction(logOpts...)
	if err != nil {
		log.Fatalf("Failed to set up logger: %v", err)
	}

	var runnableService service.RunnableService

	col := collector.New(*collectorConfigPaths, version.Version(), logOpts)

	// See if manager config file exists. If so run in remote managed mode otherwise standalone mode
	if err := checkManagerConfig(managerConfigPath); err == nil {
		logger.Info("Starting In Managed Mode")

		runnableService, err = service.NewManagedCollectorService(col, logger, *managerConfigPath, (*collectorConfigPaths)[0], *loggingConfigPath)
		if err != nil {
			logger.Fatal("Failed to initiate managed mode", zap.Error(err))
		}
	} else if errors.Is(err, os.ErrNotExist) {
		logger.Info("Starting Standalone Mode")
		runnableService = service.NewStandaloneCollectorService(col)
	} else {
		logger.Fatal("Error while searching for management config", zap.Error(err))
	}

	// Run service
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

func checkManagerConfig(configPath *string) error {
	if _, err := os.Stat(*configPath); err == nil {
		// manager config file exists, return and verify
		return err

	} else if errors.Is(err, os.ErrNotExist) {
		// manager.ymal file does *not* exist, create file using env variables
		newConfig := &opamp.Config{}
		if ep, ok := os.LookupEnv(endpoint); ok {
			newConfig.Endpoint = ep

			var ai string
			if ai, ok = os.LookupEnv(agentID); !ok {
				ai = uuid.New().String()
			}
			newConfig.AgentID = ai

			if sk, ok := os.LookupEnv(secretKey); ok {
				newConfig.SecretKey = &sk
			}

			if label, ok := os.LookupEnv(labels); ok {
				newConfig.Labels = &label
			}

			var data []byte
			if data, err = yaml.Marshal(newConfig); err != nil {
				panic(err)
			}
			// write data to a manager.yaml file, with 0777 file permission
			if err := os.WriteFile(*configPath, data, 0777); err != nil {
				panic(err)
			}
			return err
		}
		return os.ErrNotExist
	}
	return os.ErrInvalid
}
