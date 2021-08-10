package collector

import (
	"fmt"

	"go.opentelemetry.io/collector/config/configparser"
)

// FileProvider is a parser provider that uses a file.
type FileProvider struct {
	filePath string
}

// NewFileProvider returns a new file provider.
func NewFileProvider(filePath string) *FileProvider {
	return &FileProvider{
		filePath: filePath,
	}
}

// Get returns a config parser from the provider.
func (f *FileProvider) Get() (*configparser.Parser, error) {
	cp, err := configparser.NewParserFromFile(f.filePath)
	if err != nil {
		return nil, fmt.Errorf("error loading config file: %w", err)
	}

	return cp, nil
}
