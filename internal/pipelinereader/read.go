package pipelinereader

import (
	"os"

	"github.com/observiq/observiq-collector/collector"
	"github.com/observiq/observiq-collector/receiver/logsreceiver"
	"gopkg.in/yaml.v2"
)

type ConfigMap struct {
	Receivers map[string]baseStanzaReceiver `mapstructure:"receivers"`
	Exporters map[string]interface{}        `mapstructure:"exporters"`
}

type baseStanzaReceiver logsreceiver.Config

// Read reads the stanza portion of the collector config, this is so that the manager
// has visibility into the collector's config
func Read(c *collector.Collector) (logsreceiver.OperatorConfigs, error) {
	configBytes, err := os.ReadFile(c.ConfigPath())
	if err != nil {
		return nil, err
	}

	configMap := ConfigMap{}
	err = yaml.Unmarshal(configBytes, &configMap)
	if err != nil {
		return nil, err
	}

	stanza := configMap.Receivers["stanza"]
	return stanza.Pipeline, nil
}
