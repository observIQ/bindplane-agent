package pipelinereader

import (
	"os"

	"github.com/observiq/observiq-collector/collector"
	"gopkg.in/yaml.v2"
)

func Read(c *collector.Collector) (map[string]interface{}, error) {
	configBytes, err := os.ReadFile(c.ConfigPath())
	if err != nil {
		return nil, err
	}

	configMap := make(map[string]interface{})
	err = yaml.Unmarshal(configBytes, &configMap)
	if err != nil {
		return nil, err
	}

	receivers, ok := configMap["receivers"].(map[string]interface{})
	if !ok {
		receivers = make(map[string]interface{})
		configMap["receivers"] = receivers
	}

	stanza, ok := receivers["stanza"].(map[string]interface{})
	if !ok {
		stanza = make(map[string]interface{})
		receivers["stanza"] = stanza
	}
	return stanza, nil
}
