package observiq

import (
	"runtime"

	ios "github.com/observiq/observiq-otel-collector/internal/os"
	"github.com/observiq/observiq-otel-collector/internal/version"
	"github.com/observiq/observiq-otel-collector/opamp"
	"github.com/open-telemetry/opamp-go/protobufs"
	"go.uber.org/zap"
)

// identity contains identifying information about the Collector
type identity struct {
	agentID     string
	agentName   *string
	serviceName string
	version     string
	labels      *string
	oSArch      string
	oSDetails   string
	oSFamily    string
	hostname    string
	mac         string
}

// NewIdentity constructs a new identity for this collector
func NewIdentity(logger *zap.SugaredLogger, config opamp.Config) *identity {
	// Grab various fields from OS
	hostname, err := ios.Hostname()
	if err != nil {
		logger.Warn("Failed to retrieve hostname for collector. Creating partial identity", zap.Error(err))
	}

	name, err := ios.Name()
	if err != nil {
		logger.Warn("Failed to retrieve host details on collector. Creating partial identity", zap.Error(err))
	}

	return &identity{
		agentID:     config.AgentID,
		agentName:   config.AgentName,
		serviceName: "com.observiq.collector", // TODO figure this out
		version:     version.Version(),
		labels:      config.Labels,
		oSArch:      runtime.GOARCH,
		oSDetails:   name,
		oSFamily:    runtime.GOOS,
		hostname:    hostname,
		mac:         ios.MACAddress(),
	}
}

func (i *identity) ToAgentDescription() *protobufs.AgentDescription {
	identifyingAttributes := []*protobufs.KeyValue{
		opamp.StringKeyValue("service.instance.id", i.agentID),
		opamp.StringKeyValue("service.name", i.serviceName),
		opamp.StringKeyValue("service.version", i.version),
	}

	if i.agentName != nil {
		identifyingAttributes = append(identifyingAttributes, opamp.StringKeyValue("service.instance.name", *i.agentName))
	} else {
		identifyingAttributes = append(identifyingAttributes, opamp.StringKeyValue("service.instance.name", i.hostname))
	}

	nonIdentifyingAttributes := []*protobufs.KeyValue{
		opamp.StringKeyValue("os.arch", i.oSArch),
		opamp.StringKeyValue("os.details", i.oSDetails),
		opamp.StringKeyValue("os.family", i.oSFamily),
		opamp.StringKeyValue("host.name", i.hostname),
		opamp.StringKeyValue("host.mac_address", i.mac),
	}

	if i.labels != nil {
		nonIdentifyingAttributes = append(nonIdentifyingAttributes, opamp.StringKeyValue("service.labels", *i.labels))
	}

	// Create agent description.
	return &protobufs.AgentDescription{
		IdentifyingAttributes:    identifyingAttributes,
		NonIdentifyingAttributes: nonIdentifyingAttributes,
	}
}
