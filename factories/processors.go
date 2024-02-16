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

package factories

import (
	"github.com/observiq/bindplane-agent/internal/processor/snapshotprocessor"
	"github.com/observiq/bindplane-agent/processor/datapointcountprocessor"
	"github.com/observiq/bindplane-agent/processor/logcountprocessor"
	"github.com/observiq/bindplane-agent/processor/logdeduplicationprocessor"
	"github.com/observiq/bindplane-agent/processor/lookupprocessor"
	"github.com/observiq/bindplane-agent/processor/maskprocessor"
	"github.com/observiq/bindplane-agent/processor/metricextractprocessor"
	"github.com/observiq/bindplane-agent/processor/metricstatsprocessor"
	"github.com/observiq/bindplane-agent/processor/removeemptyvaluesprocessor"
	"github.com/observiq/bindplane-agent/processor/resourceattributetransposerprocessor"
	"github.com/observiq/bindplane-agent/processor/samplingprocessor"
	"github.com/observiq/bindplane-agent/processor/spancountprocessor"
	"github.com/observiq/bindplane-agent/processor/throughputmeasurementprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/cumulativetodeltaprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/deltatorateprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/filterprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbyattrsprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbytraceprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/logstransformprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricsgenerationprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricstransformprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/probabilisticsamplerprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourceprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/routingprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/spanmetricsprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/spanprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/tailsamplingprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/batchprocessor"
	"go.opentelemetry.io/collector/processor/memorylimiterprocessor"
	"go.opentelemetry.io/collector/processor/processortest"
)

var defaultProcessors = []processor.Factory{
	attributesprocessor.NewFactory(),
	batchprocessor.NewFactory(),
	processortest.NewNopFactory(),
	cumulativetodeltaprocessor.NewFactory(),
	datapointcountprocessor.NewFactory(),
	deltatorateprocessor.NewFactory(),
	filterprocessor.NewFactory(),
	groupbyattrsprocessor.NewFactory(),
	groupbytraceprocessor.NewFactory(),
	k8sattributesprocessor.NewFactory(),
	logcountprocessor.NewFactory(),
	logdeduplicationprocessor.NewFactory(),
	logstransformprocessor.NewFactory(),
	lookupprocessor.NewFactory(),
	maskprocessor.NewFactory(),
	memorylimiterprocessor.NewFactory(),
	metricextractprocessor.NewFactory(),
	metricsgenerationprocessor.NewFactory(),
	metricstatsprocessor.NewFactory(),
	metricstransformprocessor.NewFactory(),
	probabilisticsamplerprocessor.NewFactory(),
	removeemptyvaluesprocessor.NewFactory(),
	resourceattributetransposerprocessor.NewFactory(),
	resourcedetectionprocessor.NewFactory(),
	resourceprocessor.NewFactory(),
	routingprocessor.NewFactory(),
	samplingprocessor.NewFactory(),
	snapshotprocessor.NewFactory(),
	spancountprocessor.NewFactory(),
	spanmetricsprocessor.NewFactory(),
	spanprocessor.NewFactory(),
	throughputmeasurementprocessor.NewFactory(),
	tailsamplingprocessor.NewFactory(),
	transformprocessor.NewFactory(),
}
