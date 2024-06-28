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

//go:build aix

package factories

import (
	"github.com/observiq/bindplane-agent/receiver/pluginreceiver"
	"github.com/observiq/bindplane-agent/receiver/sapnetweaverreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/collectdreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filelogreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/jmxreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/saphanareceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/statsdreceiver"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
	"go.opentelemetry.io/collector/receiver/receivertest"
)

// Restricted list. Many receivers aren't supported on AIX.
// Additionally, we're only including those that make sense for
// An "unmanaged" config.
var defaultReceivers = []receiver.Factory{
	collectdreceiver.NewFactory(),
	filelogreceiver.NewFactory(),
	hostmetricsreceiver.NewFactory(),
	jmxreceiver.NewFactory(),
	otlpreceiver.NewFactory(),
	pluginreceiver.NewFactory(),
	receivertest.NewNopFactory(),
	saphanareceiver.NewFactory(),
	sapnetweaverreceiver.NewFactory(),
	statsdreceiver.NewFactory(),
}
