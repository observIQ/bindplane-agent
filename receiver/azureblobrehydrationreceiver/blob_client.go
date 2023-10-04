package azureblobrehydrationreceiver //import "github.com/observiq/bindplane-agent/receiver/azureblobrehydrationreceiver"

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

// blobClient
type blobClient interface {
	// ListBlobs returns a list of blob names present in the container with the given prefix
	ListBlobs(ctx context.Context, container string, prefix, marker *string) ([]string, *string, error)

	// DownloadBlob downloads the contents of the blob into the supplied buffer.
	// It will return the count of bytes used in the buffer.
	DownloadBlob(ctx context.Context, container, blobPath string, buf []byte) (int64, error)
}

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

// ListBlobs returns a list of blob names present in the container with the given prefix
func (a *azureBlobClient) ListBlobs(ctx context.Context, container string, prefix, marker *string) ([]string, *string, error) {
	listOptions := &azblob.ListBlobsFlatOptions{
		Marker: marker,
		Prefix: prefix,
	}

	pager := a.azClient.NewListBlobsFlatPager(container, listOptions)

	var nextMarker *string
	blobNames := make([]string, 0)
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("listBlobs: %w", err)
		}

		for _, blob := range resp.Segment.BlobItems {
			// Skip deleted blobs
			if blob.Deleted != nil && *blob.Deleted {
				continue
			}

			if blob.Name != nil {
				blobNames = append(blobNames, *blob.Name)
			}
		}
		nextMarker = resp.NextMarker
	}

	return blobNames, nextMarker, nil
}

func (a *azureBlobClient) DownloadBlob(ctx context.Context, container, blobPath string, buf []byte) (int64, error) {
	bytesDownloaded, err := a.azClient.DownloadBuffer(ctx, container, blobPath, buf, nil)
	if err != nil {
		return 0, fmt.Errorf("download: %w", err)
	}

	return bytesDownloaded, nil
}
