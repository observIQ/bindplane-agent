package env

import (
	"os"
	"path"
	"runtime"
)

const bpHomeEnvVar = "BP_AGENT_HOME"

type BPEnvProvider interface {
	RemoteConfig() string
	LoggingConfig() string
}

type defaultBPEnvProvider struct{}

func (p defaultBPEnvProvider) RemoteConfig() string {
	return path.Join(bpConfigDir(), "remote.yaml")
}

func (p defaultBPEnvProvider) LoggingConfig() string {
	return path.Join(bpConfigDir(), "logging.yaml")
}

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

func bpConfigDir() string {
	return path.Join(bpHomeDir(), "config")
}

var DefaultBPEnvProvider BPEnvProvider = defaultBPEnvProvider{}
