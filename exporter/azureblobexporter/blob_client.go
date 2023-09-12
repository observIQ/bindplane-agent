package azureblobexporter // import "github.com/observiq/bindplane-agent/exporter/azureblobexporter"

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

// blobClient is a wrapper for an Azure Blob client to allow mocking for testing.
//
//go:generate mockery --name blobClient --output ./internal/mocks --with-expecter --filename mock_blob_client.go --structname mockBlobClient
type blobClient interface {
	// UploadBuffer uploads a buffer in blocks to a block blob.
	UploadBuffer(context.Context, string, string, []byte) error
}

// azureBlobClient is the azure implementation of the blobClient
type azureBlobClient struct {
	azClient *azblob.Client
}

// newAzureBlobClient creates a new azureBlobClient with the given connection string
func newAzureBlobClient(connectionString string) (*azureBlobClient, error) {
	azClient, err := azblob.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, err
	}

	return &azureBlobClient{
		azClient: azClient,
	}, nil
}

func (a *azureBlobClient) UploadBuffer(ctx context.Context, containerName, blobName string, buffer []byte) error {
	_, err := a.azClient.UploadBuffer(ctx, containerName, blobName, buffer, nil)
	return err
}
