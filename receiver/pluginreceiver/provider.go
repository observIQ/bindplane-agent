package pluginreceiver

import (
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/configunmarshaler"
	"gopkg.in/yaml.v2"
)

// ConfigProvider implements the service.ConfigProvider interface
type ConfigProvider struct {
	components   *ComponentMap
	errChan      chan error
	unmarshaller configunmarshaler.ConfigUnmarshaler
}

// createConfigProvider creates a config provider
func createConfigProvider(components *ComponentMap) *ConfigProvider {
	return &ConfigProvider{
		components:   components,
		errChan:      make(chan error),
		unmarshaller: configunmarshaler.NewDefault(),
	}
}

// Get returns the underlying config of the provider
func (c *ConfigProvider) Get(_ context.Context, factories component.Factories) (*config.Config, error) {
	configMap := c.components.ToConfigMap()
	config, err := c.unmarshaller.Unmarshal(configMap, factories)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config map: %w", err)
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

// GetRequiredFactories returns the factories required for the configured components
func (c *ConfigProvider) GetRequiredFactories(host component.Host, emitterFactory component.ExporterFactory) (*component.Factories, error) {
	receiverFactories, err := c.getReceiverFactories(host)
	if err != nil {
		return nil, fmt.Errorf("failed to get receiver factories: %w", err)
	}

	processorFactories, err := c.getProcessorFactories(host)
	if err != nil {
		return nil, fmt.Errorf("failed to get processor factories: %w", err)
	}

	extensionFactories, err := c.getExtensionFactories(host)
	if err != nil {
		return nil, fmt.Errorf("failed to get extension factories: %w", err)
	}

	exporterFactories := map[config.Type]component.ExporterFactory{
		emitterFactory.Type(): emitterFactory,
	}

	return &component.Factories{
		Receivers:  receiverFactories,
		Processors: processorFactories,
		Exporters:  exporterFactories,
		Extensions: extensionFactories,
	}, nil
}

// getReceiverFactories returns the receiver factories required for the configured receivers
func (c *ConfigProvider) getReceiverFactories(host component.Host) (map[config.Type]component.ReceiverFactory, error) {
	factories := map[config.Type]component.ReceiverFactory{}
	for receiver := range c.components.Receivers {
		receiverType := receiver.Type()
		if _, ok := factories[receiverType]; ok {
			continue
		}

		factory := host.GetFactory(component.KindReceiver, receiverType)
		receiverFactory, ok := factory.(component.ReceiverFactory)
		if !ok {
			return nil, fmt.Errorf("receiver factory %s is missing from host", receiverType)
		}

		factories[receiverType] = receiverFactory
	}

	return factories, nil
}

// getProcessorFactories returns the processor factories required for the configured processors
func (c *ConfigProvider) getProcessorFactories(host component.Host) (map[config.Type]component.ProcessorFactory, error) {
	factories := map[config.Type]component.ProcessorFactory{}
	for processor := range c.components.Processors {
		processorType := processor.Type()
		if _, ok := factories[processorType]; ok {
			continue
		}

		factory := host.GetFactory(component.KindProcessor, processorType)
		processorFactory, ok := factory.(component.ProcessorFactory)
		if !ok {
			return nil, fmt.Errorf("processor factory %s is missing from host", processorType)
		}

		factories[processorType] = processorFactory
	}

	return factories, nil
}

// getExtensionFactories returns the extension factories required for the configured extensions
func (c *ConfigProvider) getExtensionFactories(host component.Host) (map[config.Type]component.ExtensionFactory, error) {
	factories := map[config.Type]component.ExtensionFactory{}
	for extension := range c.components.Extensions {
		extensionType := extension.Type()
		if _, ok := factories[extensionType]; ok {
			continue
		}

		factory := host.GetFactory(component.KindExtension, extensionType)
		extensionFactory, ok := factory.(component.ExtensionFactory)
		if !ok {
			return nil, fmt.Errorf("extension factory %s is missing from host", extensionType)
		}

		factories[extensionType] = extensionFactory
	}

	return factories, nil
}

// unmarshalComponents unmarshals a component map from yaml
func unmarshalComponentMap(bytes []byte) (*ComponentMap, error) {
	var components ComponentMap
	if err := yaml.Unmarshal(bytes, &components); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	emitterID := config.NewComponentID(emitterTypeStr)
	components.Exporters = map[config.ComponentID]map[string]interface{}{
		emitterID: nil,
	}

	for key, pipeline := range components.Service.Pipelines {
		pipeline.Exporters = []config.ComponentID{emitterID}
		components.Service.Pipelines[key] = pipeline
	}

	return &components, nil
}

// ComponentMap is a map of configured open telemetry components
type ComponentMap struct {
	Receivers  map[config.ComponentID]map[string]interface{} `yaml:"receivers" mapstructure:"receivers"`
	Processors map[config.ComponentID]map[string]interface{} `yaml:"processors" mapstructure:"processors"`
	Exporters  map[config.ComponentID]map[string]interface{} `yaml:"exporters" mapstructure:"exporters"`
	Extensions map[config.ComponentID]map[string]interface{} `yaml:"extensions" mapstructure:"extensions"`
	Service    ServiceMap                                    `yaml:"service" mapstructure:"service"`
}

// ToConfigMap returns the component map as a config map
func (c *ComponentMap) ToConfigMap() *config.Map {
	var mapString map[string]interface{}
	_ = mapstructure.Decode(c, &mapString)
	return config.NewMapFromStringMap(mapString)
}

// ServiceMap is a map of service components
type ServiceMap struct {
	Extensions []config.ComponentID `yaml:"extensions" mapstructure:"extensions"`
	Pipelines  map[string]Pipeline  `yaml:"pipelines" mapstructure:"pipelines"`
}

// Pipeline is a component pipeline
type Pipeline struct {
	Receivers  []config.ComponentID `yaml:"receivers" mapstructure:"receivers"`
	Processors []config.ComponentID `yaml:"processors" mapstructure:"processors"`
	Exporters  []config.ComponentID `yaml:"exporters" mapstructure:"exporters"`
}
