package env

import (
	"os"
	"path"
	"runtime"
)

const bpHomeEnvVar = "BP_AGENT_HOME"

// A BPEnvProvider gives information about paths and the environment in which an install of BPAgent may be located.
type BPEnvProvider interface {
	// RemoteConfig returns the path to BPAgent's remote.yaml
	RemoteConfig() string
	// LoggingConfig returns the path to BPAgent's logging.yaml
	LoggingConfig() string
}

type defaultBPEnvProvider struct{}

// RemoteConfig returns the path to BPAgent's remote.yaml
func (p defaultBPEnvProvider) RemoteConfig() string {
	return path.Join(bpConfigDir(), "remote.yaml")
}

// LoggingConfig returns the path to BPAgent's logging.yaml
func (p defaultBPEnvProvider) LoggingConfig() string {
	return path.Join(bpConfigDir(), "logging.yaml")
}

// bpHomeDir returns the "best guess" at where BP_AGENT_HOME is (or just the BP_AGENT_HOME environment variable, if defined)
func bpHomeDir() string {
	// This will return BPHome if BPAgent is installed on Windows (environment variable is system-wide)
	//  This also allows a point to override BP_HOME
	if home, ok := os.LookupEnv(bpHomeEnvVar); ok {
		return home
	}
	// On other OS's, the BP_AGENT_HOME variable is at the service level, so we need to replicate script logic for default paths
	switch runtime.GOOS {
	case "darwin":
		home := os.Getenv("HOME")
		return path.Join(home, "observiq-agent")
	case "linux":
		return path.Join("opt", "observiq-agent")
	default:
		return path.Join("observiq-agent")
	}
}

// bpConfigDir returns the path to BPAgent's config directory.
func bpConfigDir() string {
	return path.Join(bpHomeDir(), "config")
}

// DefaultEnvProvider is the default provider for BPAgent environment information.
var DefaultBPEnvProvider BPEnvProvider = defaultBPEnvProvider{}
