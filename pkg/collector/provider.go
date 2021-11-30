package collector

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/configmapprovider"
)

// FileProvider is a config provider that uses a file.
type FileProvider struct {
	retriever *FileRetriever
}

// NewFileProvider returns a new file provider.
func NewFileProvider(filePath string) *FileProvider {
	return &FileProvider{
		retriever: &FileRetriever{
			filePath: filePath,
		},
	}
}

// Retrieve returns a FileRetriever
func (f *FileProvider) Retrieve(ctx context.Context, onChange func(*configmapprovider.ChangeEvent)) (configmapprovider.Retrieved, error) {
	return f.retriever, nil
}

// Shutdown stops the file provider.
func (f *FileProvider) Shutdown(ctx context.Context) error {
	return nil
}

// FileRetriever is a retriever that retrieves a configuration from a file
type FileRetriever struct {
	filePath string
}

// Get returns a config map from the file.
func (f *FileRetriever) Get(ctx context.Context) (*config.Map, error) {
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

// Close closes the file retriever.
func (f *FileRetriever) Close(ctx context.Context) error {
	return nil
}
