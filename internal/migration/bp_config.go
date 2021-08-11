package migration

import (
	"io/ioutil"
	"path"

	"github.com/mitchellh/mapstructure"
	"github.com/observiq/observiq-collector/internal/env"
	"github.com/observiq/observiq-collector/internal/logging"
	"github.com/observiq/observiq-collector/manager"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

type bpLogConfig struct {
	Level        *zapcore.Level `mapstructure:"level" yaml:"level"`
	MaxBackups   *int           `mapstructure:"max_backups" yaml:"max_backups"`
	MaxMegabytes *int           `mapstructure:"max_megabytes" yaml:"max_megabytes"`
	MaxDays      *int           `mapstructure:"max_days" yaml:"max_days"`
}

func BPRemoteConfig() (*manager.Config, error) {
	// bpagents remote.yaml is compatible with manager.Config (manager.Config is a superset of it)
	remoteConfigPath := path.Join(env.BPConfigDir(), "remote.yaml")
	return manager.ReadConfig(remoteConfigPath)
}

func LoadBPLogConfig() (*bpLogConfig, error) {
	loggingConfigPath := path.Join(env.BPConfigDir(), "logging.yaml")

	bytes, err := ioutil.ReadFile(loggingConfigPath)

	if err != nil {
		return nil, err
	}

	r := &bpLogConfig{}
	err = yaml.Unmarshal(bytes, r)

	return r, err
}

func BPLogConfigToLogConfig(c bpLogConfig) (*logging.Config, error) {
	config, err := logging.DefaultConfig()
	if err != nil {
		return nil, err
	}

	err = mapstructure.Decode(c, &config.Collector)
	if err != nil {
		return nil, err
	}

	err = mapstructure.Decode(c, &config.Manager)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// bpAgentInstalled returns whether a BPAgent install is detected
//  It accomplishes this by checking if the files needed to migrate exist in certain paths
func bpAgentInstalled() (bool, error) {
	remoteConfigPath := path.Join(env.BPConfigDir(), "remote.yaml")
	loggingConfigPath := path.Join(env.BPConfigDir(), "logging.yaml")

	for _, file := range []string{remoteConfigPath, loggingConfigPath} {
		exists, err := fileExists(file)
		if err != nil {
			return false, err
		}

		if !exists {
			return false, nil
		}
	}

	return true, nil
}
