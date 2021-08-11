package env

import (
	"os"
	"path"
	"strconv"
)

const collectorHomePathEnvVar = "OIQ_COLLECTOR_HOME"
const launcherPPIDEnvVar = "COL_PPID"

// HomeDir returns the base directory of the collector.
func HomeDir() string {
	return os.Getenv(collectorHomePathEnvVar)
}

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

func LogDir() string {
	return path.Join(HomeDir(), "log")
}

func ConfigDir() string {
	// TODO: We might want to change from 'config/current' to just 'config' at some point
	return path.Join(HomeDir(), "config", "current")
}

func DefaultLoggingConfigFile() string {
	return path.Join(ConfigDir(), "collector-logging.yaml")
}

func DefaultRemoteConfigFile() string {
	return path.Join(ConfigDir(), "collector-remote.yaml")
}
