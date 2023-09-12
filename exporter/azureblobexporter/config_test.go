package azureblobexporter // import "github.com/observiq/bindplane-agent/exporter/azureblobexporter"

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigValidate(t *testing.T) {
	testCases := []struct {
		desc        string
		config      *Config
		expectedErr error
	}{
		{
			desc: "Empty connection_string",
			config: &Config{
				ConnectionString: "",
				Container:        "test",
				Partition:        minutePartition,
				Compression:      noCompression,
			},
			expectedErr: errors.New("connection_string is required"),
		},
		{
			desc: "Empty container",
			config: &Config{
				ConnectionString: "connection",
				Container:        "",
				Partition:        minutePartition,
				Compression:      noCompression,
			},
			expectedErr: errors.New("container is required"),
		},
		{
			desc: "Invalid partition",
			config: &Config{
				ConnectionString: "connection",
				Container:        "test",
				Partition:        partitionType("nope"),
				Compression:      noCompression,
			},
			expectedErr: errors.New("unsupported partition type"),
		},
		{
			desc: "Invalid compression",
			config: &Config{
				ConnectionString: "connection",
				Container:        "test",
				Partition:        minutePartition,
				Compression:      compressionType("bad"),
			},
			expectedErr: errors.New("unsupported compression type"),
		},
		{
			desc: "Valid partition type hour",
			config: &Config{
				ConnectionString: "connection",
				Container:        "test",
				Partition:        hourPartition,
				Compression:      noCompression,
			},
			expectedErr: nil,
		},
		{
			desc: "Valid compression type gzip",
			config: &Config{
				ConnectionString: "connection",
				Container:        "test",
				Partition:        hourPartition,
				Compression:      gzipCompression,
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		currentTC := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			err := currentTC.config.Validate()
			if currentTC.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, currentTC.expectedErr.Error())
			}
		})
	}
}
