package azureblobexporter // import "github.com/observiq/bindplane-agent/exporter/azureblobexporter"

import (
	"errors"
	"fmt"

	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

// partitionType is the type of partition to store blobs under
type partitionType string

const (
	minutePartition partitionType = "minute"
	hourPartition   partitionType = "hour"
)

type compressionType string

const (
	noCompression   compressionType = ""
	gzipCompression compressionType = "gzip"
)

type Config struct {
	exporterhelper.TimeoutSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct.
	exporterhelper.QueueSettings   `mapstructure:"sending_queue"`
	exporterhelper.RetrySettings   `mapstructure:"retry_on_failure"`

	// Azure Blob Storage connection key,
	// which can be found in the Azure Blob Storage resource on the Azure Portal. (no default)
	ConnectionString string `mapstructure:"connection_string"`

	// Container is the name of the user created storage container. (no default)
	Container string `mapstructure:"container"`

	// BlobPrefix is the blob prefix defined by the user. (no default)
	BlobPrefix string `mapstructure:"blob_prefix"`

	// RootFolder is the name of the root folder in path.
	// Defaults to telemetry type eg: metrics, logs, traces.
	RootFolder string `mapstructure:"root_folder"`

	// Partition is the time granularity of the blob.
	// Valid values are "hour" or "minute". Default: minute
	Partition partitionType `mapstructure:"partition"`

	// Compression is the type of compression to use
	Compression compressionType `mapstructure:"compression"`
}

func (c *Config) Validate() error {
	if c.ConnectionString == "" {
		return errors.New("connection_string is required")
	}

	if c.Container == "" {
		return errors.New("container is required")
	}

	if c.Partition != minutePartition && c.Partition != hourPartition {
		return fmt.Errorf("unsupported partition type '%s'", c.Partition)
	}

	switch c.Compression {
	case noCompression, gzipCompression:
		return nil
	default:
		return fmt.Errorf("unsupported compression type: %s", c.Compression)
	}
}
