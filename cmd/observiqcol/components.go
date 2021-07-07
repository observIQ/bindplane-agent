package main

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/storage/filestorage"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbyattrsprocessor"
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
		// TODO: Figure out dependency issues with these
		// filelogreceiver.NewFactory()
		// syslogreceiver.NewFactory(),
		// tcplogreceiver.NewFactory(),
		// udplogreceiver.NewFactory(),
	}
	for _, rf := range defaultFactories.Receivers {
		receivers = append(receivers, rf)
	}

	processors := []component.ProcessorFactory{
		groupbyattrsprocessor.NewFactory(),
		// k8sprocessor.NewFactory(),
	}
	for _, pf := range defaultFactories.Processors {
		processors = append(processors, pf)
	}

	exporters := []component.ExporterFactory{}
	for _, ef := range defaultFactories.Exporters {
		exporters = append(exporters, ef)
	}

	extensions := []component.ExtensionFactory{
		filestorage.NewFactory(),
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
