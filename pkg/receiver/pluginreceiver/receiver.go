package pluginreceiver

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"text/template"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configunmarshaler"
	"go.opentelemetry.io/collector/service"
	"go.opentelemetry.io/contrib/zpages"
)

// Receiver is a receiver for running templated plugins
type Receiver struct {
	template   string
	parameters map[string]interface{}
	consumer   Consumer
	telemetry  component.TelemetrySettings
	buildInfo  component.BuildInfo
	svc        *service.Service
}

// Start will start the receiver
func (r *Receiver) Start(ctx context.Context, host component.Host) error {
	bytes, err := r.renderTemplate()
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	pipeline, err := createPipeline(bytes)
	if err != nil {
		return fmt.Errorf("failed to create pipeline: %w", err)
	}

	factories, err := pipeline.getRequiredFactories(host, r.consumer)
	if err != nil {
		return fmt.Errorf("failed to get required factories: %w", err)
	}

	unmarshaller := configunmarshaler.NewDefault()
	config, err := unmarshaller.Unmarshal(pipeline.Map, factories)
	if err != nil {
		return fmt.Errorf("failed to unmarshal rendered template: %w", err)
	}

	settings := &service.SvcSettings{
		Factories:           factories,
		BuildInfo:           r.buildInfo,
		Config:              config,
		Telemetry:           r.telemetry,
		ZPagesSpanProcessor: zpages.NewSpanProcessor(),
		AsyncErrorChannel:   make(chan error),
	}

	svc, err := service.NewService(settings)
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}
	r.svc = svc

	if err = svc.Start(ctx); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	return nil
}

// Shutdown will shutdown the receiver
func (r *Receiver) Shutdown(ctx context.Context) error {
	if r.svc == nil {
		return errors.New("service not initialized")
	}

	err := r.svc.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("failed to shutdown service: %w", err)
	}

	return nil
}

// renderTemplate renders the receiver's templated config
func (r *Receiver) renderTemplate() ([]byte, error) {
	templateContents, err := ioutil.ReadFile(r.template)
	if err != nil {
		return nil, fmt.Errorf("failed to read template from path: %w", err)
	}

	template, err := template.New(r.template).Parse(string(templateContents))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	var writer bytes.Buffer
	if err := template.Execute(&writer, r.parameters); err != nil {
		return nil, fmt.Errorf("failed to execute template with parameters: %w", err)
	}

	return writer.Bytes(), nil
}
