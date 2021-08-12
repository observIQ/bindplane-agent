package env

import (
	"os"
	"path"
	"path/filepath"
	"strconv"
)

const fileLogEnableEnvVar = "OIQ_COLLECTOR_FILE_LOG"
const collectorHomePathEnvVar = "OIQ_COLLECTOR_HOME"
const launcherPPIDEnvVar = "COL_PPID"

/*
	Check if logging to file is enabled (env variable is set and non-zero).
	If it is enabled, OIQ_COLLECTOR_HOME is expected to be set, as well.
*/
func IsFileLoggingEnabled() bool {
	s := os.Getenv(fileLogEnableEnvVar)
	if s == "" || s == "0" {
		return false
	}
	return true
}

/*
	Gets the logging path relative to OIQ_COLLECTOR_HOME environment variable
*/
func GetLoggingPath() (string, bool) {
	p, ok := os.LookupEnv(collectorHomePathEnvVar)
	if !ok {
		return "", false
	}

	fp, err := filepath.Abs(path.Join(p, "log", "collector.log"))
	if err != nil {
		return "", false
	}

	return fp, true
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
