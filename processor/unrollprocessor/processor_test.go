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
