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
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

func Test_newLogConsumer(t *testing.T) {
	nopLogger := zap.NewNop()
	componentID := "id"
	baseConsumer := consumertest.NewNop()
	lConsumer := newLogConsumer(nopLogger, componentID, baseConsumer)

	require.Equal(t, nopLogger, lConsumer.logger)
	require.Equal(t, baseConsumer, lConsumer.baseConsumer)
	require.Len(t, lConsumer.mutators, 1)
	require.Equal(t, &plog.ProtoMarshaler{}, lConsumer.logsSizer)
}

func Test_logConsumer_ConsumeLogs(t *testing.T) {
	nopLogger := zap.NewNop()
	componentID := "id"
	baseConsumer := new(consumertest.LogsSink)
	lConsumer := newLogConsumer(nopLogger, componentID, baseConsumer)

	ld := plog.NewLogs()
	ld.ResourceLogs().AppendEmpty().ScopeLogs().AppendEmpty().LogRecords().AppendEmpty()

	err := lConsumer.ConsumeLogs(context.Background(), ld)
	require.NoError(t, err)

	require.Equal(t, 1, baseConsumer.LogRecordCount())
}

func Test_logConsumer_Capabilities(t *testing.T) {
	nopLogger := zap.NewNop()
	componentID := "id"
	baseConsumer := consumertest.NewNop()
	lConsumer := newLogConsumer(nopLogger, componentID, baseConsumer)

	require.Equal(t, baseConsumer.Capabilities(), lConsumer.Capabilities())
}
