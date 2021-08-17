package startup

import (
	"github.com/observiq/observiq-collector/collector"
	"github.com/observiq/observiq-collector/internal/env"
	"github.com/observiq/observiq-collector/internal/pipelinereader"
	"github.com/observiq/observiq-collector/manager/message"
)

// Startup is an object that is reported when the collector gets started up
type Startup struct {
	// note that that this needs to be reported back up as bpHome
	// can probably change server side
	OIQHome        string                 `json:"bpHome" mapstructure:"bpHome"`
	TemplateID     string                 `json:"configurationID" mapstructure:"configuration_id"`
	MacAddress     string                 `json:"macAddress" mapstructure:"macAddress"`
	OSDetails      string                 `json:"osDetails" mapstructure:"osDetails"`
	AgentName      string                 `json:"agentName" mapstructure:"agentName"`
	LogAgentConfig map[string]interface{} `mapstructure:"log_agent_config"`
}

// New returns a populated struct for a startup message
func New(templateID, name string, col *collector.Collector) Startup {
	return Startup{
		OIQHome:        env.HomeDir(),
		AgentName:      name,
		TemplateID:     templateID,
		MacAddress:     FindMACAddressOrUnknown(),
		OSDetails:      GetDetails(),
		LogAgentConfig: getLogAgentConfig(col),
	}
}

// ToMessage converts the StartupStruct to a message
func (st Startup) ToMessage() *message.Message {
	msg, _ := message.New("onStartup", st)
	return msg
}

func getLogAgentConfig(col *collector.Collector) map[string]interface{} {
	pipeline, err := pipelinereader.Read(col)
	if err != nil {
		emptyMap := make(map[string]interface{})
		pipeline = emptyMap
	}
	return pipeline
}
