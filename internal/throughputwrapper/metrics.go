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

package throughputwrapper

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

const tagComponentKey = "component"

var (
	componentTagKey      = tag.MustNewKey(tagComponentKey)
	logThroughputSize    = stats.Int64("log_throughput_size", "Size of the log package emitted from the component", stats.UnitBytes)
	metricThroughputSize = stats.Int64("metric_throughput_size", "Size of the metric package emitted from the component", stats.UnitBytes)
	traceThroughputSize  = stats.Int64("trace_throughput_size", "Size of the trace package emitted from the component", stats.UnitBytes)
)

func metricViews() []*view.View {
	componentTagKeys := []tag.Key{componentTagKey}

	return []*view.View{
		{
			Name:        tagComponentKey + "/" + logThroughputSize.Name(),
			Description: logThroughputSize.Description(),
			Measure:     logThroughputSize,
			TagKeys:     componentTagKeys,
			Aggregation: view.Sum(),
		},
		{
			Name:        tagComponentKey + "/" + metricThroughputSize.Name(),
			Description: metricThroughputSize.Description(),
			Measure:     metricThroughputSize,
			TagKeys:     componentTagKeys,
			Aggregation: view.Sum(),
		},
		{
			Name:        tagComponentKey + "/" + traceThroughputSize.Name(),
			Description: traceThroughputSize.Description(),
			Measure:     traceThroughputSize,
			TagKeys:     componentTagKeys,
			Aggregation: view.Sum(),
		},
	}
}
