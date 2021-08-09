package task

import (
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/observiq/observiq-collector/collector"
	"go.opentelemetry.io/collector/config/configtest"
	"gopkg.in/yaml.v3"
)

// Reconfigure is the task type assigned to reconfiguring the collector.
const Reconfigure Type = "reconfigure"

// ReconfigureParams are the parameters supplied in a reconfigure task.
type ReconfigureParams struct {
	Config struct {
		Pipeline []map[string]interface{}
	}
}

// getStanzaPipeline returns the stanza pipeline contained in the reconfigure task.
// Prior to using open telemetry, the observiq agent relied on a cabin operator.
// Until this is removed, we must sanitize the supplied config.
func (r *ReconfigureParams) getStanzaPipeline() []map[string]interface{} {
	pipeline := []map[string]interface{}{}

	for _, c := range r.Config.Pipeline {
		if c["type"] == "cabin_output" {
			continue
		}

		delete(c, "output")
		pipeline = append(pipeline, c)
	}

	return pipeline
}

// ExecuteReconfigure will execute a reconfigure task.
func ExecuteReconfigure(task Task, observiqCollector *collector.Collector, configPath string) Response {
	response := Response{
		ID:   task.ID,
		Type: task.Type,
	}

	if observiqCollector == nil {
		response.Status = Exception
		response.Message = "task received nil collector"
		return response
	}

	if task.Type != Reconfigure {
		response.Status = Exception
		response.Message = "task type does not match reconfigure"
		return response
	}

	var params ReconfigureParams
	err := mapstructure.Decode(task.Parameters, &params)
	if err != nil {
		response.Status = Exception
		response.Message = fmt.Sprintf("unable to decode task parameters: %s", err)
		return response
	}

	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		response.Status = Exception
		response.Message = fmt.Sprintf("failed to read existing config: %s", err)
		return response
	}

	configMap := make(map[string]interface{})
	err = yaml.Unmarshal(configBytes, &configMap)
	if err != nil {
		response.Status = Exception
		response.Message = fmt.Sprintf("failed to decode existing config: %s", err)
		return response
	}

	receivers, ok := configMap["receivers"].(map[string]interface{})
	if !ok {
		response.Status = Exception
		response.Message = fmt.Sprintf("failed to get receivers from config: %s", err)
		return response
	}

	stanza := receivers["stanza"].(map[string]interface{})
	if !ok {
		response.Status = Exception
		response.Message = fmt.Sprintf("failed to get stanza from receivers: %s", err)
		return response
	}

	stanza["pipeline"] = params.getStanzaPipeline()
	newBytes, err := yaml.Marshal(&configMap)
	if err != nil {
		response.Status = Exception
		response.Message = fmt.Sprintf("failed to convert new config to yaml: %s", err)
		return response
	}

	err = os.WriteFile(configPath, newBytes, 0666)
	if err != nil {
		response.Status = Exception
		response.Message = fmt.Sprintf("failed to write new config: %s", err)
		return response
	}

	factories, err := collector.DefaultFactories()
	if err != nil {
		response.Status = Exception
		response.Message = fmt.Sprintf("failed to get default factories: %s", err)
		return response
	}

	_, err = configtest.LoadConfigAndValidate(configPath, factories)
	if err != nil {
		response.Status = Exception
		response.Message = fmt.Sprintf("failed dry run of new config: %s", err)
		_ = os.WriteFile(configPath, configBytes, 0666)
		return response
	}

	observiqCollector.Stop()
	err = observiqCollector.Run()
	if err != nil {
		response.Status = Exception
		response.Message = fmt.Sprintf("failed to restart collector: %s", err)
		return response
	}

	return response
}
