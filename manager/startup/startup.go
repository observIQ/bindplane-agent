package startup

import (
	"github.com/observiq/observiq-collector/collector"
	"github.com/observiq/observiq-collector/internal/env"
	"github.com/observiq/observiq-collector/internal/pipelinereader"
	"github.com/observiq/observiq-collector/manager/message"
	"github.com/observiq/observiq-collector/receiver/logsreceiver"
)

// Startup is an object that is reported when the collector gets started up
type Startup struct {
	// note that that this needs to be reported back up as bpHome
	// can probably change server side
	OIQHome        string         `json:"bpHome" mapstructure:"bpHome"`
	TemplateID     string         `json:"configurationID" mapstructure:"configuration_id"`
	MacAddress     string         `json:"macAddress" mapstructure:"macAddress"`
	OSDetails      string         `json:"osDetails" mapstructure:"osDetails"`
	AgentName      string         `json:"agentName" mapstructure:"agentName"`
	LogAgentConfig logAgentConfig `mapstructure:"log_agent_config"`
}

type logAgentConfig struct {
	Pipeline *logsreceiver.OperatorConfigs `json:"pipeline" mapstructure:"pipeline"`
}

// New returns a populated struct for a startup message
func New(templateID, name string, col *collector.Collector) Startup {
	laCfg := getLogAgentConfig(col)
	return Startup{
		OIQHome:        env.HomeDir(),
		AgentName:      name,
		TemplateID:     templateID,
		MacAddress:     FindMACAddressOrUnknown(),
		OSDetails:      GetDetails(),
		LogAgentConfig: logAgentConfig{Pipeline: laCfg},
	}
}

// ToMessage converts the StartupStruct to a message
func (st Startup) ToMessage() *message.Message {
	msg, _ := message.New("onStartup", st)
	return msg
}

func getLogAgentConfig(col *collector.Collector) *logsreceiver.OperatorConfigs {
	pipeline, err := pipelinereader.Read(col)
	if err != nil {
		emptyList := logsreceiver.OperatorConfigs{}
		pipeline = emptyList
	}
	return &pipeline
}
