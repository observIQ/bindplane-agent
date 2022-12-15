package routerprocessor

import (
	"context"
	"fmt"

	"github.com/observiq/observiq-otel-collector/internal/expr"
	"github.com/observiq/observiq-otel-collector/receiver/routereceiver"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

const (
	// defaultRoute is the name of the default route.
	defaultRoute = "default"
)

// routerProcessor is a processor that routes OTel objects to route receivers based on a match expression.
type routerProcessor struct {
	routes   []*route
	consumer consumer.Logs
	logger   *zap.Logger
}

// route represents a route to a receiver with it's compiled match expression.
type route struct {
	name  string
	match *expr.Expression
}

// newProcessor creates a new router processor.
func newProcessor(config *Config, consumer consumer.Logs, logger *zap.Logger) (*routerProcessor, error) {
	routes := make([]*route, len(config.Routes))
	for i, r := range config.Routes {
		match, err := expr.CreateBoolExpression(r.Match)
		if err != nil {
			return nil, fmt.Errorf("invalid match expression '%s': %w", r.Match, err)
		}

		routes[i] = &route{
			name:  r.Route,
			match: match,
		}
	}

	return &routerProcessor{
		routes:   routes,
		consumer: consumer,
		logger:   logger,
	}, nil
}

// Start starts the processor.
func (p *routerProcessor) Start(_ context.Context, _ component.Host) error {
	return nil
}

// Capabilities returns the consumer's capabilities.
func (p *routerProcessor) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}

// Shutdown stops the processor.
func (p *routerProcessor) Shutdown(_ context.Context) error {
	return nil
}

// ConsumeLogs processes the logs.
func (p *routerProcessor) ConsumeLogs(ctx context.Context, pl plog.Logs) error {
	routeBatchMap := p.createRouteBatchMap()

	// Iterate through plogs saving off the resource
	for ri := 0; ri < pl.ResourceLogs().Len(); ri++ {
		resourceLogs := pl.ResourceLogs().At(ri)
		resourceLogsMap := p.processResourceLogs(resourceLogs)

		// Add the resource logs to the route batch
		for route, routeResourceLogs := range resourceLogsMap {
			if routeResourceLogs.ScopeLogs().Len() == 0 {
				continue
			}
			routeResourceLogs.CopyTo(routeBatchMap[route].ResourceLogs().AppendEmpty())
		}
	}

	var multiErr error
	// Route batches to their corresponding destinations.
	for route, batch := range routeBatchMap {
		switch {
		case batch.LogRecordCount() == 0: // No log records then skip
			continue
		case route == defaultRoute: // Default route process as normal
			p.consumer.ConsumeLogs(ctx, batch)
		default:
			// Route logs to the correct route receiver
			if err := routereceiver.RouteLogs(ctx, route, batch); err != nil {
				p.logger.Error("failed to route logs", zap.String("route", route), zap.Error(err))
				multiErr = multierr.Append(multiErr, err)
			}
		}
	}
	return multiErr
}

// findRoute returns the name of the route that the records matches
func (p *routerProcessor) findRoute(record expr.Record) string {
	for _, route := range p.routes {
		if route.match.MatchRecord(record) {
			return route.name
		}
	}

	return defaultRoute
}

func (p *routerProcessor) createRouteBatchMap() map[string]plog.Logs {
	// Initialize map with default route batch.
	routeMap := map[string]plog.Logs{
		defaultRoute: plog.NewLogs(),
	}

	for _, route := range p.routes {
		routeMap[route.name] = plog.NewLogs()
	}

	return routeMap
}

func (p *routerProcessor) createRouteResourceMap(baseResource pcommon.Resource) map[string]plog.ResourceLogs {
	// Create a base resource logs
	resourceLogs := plog.NewResourceLogs()
	baseResource.CopyTo(resourceLogs.Resource())

	// Initialize map with default route batch.
	resourceLogsMap := map[string]plog.ResourceLogs{
		defaultRoute: plog.NewResourceLogs(),
	}

	// Create an key for each route
	for _, route := range p.routes {
		resourceLogsMap[route.name] = plog.NewResourceLogs()
	}

	// Copy over base to each map entry
	for _, routeResourceLogs := range resourceLogsMap {
		resourceLogs.CopyTo(routeResourceLogs)
	}

	return resourceLogsMap
}

func (p *routerProcessor) createRouteScopeMap(baseScopeLogs plog.ScopeLogs) map[string]plog.ScopeLogs {
	// Create a scope logs with the same scope and schemas as the passed in base
	scopeLogs := plog.NewScopeLogs()
	baseScopeLogs.Scope().CopyTo(scopeLogs.Scope())
	scopeLogs.SetSchemaUrl(baseScopeLogs.SchemaUrl())

	// Initialize map with default route batch.
	scopeLogsMap := map[string]plog.ScopeLogs{
		defaultRoute: plog.NewScopeLogs(),
	}

	// Create an key for each route
	for _, route := range p.routes {
		scopeLogsMap[route.name] = plog.NewScopeLogs()
	}

	// Copy over base to each map entry
	for _, routeScopeLogs := range scopeLogsMap {
		scopeLogs.CopyTo(routeScopeLogs)
	}

	return scopeLogsMap
}

func (p *routerProcessor) processResourceLogs(resourceLogs plog.ResourceLogs) map[string]plog.ResourceLogs {
	resourceAttrs := resourceLogs.Resource().Attributes().AsRaw()
	resourceLogsMap := p.createRouteResourceMap(resourceLogs.Resource())

	// Iterate through scope logs
	for si := 0; si < resourceLogs.ScopeLogs().Len(); si++ {
		scopeLogs := resourceLogs.ScopeLogs().At(si)

		// Iterate through log records and route them to correct batch if the expression matches
		scopeLogsMap := p.processScopeLogs(scopeLogs, resourceAttrs)

		// Add the scope logs to the resource logs
		for route, routeScopedLogs := range scopeLogsMap {
			if routeScopedLogs.LogRecords().Len() == 0 {
				continue
			}
			routeScopedLogs.CopyTo(resourceLogsMap[route].ScopeLogs().AppendEmpty())
		}
	}

	return resourceLogsMap
}

func (p *routerProcessor) processScopeLogs(scopeLogs plog.ScopeLogs, resourceAttrs map[string]any) map[string]plog.ScopeLogs {
	scopeLogsMap := p.createRouteScopeMap(scopeLogs)

	for li := 0; li < scopeLogs.LogRecords().Len(); li++ {
		logRecord := scopeLogs.LogRecords().At(li)
		record := expr.ConvertToRecord(logRecord, resourceAttrs)
		matchingRoute := p.findRoute(record)

		routeScopeLogs := scopeLogsMap[matchingRoute]
		logRecord.CopyTo(routeScopeLogs.LogRecords().AppendEmpty())
		continue
	}

	return scopeLogsMap
}
