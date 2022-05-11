package opamp

import "github.com/open-telemetry/opamp-go/protobufs"

// ConfigManager handles remote configuration of local configs
type ConfigManager interface {
	// AddConfig adds a config to be tracked by the config manager with it's corresponding validator function.
	AddConfig(configName, configPath string, validator ValidatorFunc)

	// ComposeEffectiveConfig reads in all config files and calculates the effective config
	ComposeEffectiveConfig() (*protobufs.EffectiveConfig, error)

	// ApplyConfigChanges compares the remoteConfig to the existing and applies changes.
	// Calculates new effective config
	ApplyConfigChanges(remoteConfig *protobufs.AgentRemoteConfig) (effectiveConfig *protobufs.EffectiveConfig, changed bool, err error)
}
