package collector

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/config"
)

// FileProvider is a config provider that uses a file.
type FileProvider struct {
	filePath string
}

// NewFileProvider returns a new file provider.
func NewFileProvider(filePath string) *FileProvider {
	return &FileProvider{
		filePath: filePath,
	}
}

// Get returns a config map from the provider.
func (f *FileProvider) Get(ctx context.Context) (*config.Map, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	configMap, err := config.NewMapFromFile(f.filePath)
	if err != nil {
		return nil, fmt.Errorf("error loading config file: %w", err)
	}

	return configMap, nil
}

// Close closes the file provider.
func (f *FileProvider) Close(ctx context.Context) error {
	return nil
}
