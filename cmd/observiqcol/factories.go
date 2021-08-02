package main

import (
	"github.com/observiq/observiq-collector/extension/orphandetectorextension"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/observiqexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/storage/filestorage"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbyattrsprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filelogreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/syslogreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/tcplogreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/udplogreceiver"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/consumererror"
	"go.opentelemetry.io/collector/exporter/fileexporter"
	"go.opentelemetry.io/collector/exporter/loggingexporter"
	"go.opentelemetry.io/collector/exporter/otlpexporter"
	"go.opentelemetry.io/collector/exporter/otlphttpexporter"
	"go.opentelemetry.io/collector/extension/bearertokenauthextension"
	"go.opentelemetry.io/collector/extension/healthcheckextension"
	"go.opentelemetry.io/collector/extension/oidcauthextension"
	"go.opentelemetry.io/collector/extension/pprofextension"
	"go.opentelemetry.io/collector/extension/zpagesextension"
	"go.opentelemetry.io/collector/processor/attributesprocessor"
	"go.opentelemetry.io/collector/processor/batchprocessor"
	"go.opentelemetry.io/collector/processor/memorylimiter"
	"go.opentelemetry.io/collector/processor/probabilisticsamplerprocessor"
	"go.opentelemetry.io/collector/processor/resourceprocessor"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
)

var defaultReceivers = []component.ReceiverFactory{
	otlpreceiver.NewFactory(),
	filelogreceiver.NewFactory(),
	syslogreceiver.NewFactory(),
	tcplogreceiver.NewFactory(),
	udplogreceiver.NewFactory(),
}

var defaultProcessors = []component.ProcessorFactory{
	groupbyattrsprocessor.NewFactory(),
	k8sprocessor.NewFactory(),
	attributesprocessor.NewFactory(),
	resourceprocessor.NewFactory(),
	batchprocessor.NewFactory(),
	memorylimiter.NewFactory(),
	probabilisticsamplerprocessor.NewFactory(),
}

var defaultExporters = []component.ExporterFactory{
	fileexporter.NewFactory(),
	otlpexporter.NewFactory(),
	otlphttpexporter.NewFactory(),
	observiqexporter.NewFactory(),
	loggingexporter.NewFactory(),
}

var defaultExtensions = []component.ExtensionFactory{
	bearertokenauthextension.NewFactory(),
	healthcheckextension.NewFactory(),
	oidcauthextension.NewFactory(),
	pprofextension.NewFactory(),
	zpagesextension.NewFactory(),
	filestorage.NewFactory(),
	orphandetectorextension.NewFactory(),
}

// defaultFactories returns the default factories used by the observIQ collector
func defaultFactories() (component.Factories, error) {
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
	}, consumererror.Combine(errs)
}
