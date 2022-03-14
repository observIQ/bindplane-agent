package pluginreceiver

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/configunmarshaler"
)

// ConfigProvider implements the service.ConfigProvider interface
type ConfigProvider struct {
	configMap    *config.Map
	errChan      chan error
	unmarshaller configunmarshaler.ConfigUnmarshaler
}

// createConfigProvider creates a config provider
func createConfigProvider(configMap *config.Map) *ConfigProvider {
	return &ConfigProvider{
		configMap:    configMap,
		errChan:      make(chan error),
		unmarshaller: configunmarshaler.NewDefault(),
	}
}

// Get returns the underlying config of the provider
func (c *ConfigProvider) Get(_ context.Context, factories component.Factories) (*config.Config, error) {
	config, err := c.unmarshaller.Unmarshal(c.configMap, factories)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}

// Watch returns a channel that indicates updates to the config
func (c *ConfigProvider) Watch() <-chan error {
	return c.errChan
}

// Shutdown always returns nil and is a no-op
func (c *ConfigProvider) Shutdown(_ context.Context) error {
	return nil
}

// FactoryProvider finds and provides factories for a config
type FactoryProvider struct {
	factories component.Factories
}

// createFactoryProvider creates a factory provider
func createFactoryProvider(factories component.Factories) *FactoryProvider {
	return &FactoryProvider{
		factories: factories,
	}
}

// GetFactories finds and returns all factories required for the supplied config.
// The provider first searches its own list of factories, before then searching the host component.
// An error is returned if the factory does not exist in either location.
func (f *FactoryProvider) GetFactories(host component.Host, configMap *config.Map) (*component.Factories, error) {
	var componentMap ComponentMap
	if err := configMap.Unmarshal(&componentMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal component map: %w", err)
	}

	receiverFactories, err := f.getReceiverFactories(host, &componentMap)
	if err != nil {
		return nil, fmt.Errorf("failed to get receiver factories: %w", err)
	}

	processorFactories, err := f.getProcessorFactories(host, &componentMap)
	if err != nil {
		return nil, fmt.Errorf("failed to get processor factories: %w", err)
	}

	exporterFactories, err := f.getExporterFactories(host, &componentMap)
	if err != nil {
		return nil, fmt.Errorf("failed to get exporter factories: %w", err)
	}

	extensionFactories, err := f.getExtensionFactories(host, &componentMap)
	if err != nil {
		return nil, fmt.Errorf("failed to get extension factories: %w", err)
	}

	return &component.Factories{
		Receivers:  receiverFactories,
		Processors: processorFactories,
		Exporters:  exporterFactories,
		Extensions: extensionFactories,
	}, nil
}

// getReceiverFactories returns the receiver factories required for the supplied config
func (f *FactoryProvider) getReceiverFactories(host component.Host, componentMap *ComponentMap) (map[config.Type]component.ReceiverFactory, error) {
	factories := map[config.Type]component.ReceiverFactory{}
	for key, factory := range f.factories.Receivers {
		factories[key] = factory
	}

	for receiver := range componentMap.Receivers {
		receiverType := receiver.Type()
		if _, ok := factories[receiverType]; ok {
			continue
		}

		factory := host.GetFactory(component.KindReceiver, receiverType)
		receiverFactory, ok := factory.(component.ReceiverFactory)
		if !ok {
			return nil, fmt.Errorf("%s factory does not exist", receiverType)
		}

		factories[receiverType] = receiverFactory
	}

	return factories, nil
}

// getProcessorFactories returns the processor factories required for the supplied config
func (f *FactoryProvider) getProcessorFactories(host component.Host, componentMap *ComponentMap) (map[config.Type]component.ProcessorFactory, error) {
	factories := map[config.Type]component.ProcessorFactory{}
	for key, factory := range f.factories.Processors {
		factories[key] = factory
	}

	for processor := range componentMap.Processors {
		processorType := processor.Type()
		if _, ok := factories[processorType]; ok {
			continue
		}

		factory := host.GetFactory(component.KindProcessor, processorType)
		processorFactory, ok := factory.(component.ProcessorFactory)
		if !ok {
			return nil, fmt.Errorf("%s factory does not exist", processorType)
		}

		factories[processorType] = processorFactory
	}

	return factories, nil
}

// getExtensionFactories returns the extension factories required for the supplied config
func (f *FactoryProvider) getExtensionFactories(host component.Host, componentMap *ComponentMap) (map[config.Type]component.ExtensionFactory, error) {
	factories := map[config.Type]component.ExtensionFactory{}
	for key, factory := range f.factories.Extensions {
		factories[key] = factory
	}

	for extension := range componentMap.Extensions {
		extensionType := extension.Type()
		if _, ok := factories[extensionType]; ok {
			continue
		}

		factory := host.GetFactory(component.KindExtension, extensionType)
		extensionFactory, ok := factory.(component.ExtensionFactory)
		if !ok {
			return nil, fmt.Errorf("%s factory does not exist", extensionType)
		}

		factories[extensionType] = extensionFactory
	}

	return factories, nil
}

// getExporterFactories returns the exporter factories required for the supplied config
func (f *FactoryProvider) getExporterFactories(host component.Host, componentMap *ComponentMap) (map[config.Type]component.ExporterFactory, error) {
	factories := map[config.Type]component.ExporterFactory{}
	for key, factory := range f.factories.Exporters {
		factories[key] = factory
	}

	for exporter := range componentMap.Exporters {
		exporterType := exporter.Type()
		if _, ok := factories[exporterType]; ok {
			continue
		}

		factory := host.GetFactory(component.KindExporter, exporterType)
		exporterFactory, ok := factory.(component.ExporterFactory)
		if !ok {
			return nil, fmt.Errorf("%s factory does not exist", exporterType)
		}

		factories[exporterType] = exporterFactory
	}

	return factories, nil
}

// ComponentMap is a map of configured open telemetry components
type ComponentMap struct {
	Receivers  map[config.ComponentID]map[string]interface{} `mapstructure:"receivers"`
	Processors map[config.ComponentID]map[string]interface{} `mapstructure:"processors"`
	Exporters  map[config.ComponentID]map[string]interface{} `mapstructure:"exporters"`
	Extensions map[config.ComponentID]map[string]interface{} `mapstructure:"extensions"`
}
