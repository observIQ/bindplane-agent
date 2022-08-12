package throughputmeasurementprocessor

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opentelemetry.io/collector/obsreport"
)

var (
	logDataSize    = stats.Int64("log_data_size", "Size of the log package passed to the processor", stats.UnitBytes)
	metricDataSize = stats.Int64("metric_data_size", "Size of the metric package passed to the processor", stats.UnitBytes)
	traceDataSize  = stats.Int64("trace_data_size", "Size of the trace package passed to the processor", stats.UnitBytes)
)

func MetricViews() []*view.View {
	return []*view.View{
		{
			Name:        obsreport.BuildProcessorCustomMetricName(string(typeStr), logDataSize.Name()),
			Description: logDataSize.Description(),
			Measure:     logDataSize,
			Aggregation: view.Sum(),
		},
		{
			Name:        obsreport.BuildProcessorCustomMetricName(string(typeStr), metricDataSize.Name()),
			Description: metricDataSize.Description(),
			Measure:     metricDataSize,
			Aggregation: view.Sum(),
		},
		{
			Name:        obsreport.BuildProcessorCustomMetricName(string(typeStr), traceDataSize.Name()),
			Description: traceDataSize.Description(),
			Measure:     traceDataSize,
			Aggregation: view.Sum(),
		},
	}
}
