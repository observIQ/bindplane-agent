package logging

import (
	"testing"
)

func TestCreateFileCore(t *testing.T) {
	// Test that CreateFileCore doesn't panic for default config
	config := DefaultConfig()

	_ = CreateFileCore(&config.Collector)
	_ = CreateFileCore(&config.Manager)
}
