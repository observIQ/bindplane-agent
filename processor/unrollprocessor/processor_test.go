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

package unrollprocessor

import (
	"context"
	"testing"

	"go.opentelemetry.io/collector/pdata/plog"
)

func BenchmarkUnroll(b *testing.B) {
	unrollProcessor := &unrollProcessor{
		cfg: createDefaultConfig().(*Config),
	}
	testLogs := createTestResourceLogs()

	for n := 0; n < b.N; n++ {
		unrollProcessor.ProcessLogs(context.Background(), testLogs)
	}
}

func createTestResourceLogs() plog.Logs {
	rl := plog.NewLogs()
	for i := 0; i < 10; i++ {
		resourceLog := rl.ResourceLogs().AppendEmpty()
		for j := 0; j < 10; j++ {
			scopeLogs := resourceLog.ScopeLogs().AppendEmpty()
			scopeLogs.LogRecords().AppendEmpty().Body().SetEmptySlice().FromRaw([]any{1, 2, 3, 4, 5, 6, 7})
		}
	}
	return rl
}
