package azureblobexporter // import "github.com/observiq/bindplane-agent/exporter/azureblobexporter"

import (
	"errors"
	"fmt"
)

// partitionType is the type of partition to store blobs under
type partitionType string

const (
	minute partitionType = "minute"
	hour   partitionType = "hour"
)

type Config struct {
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

	// TODO compression
}

func (c *Config) Validate() error {
	if c.ConnectionString == "" {
		return errors.New("connection_string is required")
	}

	if c.Container == "" {
		return errors.New("container is required")
	}

	if c.Partition != minute && c.Partition != hour {
		return fmt.Errorf("invalid partition type '%s'", c.Partition)
	}

	return nil
}
