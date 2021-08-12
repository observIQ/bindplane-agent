package env

import (
	"os"
	"path"
	"runtime"
)

const bpHomeEnvVar = "BP_AGENT_HOME"

type BPEnvProvider interface {
	HasBpHome() bool
	BPHomeDir() string
	BPConfigDir() string
	RemoteConfig() string
	LoggingConfig() string
}

type defaultBPEnvProvider struct{}

func (defaultBPEnvProvider) HasBpHome() bool {
	_, hasEnv := os.LookupEnv(bpHomeEnvVar)
	return hasEnv
}

func (defaultBPEnvProvider) BPHomeDir() string {
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

func (p defaultBPEnvProvider) BPConfigDir() string {
	return path.Join(p.BPHomeDir(), "config")
}

func (p defaultBPEnvProvider) RemoteConfig() string {
	return path.Join(p.BPConfigDir(), "remote.yaml")
}

func (p defaultBPEnvProvider) LoggingConfig() string {
	return path.Join(p.BPConfigDir(), "logging.yaml")
}

var DefaultBPEnvProvider BPEnvProvider = defaultBPEnvProvider{}
