package task

import (
	"errors"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/observiq/observiq-collector/collector"
	"gopkg.in/yaml.v3"
)

// Reconfigure is the task type assigned to reconfiguring the collector.
const Reconfigure Type = "reconfigure"

// ReconfigureParams are the parameters supplied in a reconfigure task.
type ReconfigureParams struct {
	Config StanzaConfig
}

// StanzaConfig is the configuration of a stanza receiver.
type StanzaConfig struct {
	Pipeline StanzaPipeline
}

// StanzaPipeline represents a stanza pipeline.
type StanzaPipeline []map[string]interface{}

// getStanzaPipeline returns the stanza pipeline contained in the reconfigure task.
//
// Since the stanza receiver always sends from the last defined operator in its config,
// we need to ensure that the cabin_output operator always comes last.
//
// Also, if no operators exist in the pipeline, we need to ensure that a noop
// operator is configured. Otherwise, the collector will fail to start.
func (r *ReconfigureParams) getStanzaPipeline() StanzaPipeline {
	pipeline := r.Config.Pipeline

	for i, c := range r.Config.Pipeline {
		if c["type"] == "cabin_output" {
			pipeline = append(pipeline[:i], pipeline[i+1:]...)
			pipeline = append(pipeline, c)
		}
	}

	if len(pipeline) == 0 {
		noop := map[string]interface{}{"type": "noop"}
		pipeline = append(pipeline, noop)
	}

	return pipeline
}

// ExecuteReconfigure will execute a reconfigure task.
func ExecuteReconfigure(task *Task, collector *collector.Collector) Response {
	if task.Type != Reconfigure {
		err := errors.New("invalid type")
		return task.Failure("task is not a reconfigure", err)
	}

	var params ReconfigureParams
	err := mapstructure.Decode(task.Parameters, &params)
	if err != nil {
		return task.Failure("unable to decode parameters", err)
	}

	configBytes, err := os.ReadFile(collector.ConfigPath())
	if err != nil {
		return task.Failure("failed to read existing config", err)
	}

	configMap := make(map[string]interface{})
	err = yaml.Unmarshal(configBytes, &configMap)
	if err != nil {
		return task.Failure("failed to decode existing config", err)
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

	stanza["pipeline"] = params.getStanzaPipeline()
	updatedConfig, err := yaml.Marshal(&configMap)
	if err != nil {
		return task.Failure("failed to convert new config to yaml", err)
	}

	err = os.WriteFile(collector.ConfigPath(), updatedConfig, 0666)
	if err != nil {
		return task.Failure("failed to write new config", err)
	}

	err = collector.ValidateConfig()
	if err != nil {
		_ = os.WriteFile(collector.ConfigPath(), configBytes, 0666)
		return task.Failure("new config failed validation", err)
	}

	err = collector.Restart()
	if err != nil {
		return task.Failure("failed to restart collector", err)
	}

	return task.Success()
}
