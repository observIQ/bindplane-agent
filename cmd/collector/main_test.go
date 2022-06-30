package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/observiq/observiq-otel-collector/opamp"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestCheckManagerConfig(t *testing.T) {
	manager := "./manager.yaml"
	err := checkManagerConfig(&manager)
	require.NoError(t, err)
}

func TestCheckManagerConfigNoFile(t *testing.T) {
	exec.Command("rm", "-r", "./manager.yaml").Run()
	manager := "./manager.yaml"
	err := checkManagerConfig(&manager)
	require.Error(t, err)

	os.Setenv(ENDPOINT, "0.0.0.0")
	defer os.Unsetenv(ENDPOINT)

	os.Setenv(SECRET_KEY, "secret_key")
	defer os.Unsetenv(SECRET_KEY)

	os.Setenv(LABELS, "this is a label")
	defer os.Unsetenv(LABELS)
	defer os.Unsetenv(AGENT_ID)

	manager = "./manager.yaml"
	err = checkManagerConfig(&manager)
	require.NoError(t, err)

	dat, err := os.ReadFile("./manager.yaml")
	out := &opamp.Config{}
	err = yaml.Unmarshal(dat, out)
	require.Equal(t,
		&opamp.Config{
			Endpoint: "0.0.0.0",
		},
		&opamp.Config{
			Endpoint: out.Endpoint,
		})
}
