// Copyright observIQ, Inc.
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

package oktareceiver

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
)

type oktaLogsReceiver struct {
	cfg      Config
	consumer consumer.Logs
}

// newOktaLogsReceiver returns a newly configured oktaLogsReceiver
func newOktaLogsReceiver(cfg *Config, consumer consumer.Logs) (*oktaLogsReceiver, error) {
	return &oktaLogsReceiver{
		cfg:      *cfg,
		consumer: consumer,
	}, nil
}

func (r *oktaLogsReceiver) Start(ctx context.Context, host component.Host) error {
	return nil
}

func (r *oktaLogsReceiver) Shutdown(ctx context.Context) error {
	return nil
}
