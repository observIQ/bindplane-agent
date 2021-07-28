package main

import (
	"errors"
	"testing"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/observiqexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/tcplogreceiver"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension/bearertokenauthextension"
	"go.opentelemetry.io/collector/processor/attributesprocessor"
)

func TestCombineFactories(t *testing.T) {
	testCases := []struct {
		name          string
		receivers     []component.ReceiverFactory
		processors    []component.ProcessorFactory
		exporters     []component.ExporterFactory
		extensions    []component.ExtensionFactory
		expectedError error
	}{
		{
			name:       "With valid combination",
			receivers:  defaultReceivers,
			processors: defaultProcessors,
			exporters:  defaultExporters,
			extensions: defaultExtensions,
		},
		{
			name: "With single error",
			receivers: []component.ReceiverFactory{
				tcplogreceiver.NewFactory(),
				tcplogreceiver.NewFactory(),
			},
			expectedError: errors.New(`duplicate receiver factory "tcplog"`),
		},
		{
			name: "With multiple errors",
			receivers: []component.ReceiverFactory{
				tcplogreceiver.NewFactory(),
				tcplogreceiver.NewFactory(),
			},
			processors: []component.ProcessorFactory{
				attributesprocessor.NewFactory(),
				attributesprocessor.NewFactory(),
			},
			exporters: []component.ExporterFactory{
				observiqexporter.NewFactory(),
				observiqexporter.NewFactory(),
			},
			extensions: []component.ExtensionFactory{
				bearertokenauthextension.NewFactory(),
				bearertokenauthextension.NewFactory(),
			},
			expectedError: errors.New(`[duplicate receiver factory "tcplog"; duplicate processor factory "attributes"; duplicate exporter factory "observiq"; duplicate extension factory "bearertokenauth"]`),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			factories, err := combineFactories(tc.receivers, tc.processors, tc.exporters, tc.extensions)

			if tc.expectedError != nil {
				assert.Error(t, tc.expectedError, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
				return
			}

			assert.NoError(t, err)

			for _, receiver := range tc.receivers {
				assert.Equal(t, factories.Receivers[receiver.Type()], receiver)
			}

			for _, processor := range tc.processors {
				assert.Equal(t, factories.Processors[processor.Type()], processor)
			}

			for _, exporter := range tc.exporters {
				assert.Equal(t, factories.Exporters[exporter.Type()], exporter)
			}

			for _, extension := range tc.extensions {
				assert.Equal(t, factories.Extensions[extension.Type()], extension)
			}
		})
	}
}

func TestDefaultFactories(t *testing.T) {
	factories, err := defaultFactories()
	assert.NoError(t, err)

	for _, receiver := range defaultReceivers {
		assert.Equal(t, factories.Receivers[receiver.Type()], receiver)
	}

	for _, processor := range defaultProcessors {
		assert.Equal(t, factories.Processors[processor.Type()], processor)
	}

	for _, exporter := range defaultExporters {
		assert.Equal(t, factories.Exporters[exporter.Type()], exporter)
	}

	for _, extension := range defaultExtensions {
		assert.Equal(t, factories.Extensions[extension.Type()], extension)
	}
}
