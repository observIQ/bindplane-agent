package main

import (
	"github.com/observIQ/observiq-collector/extension/orphandetectorextension"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/observiqexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/httpforwarder"
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

// Get the factories for components we want to use.
// This includes all the defaults.
func components() (component.Factories, error) {
	var errs []error

	receiverMap, err := component.MakeReceiverFactoryMap(
		otlpreceiver.NewFactory(),
		filelogreceiver.NewFactory(),
		syslogreceiver.NewFactory(),
		tcplogreceiver.NewFactory(),
		udplogreceiver.NewFactory(),
	)

	if err != nil {
		errs = append(errs, err)
	}

	processorMap, err := component.MakeProcessorFactoryMap(
		groupbyattrsprocessor.NewFactory(),
		k8sprocessor.NewFactory(),
		attributesprocessor.NewFactory(),
		resourceprocessor.NewFactory(),
		batchprocessor.NewFactory(),
		memorylimiter.NewFactory(),
		probabilisticsamplerprocessor.NewFactory(),
	)

	if err != nil {
		errs = append(errs, err)
	}

	exporterMap, err := component.MakeExporterFactoryMap(
		fileexporter.NewFactory(),
		otlpexporter.NewFactory(),
		otlphttpexporter.NewFactory(),
		observiqexporter.NewFactory(),
		loggingexporter.NewFactory(),
	)

	if err != nil {
		errs = append(errs, err)
	}

	extensionMap, err := component.MakeExtensionFactoryMap(
		bearertokenauthextension.NewFactory(),
		healthcheckextension.NewFactory(),
		oidcauthextension.NewFactory(),
		pprofextension.NewFactory(),
		zpagesextension.NewFactory(),
		filestorage.NewFactory(),
		httpforwarder.NewFactory(),
		orphandetectorextension.NewFactory(),
	)

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
