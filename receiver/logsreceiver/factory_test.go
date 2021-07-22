// Copyright The OpenTelemetry Authors
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

package logsreceiver

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.uber.org/zap"
)

// TODO Test Plugins

func TestCreateReceiver(t *testing.T) {
	params := component.ReceiverCreateSettings{
		Logger: zap.NewNop(),
	}

	t.Run("Success", func(t *testing.T) {
		factory := NewFactory()
		cfg := factory.CreateDefaultConfig().(*Config)
		cfg.Pipeline = []map[string]interface{}{
			{
				"type": "json_parser", // TODO pipeline should require an input operator
			},
		}
		receiver, err := factory.CreateLogsReceiver(context.Background(), params, cfg, consumertest.NewNop())
		require.NoError(t, err, "receiver creation failed")
		require.NotNil(t, receiver, "receiver creation failed")
	})

	t.Run("Success with ConverterConfig", func(t *testing.T) {
		factory := NewFactory()
		cfg := factory.CreateDefaultConfig().(*Config)
		cfg.Pipeline = []map[string]interface{}{
			{
				"type": "json_parser", // TODO pipeline should require an input operator
			},
		}
		cfg.Converter = ConverterConfig{
			MaxFlushCount: 1,
			FlushInterval: 3 * time.Second,
		}
		receiver, err := factory.CreateLogsReceiver(context.Background(), params, cfg, consumertest.NewNop())
		require.NoError(t, err, "receiver creation failed")
		require.NotNil(t, receiver, "receiver creation failed")
	})

	t.Run("DecodeOperatorConfigsFailureMissingFields", func(t *testing.T) {
		factory := NewFactory()
		badCfg := factory.CreateDefaultConfig().(*Config)
		badCfg.Pipeline = []map[string]interface{}{
			{
				"badparam": "badvalue",
			},
		}
		receiver, err := factory.CreateLogsReceiver(context.Background(), params, badCfg, consumertest.NewNop())
		require.Error(t, err, "receiver creation should fail if parser configs aren't valid")
		require.Nil(t, receiver, "receiver creation should fail if parser configs aren't valid")
	})
}
