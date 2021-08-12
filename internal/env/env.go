package env

import (
	"os"
	"path"
	"strconv"
)

const collectorHomePathEnvVar = "OIQ_COLLECTOR_HOME"
const launcherPPIDEnvVar = "COL_PPID"

type EnvProvider interface {
	LogDir() string
	DefaultLoggingConfigFile() string
	DefaultManagerConfigFile() string
}

type defaultEnvProvider struct{}

// homeDir returns the base directory of the collector.
func homeDir() string {
	return os.Getenv(collectorHomePathEnvVar)
}

func configDir() string {
	// TODO: We might want to change from 'config/current' to just 'config' at some point
	return path.Join(homeDir(), "config", "current")
}

func (p defaultEnvProvider) LogDir() string {
	return path.Join(homeDir(), "log")
}

func (p defaultEnvProvider) DefaultLoggingConfigFile() string {
	return path.Join(configDir(), "collector-logging.yaml")
}

func (p defaultEnvProvider) DefaultManagerConfigFile() string {
	return path.Join(configDir(), "collector-remote.yaml")
}

var DefaultEnvProvider EnvProvider = defaultEnvProvider{}

// GetLauncherPPID returns the launcher ppid contained in the `COL_PPID` environment variable.
func GetLauncherPPID() int {
	value, ok := os.LookupEnv(launcherPPIDEnvVar)
	if !ok {
		return 0
	}

	ppid, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return ppid
}
