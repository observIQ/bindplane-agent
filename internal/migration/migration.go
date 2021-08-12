/*
* The migration package contains code to migrate from bpagent
* to out opentelemetry collector distro with minimal changes in behaviour
 */
package migration

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/observiq/observiq-collector/internal/env"
	"gopkg.in/yaml.v3"
)

// ShouldMigrate returns whether the collector should perform a migration from bpagent configs to the collector.
//  This will only be true if an install of bpAgent is detected (BP_AGENT_HOME is set)
func ShouldMigrate(envProvider env.BPEnvProvider) (bool, error) {
	return bpAgentInstalled(envProvider)
}

// Migrate moves and migrates all configs from BP_HOME to OIQ_COLLECTOR_HOME.
//  OIQ_COLLECTOR_HOME may be equal to BP_HOME with no conflicts.
func Migrate(envProvider env.EnvProvider, bpEnvProvider env.BPEnvProvider) error {
	migrateLogging, err := shouldMigrateLoggingConfig(envProvider)
	if err != nil {
		return fmt.Errorf("failed to determine if logging config should be migrated: %w", err)
	}

	if migrateLogging {
		err := migrateLoggingConfig(envProvider, bpEnvProvider)
		if err != nil {
			return fmt.Errorf("failed to migrate logging config: %w", err)
		}
	}

	migrateRemote, err := shouldMigrateRemoteConfig(envProvider)
	if err != nil {
		return fmt.Errorf("failed to determine if remote config should be migrated: %w", err)
	}

	if migrateRemote {
		err := migrateRemoteConfig(envProvider, bpEnvProvider)
		if err != nil {
			return fmt.Errorf("failed to migrate logging config: %w", err)
		}
	}

	return nil
}

// migrateLoggingConfig reads and re logging config to collector config
func migrateLoggingConfig(envProvider env.EnvProvider, bpEnvProvider env.BPEnvProvider) error {
	loggingConfigPath := bpEnvProvider.LoggingConfig()
	bpLogConfig, err := LoadBPLogConfig(loggingConfigPath)
	if err != nil {
		return err
	}

	logConf, err := BPLogConfigToLogConfig(*bpLogConfig)
	if err != nil {
		return err
	}

	configBytes, err := yaml.Marshal(logConf)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(envProvider.DefaultLoggingConfigFile(), configBytes, 0666)
	if err != nil {
		return err
	}

	return nil
}

func migrateRemoteConfig(envProvider env.EnvProvider, bpEnvProvider env.BPEnvProvider) error {
	config, err := BPRemoteConfig(bpEnvProvider)
	if err != nil {
		return err
	}

	// When this config is loaded, viper does something like this: (yaml -> mapstructure -> config)
	// So here, we do this in reverse: (config -> mapstructure -> yaml)
	// This ensures consistency between the saved file and loading code (which uses mapstructure tags)
	configAsMap := make(map[string]interface{})
	err = mapstructure.Decode(config, &configAsMap)
	if err != nil {
		return err
	}

	configBytes, err := yaml.Marshal(configAsMap)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(envProvider.DefaultManagerConfigFile(), configBytes, 0666)
	if err != nil {
		return err
	}

	return nil
}

// shouldMigrateLoggingConfig determines if the logging config should be migrated.
//  Currently, it will return true if the logging config file does not exist
func shouldMigrateLoggingConfig(envProvider env.EnvProvider) (bool, error) {
	exists, err := fileExists(envProvider.DefaultLoggingConfigFile())
	return !exists, err
}

// shouldMigrateRemoteConfig determines if the logging config should be migrated.
//  Currently, it will return true if the logging config file does not exist
func shouldMigrateRemoteConfig(envProvider env.EnvProvider) (bool, error) {
	exists, err := fileExists(envProvider.DefaultManagerConfigFile())
	return !exists, err
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
