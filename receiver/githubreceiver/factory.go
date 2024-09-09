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

package githubreceiver // import "github.com/observiq/bindplane-agent/receiver/githubreceiver"

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"

	"github.com/observiq/bindplane-agent/receiver/githubreceiver/internal/metadata"
)

// NewFactory creates a factory for GitHub receiver.
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		metadata.Type,
		createDefaultConfig,
		receiver.WithLogs(createLogsReceiver, metadata.LogsStability),
	)
}
func createDefaultConfig() component.Config {
	c := &Config{
		PollInterval:  time.Second,
		WebhookConfig: nil,
	}
	return c
}

func createLogsReceiver(
	_ context.Context,
	params receiver.Settings,
	rConf component.Config,
	consumer consumer.Logs,
) (receiver.Logs, error) {
	cfg := rConf.(*Config)
	rcvr, err := newGitHubLogsReceiver(cfg, params.Logger, consumer)
	if err != nil {
		return nil, fmt.Errorf("unable to create an log receiver instance: %w", err)
	}
	return rcvr, nil
}
