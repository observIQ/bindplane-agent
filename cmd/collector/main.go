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

	"github.com/observiq/observiq-otel-collector/collector"
	"github.com/observiq/observiq-otel-collector/internal/logging"
	"github.com/observiq/observiq-otel-collector/internal/service"
	"github.com/observiq/observiq-otel-collector/internal/version"
	"github.com/observiq/observiq-otel-collector/opamp"
	"github.com/oklog/ulid/v2"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

const (
	// env variable name constants
	endpointENV    = "OPAMP_ENDPOINT"
	agentIDENV     = "OPAMP_AGENT_ID"
	secretkeyENV   = "OPAMP_SECRET_KEY" //#nosec G101
	labelsENV      = "OPAMP_LABELS"
	agentNameENV   = "OPAMP_AGENT_NAME"
	configPathENV  = "CONFIG_YAML_PATH"
	managerPathENV = "MANAGER_YAML_PATH"
	loggingPathENV = "LOGGING_YAML_PATH"
)

func main() {
	collectorConfigPaths := pflag.StringSlice("config", getDefaultCollectorConfigPaths(), "the collector config path")
	managerConfigPath := pflag.String("manager", getDefaultManagerConfigPath(), "The configuration for remote management")
	loggingConfigPath := pflag.String("logging", getDefaultLoggingConfigPath(), "the collector logging config path")

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

		collectorConfigPath := (*collectorConfigPaths)[0]

		// Check for existing rollback files for collector config.
		// If any exist it's likely the case that the collector crashed during reconfigure.
		logger.Debug("Checking for existing rollback files")
		if err := checkForCollectorRollbackConfig(collectorConfigPath); err != nil {
			// Log an error rather than exit as we should still have a config to run with
			logger.Error("Error occurred while checking for collector config rollbacks", zap.Error(err))
		}

		runnableService, err = service.NewManagedCollectorService(col, logger, *managerConfigPath, collectorConfigPath, *loggingConfigPath)
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
	}
	return []string{"./config.yaml"}
}

func getDefaultManagerConfigPath() string {
	mp, ok := os.LookupEnv(managerPathENV)
	if ok {
		return mp
	}
	return "./manager.yaml"
}

func getDefaultLoggingConfigPath() string {
	lp, ok := os.LookupEnv(loggingPathENV)
	if ok {
		return lp
	}
	return logging.DefaultConfigPath
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

		// manager.yaml file does *not* exist, create file using env variables
		newConfig := &opamp.Config{}

		// Endpoint is only required env
		newConfig.Endpoint, ok = os.LookupEnv(endpointENV)
		if !ok {
			// Envs were not found and statErr is os.ErrNotExist so return that
			return statErr
		}

		newConfig.AgentID, ok = os.LookupEnv(agentIDENV)
		if !ok {
			newConfig.AgentID = ulid.Make().String()
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

	// If the AgentID is not a ULID (legacy ID or empty) then we need to generate a ULID as the AgentID.
	if _, err := ulid.Parse(candidateConfig.AgentID); err == nil {
		return nil
	}

	candidateConfig.AgentID = ulid.Make().String()
	newBytes, err := yaml.Marshal(candidateConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal sanitized config: %w", err)
	}

	if err = os.WriteFile(filepath.Clean(configPath), newBytes, 0600); err != nil {
		return fmt.Errorf("failed to rewrite manager config with identifying fields: %w", err)
	}
	return nil
}

func checkForCollectorRollbackConfig(configPath string) error {
	cleanPath := filepath.Clean(configPath)
	rollbackFileName := fmt.Sprintf("%s.rollback", cleanPath)

	// Check to see if rollback file exists
	_, err := os.Stat(rollbackFileName)

	// No rollback file just return
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	// Copy rollback file and delete original
	contents, err := os.ReadFile(rollbackFileName)
	if err != nil {
		return fmt.Errorf("error while reading in collector rollback file: %w", err)
	}

	if err := os.WriteFile(cleanPath, contents, 0600); err != nil {
		return fmt.Errorf("error while writing rollback contents onto config: %w", err)
	}

	if err := os.Remove(rollbackFileName); err != nil {
		return fmt.Errorf("error while cleaning up rollback file: %w", err)
	}

	return nil
}
