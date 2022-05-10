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
	instanceID   string
	instanceName string
	serviceName  string // TODO figure out how to get this
	version      string
	labels       *string
	oSArch       string
	oSDetails    string
	oSFamily     string
	hostname     string
	mac          string
}

// NewIdentity constructs a new identity for this collector
func NewIdentity(logger *zap.SugaredLogger, config opamp.Config) *identity {
	// Fill out identifying attributes
	identifyingAttrs := []*protobufs.KeyValue{
		opamp.StringKeyValue("service.name", "com.observiq.agent"), // TODO parse this out
		opamp.StringKeyValue("service.version", version.Version()),
	}

	// Fill out non-identifying attributes
	nonIdentifyingAttrs := []*protobufs.KeyValue{
		opamp.StringKeyValue("os.arch", runtime.GOARCH),
		opamp.StringKeyValue("os.family", runtime.GOOS),
	}

	// Grab various fields from OS
	hostname, err := ios.Hostname()
	if err != nil {
		logger.Warn("Failed to retrieve hostname for collector. Creating partial identity", zap.Error(err))
	} else {
		nonIdentifyingAttrs = append(nonIdentifyingAttrs, opamp.StringKeyValue("host.name", hostname))
	}

	name, err := ios.Name()
	if err != nil {
		logger.Warn("Failed to retrieve host details on collector. Creating partial identity", zap.Error(err))
	} else {
		nonIdentifyingAttrs = append(nonIdentifyingAttrs, opamp.StringKeyValue("os.details", name))
	}

	// parse config file for attributes
	configIDAttrs, configNonIDAttrs := parseConfigAttrs(config, hostname)

	identifyingAttrs = append(identifyingAttrs, configIDAttrs...)
	nonIdentifyingAttrs = append(nonIdentifyingAttrs, configNonIDAttrs...)

	// TODO move all the identifying code into a identity thing for opamp
	return &identity{
		instanceID:   config.AgentID,
		instanceName: name,
		serviceName:  name,
		version:      "",
		labels:       new(string),
		oSArch:       "",
		oSDetails:    "",
		oSFamily:     "",
		hostname:     hostname,
		mac:          "",
	}
}

func parseConfigAttrs(config opamp.Config, hostname string) (identifyingAttrs, nonIdentifyingAttrs []*protobufs.KeyValue) {
	identifyingAttrs = []*protobufs.KeyValue{
		opamp.StringKeyValue("service.instance.id", config.AgentID),
	}

	if config.AgentName != nil {
		identifyingAttrs = append(identifyingAttrs, opamp.StringKeyValue("service.instance.name", *config.AgentName))
	} else if hostname != "" {
		identifyingAttrs = append(identifyingAttrs, opamp.StringKeyValue("service.instance.name", hostname))

	}

	nonIdentifyingAttrs = make([]*protobufs.KeyValue, 0)
	if config.Labels != nil {
		nonIdentifyingAttrs = append(nonIdentifyingAttrs, opamp.StringKeyValue("service.labels", *config.Labels))
	}

	return
}
