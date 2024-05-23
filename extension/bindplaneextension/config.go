// Copyright observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bindplaneextension

import "go.opentelemetry.io/collector/component"

// Config is the configuration for the bindplane extension
type Config struct {
	// Labels in "k1=v1,k2=v2" format
	Labels string `mapstructure:"labels"`
	// Component ID of the opamp extension. If not specified, then
	// this extension will not generate any custom messages for throughput metrics.
	OpAMP component.ID `mapstructure:"opamp"`
}
