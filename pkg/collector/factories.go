package collector

import (
	"github.com/GoogleCloudPlatform/opentelemetry-operations-collector/processor/normalizesumsprocessor"
	"github.com/observiq/observiq-collector/pkg/processor/resourceattributetransposerprocessor"
	"github.com/observiq/observiq-collector/pkg/receiver/logsreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/fileexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/googlecloudexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/observiqexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/bearertokenauthextension"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/healthcheckextension"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/oidcauthextension"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/pprofextension"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/storage/filestorage"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbyattrsprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricstransformprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/probabilisticsamplerprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourceprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filelogreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/jmxreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8sclusterreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kubeletstatsreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mongodbreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mysqlreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/postgresqlreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/statsdreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/syslogreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/tcplogreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/udplogreceiver"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/exporter/loggingexporter"
	"go.opentelemetry.io/collector/exporter/otlpexporter"
	"go.opentelemetry.io/collector/exporter/otlphttpexporter"
	"go.opentelemetry.io/collector/extension/ballastextension"
	"go.opentelemetry.io/collector/extension/zpagesextension"
	"go.opentelemetry.io/collector/processor/batchprocessor"
	"go.opentelemetry.io/collector/processor/memorylimiterprocessor"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
	"go.uber.org/multierr"
)

var defaultReceivers = []component.ReceiverFactory{
	logsreceiver.NewFactory(),
	otlpreceiver.NewFactory(),
	filelogreceiver.NewFactory(),
	jmxreceiver.NewFactory(),
	syslogreceiver.NewFactory(),
	tcplogreceiver.NewFactory(),
	udplogreceiver.NewFactory(),
	componenttest.NewNopReceiverFactory(),
	kubeletstatsreceiver.NewFactory(),
	k8sclusterreceiver.NewFactory(),
	mysqlreceiver.NewFactory(),
	statsdreceiver.NewFactory(),
	postgresqlreceiver.NewFactory(),
	mongodbreceiver.NewFactory(),
}

var defaultProcessors = []component.ProcessorFactory{
	groupbyattrsprocessor.NewFactory(),
	k8sattributesprocessor.NewFactory(),
	attributesprocessor.NewFactory(),
	resourceprocessor.NewFactory(),
	batchprocessor.NewFactory(),
	memorylimiterprocessor.NewFactory(),
	probabilisticsamplerprocessor.NewFactory(),
	resourceattributetransposerprocessor.NewFactory(),
	componenttest.NewNopProcessorFactory(),
	metricstransformprocessor.NewFactory(),
	normalizesumsprocessor.NewFactory(),
	resourcedetectionprocessor.NewFactory(),
}

var defaultExporters = []component.ExporterFactory{
	fileexporter.NewFactory(),
	otlpexporter.NewFactory(),
	otlphttpexporter.NewFactory(),
	observiqexporter.NewFactory(),
	loggingexporter.NewFactory(),
	componenttest.NewNopExporterFactory(),
	googlecloudexporter.NewFactory(),
}

var defaultExtensions = []component.ExtensionFactory{
	bearertokenauthextension.NewFactory(),
	healthcheckextension.NewFactory(),
	oidcauthextension.NewFactory(),
	pprofextension.NewFactory(),
	zpagesextension.NewFactory(),
	filestorage.NewFactory(),
	ballastextension.NewFactory(),
	componenttest.NewNopExtensionFactory(),
}

// DefaultFactories returns the default factories used by the observIQ collector
func DefaultFactories() (component.Factories, error) {
	return combineFactories(defaultReceivers, defaultProcessors, defaultExporters, defaultExtensions)
}

// combineFactories combines the supplied factories into a single Factories struct.
// Any errors encountered will also be combined into a single error.
func combineFactories(receivers []component.ReceiverFactory, processors []component.ProcessorFactory, exporters []component.ExporterFactory, extensions []component.ExtensionFactory) (component.Factories, error) {
	var errs []error

	receiverMap, err := component.MakeReceiverFactoryMap(receivers...)
	if err != nil {
		errs = append(errs, err)
	}

	processorMap, err := component.MakeProcessorFactoryMap(processors...)
	if err != nil {
		errs = append(errs, err)
	}

	exporterMap, err := component.MakeExporterFactoryMap(exporters...)
	if err != nil {
		errs = append(errs, err)
	}

	extensionMap, err := component.MakeExtensionFactoryMap(extensions...)
	if err != nil {
		errs = append(errs, err)
	}

	return component.Factories{
		Receivers:  receiverMap,
		Processors: processorMap,
		Exporters:  exporterMap,
		Extensions: extensionMap,
	}, multierr.Combine(errs...)
}
