// Copyright  The OpenTelemetry Authors
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

package varnishreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/varnishreceiver"

import "go.opentelemetry.io/collector/receiver/scraperhelper" // Config defines the configuration for the various elements of the receiver agent.

// Config defines configuration for varnish metrics receiver.
type Config struct {
	scraperhelper.ScraperControllerSettings `mapstructure:",squash"`
	WorkingDir                              string `mapstructure:"working_dir"`
	ExecDir                                 string `mapstructure:"exec_dir"`
}
