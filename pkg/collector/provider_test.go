package collector

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFileProvider(t *testing.T) {
	provider := NewFileProvider("./test/valid.yaml")
	require.Equal(t, "./test/valid.yaml", provider.retriever.filePath)
}

func TestFileProviderGetValid(t *testing.T) {
	provider := NewFileProvider("./test/valid.yaml")
	parser, err := provider.retriever.Get(context.Background())
	require.NoError(t, err)
	require.NotNil(t, parser)
}

func TestFileProviderGetFail(t *testing.T) {
	provider := NewFileProvider("./test/not_existing.yaml")
	parser, err := provider.retriever.Get(context.Background())
	require.Error(t, err)
	require.Contains(t, err.Error(), "error loading config file")
	require.Nil(t, parser)
}
