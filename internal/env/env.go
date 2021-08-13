package env

import (
	"os"
	"strconv"
)

const collectorHomePathEnvVar = "OIQ_COLLECTOR_HOME"
const launcherIDEnvVar = "LAUNCHER_ID"

// HomeDir returns the base directory of the collector.
func HomeDir() string {
	return os.Getenv(collectorHomePathEnvVar)
}

// GetLauncherID returns the launcher id contained in the `LAUNCHER_ID` environment variable.
func GetLauncherID() int {
	value := os.Getenv(launcherIDEnvVar)

	id, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return id
}
