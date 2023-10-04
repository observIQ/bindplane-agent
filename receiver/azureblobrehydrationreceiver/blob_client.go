package azureblobrehydrationreceiver //import "github.com/observiq/bindplane-agent/receiver/azureblobrehydrationreceiver"

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

// blobInfo contains the necessary info to process a blob
type blobInfo struct {
	Name string
	Size int64
}

// blobClient
type blobClient interface {
	// ListBlobs returns a list of blobInfo objects present in the container with the given prefix
	ListBlobs(ctx context.Context, container string, prefix, marker *string) ([]*blobInfo, *string, error)

	// DownloadBlob downloads the contents of the blob into the supplied buffer.
	// It will return the count of bytes used in the buffer.
	DownloadBlob(ctx context.Context, container, blobPath string, buf []byte) (int64, error)

	// DeleteBlob deletes the blob in the specified container
	DeleteBlob(ctx context.Context, container, blobPath string) error
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

// contentLengthKey key for the content length metadata
const contentLengthKey = "ContentLength"

// ListBlobs returns a list of blobInfo objects present in the container with the given prefix
func (a *azureBlobClient) ListBlobs(ctx context.Context, container string, prefix, marker *string) ([]*blobInfo, *string, error) {
	listOptions := &azblob.ListBlobsFlatOptions{
		Marker: marker,
		Prefix: prefix,
	}

	pager := a.azClient.NewListBlobsFlatPager(container, listOptions)

	var nextMarker *string
	blobs := make([]*blobInfo, 0)
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
			// All blob fields are pointers so check all pointers we need before we try to process it
			if blob.Name == nil || blob.Properties == nil || blob.Properties.ContentLength == nil {
				continue
			}

			info := &blobInfo{
				Name: *blob.Name,
				Size: *blob.Properties.ContentLength,
			}

			blobs = append(blobs, info)
		}
		nextMarker = resp.NextMarker
	}

	return blobs, nextMarker, nil
}

func (a *azureBlobClient) DownloadBlob(ctx context.Context, container, blobPath string, buf []byte) (int64, error) {
	bytesDownloaded, err := a.azClient.DownloadBuffer(ctx, container, blobPath, buf, nil)
	if err != nil {
		return 0, fmt.Errorf("download: %w", err)
	}

	return bytesDownloaded, nil
}

// DeleteBlob deletes the blob in the specified container
func (a *azureBlobClient) DeleteBlob(ctx context.Context, container, blobPath string) error {
	_, err := a.azClient.DeleteBlob(ctx, container, blobPath, nil)
	return err
}
