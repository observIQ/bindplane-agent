package main

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/observiqexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/httpforwarder"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer/ecsobserver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer/hostobserver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer/k8sobserver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/storage/filestorage"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbyattrsprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filelogreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/receivercreator"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/syslogreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/tcplogreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/udplogreceiver"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/consumererror"
	"go.opentelemetry.io/collector/service/defaultcomponents"
)

// Get the factories for components we want to use.
// This includes all the defaults.
func components() (component.Factories, error) {
	defaultFactories, err := defaultcomponents.Components()

	var errs []error

	if err != nil {
		errs = append(errs, err)
	}

	receivers := []component.ReceiverFactory{
		filelogreceiver.NewFactory(),
		syslogreceiver.NewFactory(),
		tcplogreceiver.NewFactory(),
		udplogreceiver.NewFactory(),
		receivercreator.NewFactory(),
	}
	for _, rf := range defaultFactories.Receivers {
		receivers = append(receivers, rf)
	}

	processors := []component.ProcessorFactory{
		groupbyattrsprocessor.NewFactory(),
		k8sprocessor.NewFactory(),
	}
	for _, pf := range defaultFactories.Processors {
		processors = append(processors, pf)
	}

	exporters := []component.ExporterFactory{
		observiqexporter.NewFactory(), // TODO: This needs updating with a new release of the collector-contrib repo, exporter is out of date for now.
	}
	for _, ef := range defaultFactories.Exporters {
		exporters = append(exporters, ef)
	}

	//TODO: oauth2clientauthextension -- having trouble with the imports for this one, but would be nice to have.
	extensions := []component.ExtensionFactory{
		filestorage.NewFactory(),
		httpforwarder.NewFactory(),
		ecsobserver.NewFactory(),
		hostobserver.NewFactory(),
		k8sobserver.NewFactory(),
	}
	for _, ef := range defaultFactories.Extensions {
		extensions = append(extensions, ef)
	}

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
