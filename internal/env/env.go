package env

import (
	"os"
	"path"
	"strconv"
)

const collectorHomePathEnvVar = "OIQ_COLLECTOR_HOME"
const launcherPPIDEnvVar = "COL_PPID"

type EnvProvider interface {
	HomeDir() string
	LogDir() string
	ConfigDir() string
	DefaultLoggingConfigFile() string
	DefaultRemoteConfigFile() string
}

type defaultEnvProvider struct{}

// HomeDir returns the base directory of the collector.
func (defaultEnvProvider) HomeDir() string {
	return os.Getenv(collectorHomePathEnvVar)
}

func (p defaultEnvProvider) LogDir() string {
	return path.Join(p.HomeDir(), "log")
}

func (p defaultEnvProvider) ConfigDir() string {
	// TODO: We might want to change from 'config/current' to just 'config' at some point
	return path.Join(p.HomeDir(), "config", "current")
}

func (p defaultEnvProvider) DefaultLoggingConfigFile() string {
	return path.Join(p.ConfigDir(), "collector-logging.yaml")
}

func (p defaultEnvProvider) DefaultRemoteConfigFile() string {
	return path.Join(p.ConfigDir(), "collector-remote.yaml")
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
