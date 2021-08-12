package env

import (
	"os"
	"path"
	"strconv"
)

const collectorHomePathEnvVar = "OIQ_COLLECTOR_HOME"
const launcherPPIDEnvVar = "COL_PPID"

// An EnvProvider ives information about paths and the environment in which the collector is installed
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

// configDir returns the directory where all collector configs are stored
func configDir() string {
	// TODO: We might want to change from 'config/current' to just 'config' at some point
	return path.Join(homeDir(), "config", "current")
}

// LogDir returns the path to the directory containing collector logs, by default.
func (p defaultEnvProvider) LogDir() string {
	return path.Join(homeDir(), "log")
}

// DefaultLoggingConfigFile returns the path to the default Logging config file
func (p defaultEnvProvider) DefaultLoggingConfigFile() string {
	return path.Join(configDir(), "collector-logging.yaml")
}

// DefaultManagerConfigFile returns the path to the default manager config file
func (p defaultEnvProvider) DefaultManagerConfigFile() string {
	return path.Join(configDir(), "collector-remote.yaml")
}

// DefaultEnvProvider is the default provider for environment information.
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
