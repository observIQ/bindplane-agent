// Copyright  observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package factories provides factories for components in the collector
package factories

import (
	"github.com/observiq/observiq-otel-collector/internal/throughputwrapper"
	"go.opentelemetry.io/collector/component"
	"go.uber.org/multierr"
)

// DefaultFactories returns the default factories used by the observIQ Distro for OpenTelemetry Collector
func DefaultFactories() (component.Factories, error) {
	return combineFactories(defaultReceivers, defaultProcessors, defaultExporters, defaultExtensions)
}

// combineFactories combines the supplied factories into a single Factories struct.
// Any errors encountered will also be combined into a single error.
func combineFactories(receivers []component.ReceiverFactory, processors []component.ProcessorFactory, exporters []component.ExporterFactory, extensions []component.ExtensionFactory) (component.Factories, error) {
	var errs []error

	receiverMap, err := component.MakeReceiverFactoryMap(wrapReceivers(receivers)...)
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

func wrapReceivers(receivers []component.ReceiverFactory) []component.ReceiverFactory {
	wrappedReceivers := make([]component.ReceiverFactory, len(defaultReceivers))

	for i, recv := range receivers {
		wrappedReceivers[i] = throughputwrapper.WrapReceiverFactory(recv)
	}

	return wrappedReceivers
}
