package pluginreceiver

import (
	"fmt"
	"strings"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"gopkg.in/yaml.v2"
)

const (
	receiversKey  = "receivers"
	processorsKey = "processors"
	exportersKey  = "exporters"
	extensionsKey = "extensions"
	pipelinesKey  = "service::pipelines"
)

// pipeline is the internal pipeline of the receiver
type pipeline struct {
	*config.Map
}

// createPipeline creates a pipeline from the supplied bytes
func createPipeline(bytes []byte) (*pipeline, error) {
	var data map[string]interface{}
	if err := yaml.Unmarshal(bytes, &data); err != nil {
		return nil, fmt.Errorf("unable to parse content as yaml: %w", err)
	}

	configMap, err := createConfigMap(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create config map: %w", err)
	}

	return &pipeline{configMap}, nil
}

// createConfigMap creates a config map from the supplied map
func createConfigMap(data map[string]interface{}) (*config.Map, error) {
	configMap := config.NewMapFromStringMap(data)
	configMap.Set(exportersKey, map[string]interface{}{
		typeStr: nil,
	})

	if !configMap.IsSet(pipelinesKey) {
		return nil, fmt.Errorf("%s subsection is not defined", pipelinesKey)
	}

	pipelines, err := configMap.Sub(pipelinesKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s subsection: %w", pipelinesKey, err)
	}

	for pipelineKey := range pipelines.ToStringMap() {
		key := fmt.Sprintf("%s::%s::%s", pipelinesKey, pipelineKey, exportersKey)
		configMap.Set(key, []string{typeStr})
	}

	return configMap, nil
}

// getComponentTypes returns a unique slice of component types found in the pipeline
func (p *pipeline) getComponentTypes(componentKind string) []config.Type {
	foundTypes := map[string]bool{}
	uniqueTypes := []config.Type{}

	if !p.IsSet(componentKind) {
		return nil
	}

	components, err := p.Sub(componentKind)
	if err != nil {
		return nil
	}

	for key := range components.ToStringMap() {
		typeStr := strings.Split(key, "/")[0]
		if _, ok := foundTypes[typeStr]; !ok {
			foundTypes[key] = true
			uniqueTypes = append(uniqueTypes, config.Type(key))
		}
	}

	return uniqueTypes
}

// getRequiredFactories returns the factories required by the pipeline
func (p *pipeline) getRequiredFactories(host component.Host, consumer Consumer) (component.Factories, error) {
	receiverFactories, err := p.getReceiverFactories(host)
	if err != nil {
		return component.Factories{}, fmt.Errorf("failed to get receiver factories: %w", err)
	}

	processorFactories, err := p.getRequiredProcessorFactories(host)
	if err != nil {
		return component.Factories{}, fmt.Errorf("failed to get processor factories: %w", err)
	}

	extensionFactories, err := p.getRequiredExtensionFactories(host)
	if err != nil {
		return component.Factories{}, fmt.Errorf("failed to get extension factories: %w", err)
	}

	exporterFactories, err := component.MakeExporterFactoryMap(consumer.createFactory())
	if err != nil {
		return component.Factories{}, fmt.Errorf("failed to get exporter factory: %w", err)
	}

	return component.Factories{
		Receivers:  receiverFactories,
		Processors: processorFactories,
		Exporters:  exporterFactories,
		Extensions: extensionFactories,
	}, nil
}

// getRequiredReceiverFactories returns the receiver factories required by the pipeline
func (p *pipeline) getReceiverFactories(host component.Host) (map[config.Type]component.ReceiverFactory, error) {
	types := p.getComponentTypes(receiversKey)
	factories := make([]component.ReceiverFactory, len(types))
	for i, t := range types {
		factory := host.GetFactory(component.KindReceiver, t)
		receiverFactory, ok := factory.(component.ReceiverFactory)
		if !ok {
			return nil, fmt.Errorf("failed to convert factory to receiver factory: %s", t)
		}

		factories[i] = receiverFactory
	}

	return component.MakeReceiverFactoryMap(factories...)
}

// getRequiredProcessorFactories returns the processor factories required by the pipeline
func (p *pipeline) getRequiredProcessorFactories(host component.Host) (map[config.Type]component.ProcessorFactory, error) {
	types := p.getComponentTypes(processorsKey)
	factories := make([]component.ProcessorFactory, len(types))
	for i, t := range types {
		factory := host.GetFactory(component.KindProcessor, t)
		processorFactory, ok := factory.(component.ProcessorFactory)
		if !ok {
			return nil, fmt.Errorf("failed to convert factory to processor factory: %s", t)
		}

		factories[i] = processorFactory
	}

	return component.MakeProcessorFactoryMap(factories...)
}

// getRequiredExtensionFactories returns the extension factories required by the pipeline
func (p *pipeline) getRequiredExtensionFactories(host component.Host) (map[config.Type]component.ExtensionFactory, error) {
	types := p.getComponentTypes(extensionsKey)
	factories := make([]component.ExtensionFactory, len(types))
	for i, t := range types {
		factory := host.GetFactory(component.KindExtension, t)
		extensionFactory, ok := factory.(component.ExtensionFactory)
		if !ok {
			return nil, fmt.Errorf("failed to convert factory to extension factory: %s", t)
		}

		factories[i] = extensionFactory
	}

	return component.MakeExtensionFactoryMap(factories...)
}
