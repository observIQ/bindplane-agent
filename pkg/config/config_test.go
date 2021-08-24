package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadMissingFile(t *testing.T) {
	path := "/missing/file"
	config, err := Read(path)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to read file")
	require.Nil(t, config)
}

func TestReadInvalidYAML(t *testing.T) {
	file, err := ioutil.TempFile("", "config")
	require.NoError(t, err)
	defer os.Remove(file.Name())

	_, err = file.WriteString("invalid yaml")
	require.NoError(t, err)

	config, err := Read(file.Name())
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to unmarshal")
	require.Nil(t, config)
}

func TestReadValid(t *testing.T) {
	file, err := ioutil.TempFile("", "config")
	require.NoError(t, err)
	defer os.Remove(file.Name())

	contents := `
receivers:
  test_receiver:
    key: value
exporters:
  test_exporter:
    key: value
service:
  pipelines:
    test_pipeline:
      receivers: [test_receiver]
      exporters: [test_exporter]`

	_, err = file.WriteString(contents)
	require.NoError(t, err)

	config, err := Read(file.Name())
	require.NoError(t, err)

	_, ok := config.Receivers["test_receiver"]
	require.True(t, ok)

	_, ok = config.Exporters["test_exporter"]
	require.True(t, ok)

	_, ok = config.Service.Pipelines["test_pipeline"]
	require.True(t, ok)
}
