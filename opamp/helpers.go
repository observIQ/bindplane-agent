package opamp

import "github.com/open-telemetry/opamp-go/protobufs"

// StringKeyValue converts a string key-value pair into a protobuf.KeyValue struct
func StringKeyValue(key, value string) *protobufs.KeyValue {
	return &protobufs.KeyValue{
		Key: key,
		Value: &protobufs.AnyValue{
			Value: &protobufs.AnyValue_StringValue{StringValue: value},
		},
	}
}
