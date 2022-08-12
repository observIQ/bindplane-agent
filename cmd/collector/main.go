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

// Package main provides entry point for the collector
package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	endpointENV   = "OPAMP_ENDPOINT"
	agentIDENV    = "OPAMP_AGENT_ID"
	secretkeyENV  = "OPAMP_SECRET_KEY" //#nosec G101
	labelsENV     = "OPAMP_LABELS"
	agentNameENV  = "OPAMP_AGENT_NAME"
	configPathENV = "CONFIG_YAML_PATH"
)

func main() {
	collectorConfigPaths := pflag.StringSlice("config", getDefaultCollectorConfigPaths(), "the collector config path")
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

func getDefaultCollectorConfigPaths() []string {
	cp, ok := os.LookupEnv(configPathENV)
	if ok {
		return []string{cp}
	} else {
		return []string{"./config.yaml"}
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
	_, statErr := os.Stat(*configPath)
	switch {
	case statErr == nil:
		// We found the file, ensure it has identifying fields before establishing any opamp connections
		return ensureIdentity(*configPath)
	case errors.Is(statErr, os.ErrNotExist):
		var ok bool

		// manager.ymal file does *not* exist, create file using env variables
		newConfig := &opamp.Config{}

		// Endpoint is only required env
		newConfig.Endpoint, ok = os.LookupEnv(endpointENV)
		if !ok {
			// Envs were not found and statErr is os.ErrNotExist so return that
			return statErr
		}

		newConfig.AgentID, ok = os.LookupEnv(agentIDENV)
		if !ok {
			newConfig.AgentID = uuid.New().String()
		}

		if sk, ok := os.LookupEnv(secretkeyENV); ok {
			newConfig.SecretKey = &sk
		}

		if an, ok := os.LookupEnv(agentNameENV); ok {
			newConfig.AgentName = &an
		}

		if label, ok := os.LookupEnv(labelsENV); ok {
			newConfig.Labels = &label
		}

		data, err := yaml.Marshal(newConfig)
		if err != nil {
			return fmt.Errorf("failed to marshal config: %w", err)
		}

		// write data to a manager.yaml file, with 0600 file permission
		if err := os.WriteFile(*configPath, data, 0600); err != nil {
			return fmt.Errorf("failed to write config file created from ENVs: %w", err)
		}
		return nil
	}
	// Return non os.ErrNotExist
	return statErr
}

func ensureIdentity(configPath string) error {
	cBytes, err := os.ReadFile(filepath.Clean(configPath))
	if err != nil {
		return fmt.Errorf("unable to read file: %w", err)
	}
	var candidateConfig opamp.Config
	err = yaml.Unmarshal(cBytes, &candidateConfig)
	if err != nil {
		return fmt.Errorf("unable to interpret config file: %w", err)
	}

	// if they already have an AgentID, don't generate it
	if candidateConfig.AgentID != "" {
		return nil
	}

	candidateConfig.AgentID = uuid.NewString()
	newBytes, err := yaml.Marshal(candidateConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal sanitized config: %w", err)
	}

	if err = os.WriteFile(filepath.Clean(configPath), newBytes, 0600); err != nil {
		return fmt.Errorf("failed to rewrite manager config with identifying fields: %w", err)
	}
	return nil
}
