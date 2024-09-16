// Request message for importing logs.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.3
// source: chronicle_http.proto

package api

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Represents a telemetry log.
type Log struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The resource name of this log.
	// Format:
	// projects/{project}/locations/{region}/instances/{instance}/logTypes/{log_type}/logs
	// /{log}
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Raw data for the log entry.
	Data []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	// Timestamp of the log entry.
	LogEntryTime *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=log_entry_time,json=logEntryTime,proto3" json:"log_entry_time,omitempty"`
	// The time at which the log entry was collected. Must be after the
	// log_entry_time.
	CollectionTime *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=collection_time,json=collectionTime,proto3" json:"collection_time,omitempty"`
	// The user-configured environment namespace to identify the data
	// domain the logs originated from. This namespace will be used as a tag to
	// identify the appropriate data domain for indexing and enrichment
	// functionality.
	EnvironmentNamespace string `protobuf:"bytes,5,opt,name=environment_namespace,json=environmentNamespace,proto3" json:"environment_namespace,omitempty"`
	// The user-configured custom metadata labels.
	Labels map[string]*Log_LogLabel `protobuf:"bytes,6,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Log) Reset() {
	*x = Log{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chronicle_http_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Log) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Log) ProtoMessage() {}

func (x *Log) ProtoReflect() protoreflect.Message {
	mi := &file_chronicle_http_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Log.ProtoReflect.Descriptor instead.
func (*Log) Descriptor() ([]byte, []int) {
	return file_chronicle_http_proto_rawDescGZIP(), []int{0}
}

func (x *Log) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Log) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Log) GetLogEntryTime() *timestamppb.Timestamp {
	if x != nil {
		return x.LogEntryTime
	}
	return nil
}

func (x *Log) GetCollectionTime() *timestamppb.Timestamp {
	if x != nil {
		return x.CollectionTime
	}
	return nil
}

func (x *Log) GetEnvironmentNamespace() string {
	if x != nil {
		return x.EnvironmentNamespace
	}
	return ""
}

func (x *Log) GetLabels() map[string]*Log_LogLabel {
	if x != nil {
		return x.Labels
	}
	return nil
}

type ImportLogsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The parent, which owns this collection of logs.
	Parent string `protobuf:"bytes,1,opt,name=parent,proto3" json:"parent,omitempty"`
	// Types that are assignable to Source:
	//
	//	*ImportLogsRequest_InlineSource
	Source isImportLogsRequest_Source `protobuf_oneof:"source"`
	// Opaque hint to help parsing the log.
	Hint string `protobuf:"bytes,4,opt,name=hint,proto3" json:"hint,omitempty"`
}

func (x *ImportLogsRequest) Reset() {
	*x = ImportLogsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chronicle_http_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ImportLogsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ImportLogsRequest) ProtoMessage() {}

func (x *ImportLogsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_chronicle_http_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ImportLogsRequest.ProtoReflect.Descriptor instead.
func (*ImportLogsRequest) Descriptor() ([]byte, []int) {
	return file_chronicle_http_proto_rawDescGZIP(), []int{1}
}

func (x *ImportLogsRequest) GetParent() string {
	if x != nil {
		return x.Parent
	}
	return ""
}

func (m *ImportLogsRequest) GetSource() isImportLogsRequest_Source {
	if m != nil {
		return m.Source
	}
	return nil
}

func (x *ImportLogsRequest) GetInlineSource() *ImportLogsRequest_LogsInlineSource {
	if x, ok := x.GetSource().(*ImportLogsRequest_InlineSource); ok {
		return x.InlineSource
	}
	return nil
}

func (x *ImportLogsRequest) GetHint() string {
	if x != nil {
		return x.Hint
	}
	return ""
}

type isImportLogsRequest_Source interface {
	isImportLogsRequest_Source()
}

type ImportLogsRequest_InlineSource struct {
	// Logs to be imported are specified inline.
	InlineSource *ImportLogsRequest_LogsInlineSource `protobuf:"bytes,2,opt,name=inline_source,json=inlineSource,proto3,oneof"`
}

func (*ImportLogsRequest_InlineSource) isImportLogsRequest_Source() {}

// Label for a user configured custom metadata key.
type Log_LogLabel struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The value of the label.
	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	// Indicates whether this label can be used for Data RBAC.
	RbacEnabled bool `protobuf:"varint,2,opt,name=rbac_enabled,json=rbacEnabled,proto3" json:"rbac_enabled,omitempty"`
}

func (x *Log_LogLabel) Reset() {
	*x = Log_LogLabel{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chronicle_http_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Log_LogLabel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Log_LogLabel) ProtoMessage() {}

func (x *Log_LogLabel) ProtoReflect() protoreflect.Message {
	mi := &file_chronicle_http_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Log_LogLabel.ProtoReflect.Descriptor instead.
func (*Log_LogLabel) Descriptor() ([]byte, []int) {
	return file_chronicle_http_proto_rawDescGZIP(), []int{0, 1}
}

func (x *Log_LogLabel) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *Log_LogLabel) GetRbacEnabled() bool {
	if x != nil {
		return x.RbacEnabled
	}
	return false
}

// A import source with the logs to import included inline.
type ImportLogsRequest_LogsInlineSource struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The logs being imported.
	Logs []*Log `protobuf:"bytes,1,rep,name=logs,proto3" json:"logs,omitempty"`
	// The forwarder sending this import request.
	Forwarder string `protobuf:"bytes,2,opt,name=forwarder,proto3" json:"forwarder,omitempty"`
	// Source file name. Populated for certain types of files processed by the
	// outofband processor which may have metadata encoded in it for use by
	// the parser.
	SourceFilename string `protobuf:"bytes,3,opt,name=source_filename,json=sourceFilename,proto3" json:"source_filename,omitempty"`
}

func (x *ImportLogsRequest_LogsInlineSource) Reset() {
	*x = ImportLogsRequest_LogsInlineSource{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chronicle_http_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ImportLogsRequest_LogsInlineSource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ImportLogsRequest_LogsInlineSource) ProtoMessage() {}

func (x *ImportLogsRequest_LogsInlineSource) ProtoReflect() protoreflect.Message {
	mi := &file_chronicle_http_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ImportLogsRequest_LogsInlineSource.ProtoReflect.Descriptor instead.
func (*ImportLogsRequest_LogsInlineSource) Descriptor() ([]byte, []int) {
	return file_chronicle_http_proto_rawDescGZIP(), []int{1, 0}
}

func (x *ImportLogsRequest_LogsInlineSource) GetLogs() []*Log {
	if x != nil {
		return x.Logs
	}
	return nil
}

func (x *ImportLogsRequest_LogsInlineSource) GetForwarder() string {
	if x != nil {
		return x.Forwarder
	}
	return ""
}

func (x *ImportLogsRequest_LogsInlineSource) GetSourceFilename() string {
	if x != nil {
		return x.SourceFilename
	}
	return ""
}

var File_chronicle_http_proto protoreflect.FileDescriptor

var file_chronicle_http_proto_rawDesc = []byte{
	0x0a, 0x14, 0x63, 0x68, 0x72, 0x6f, 0x6e, 0x69, 0x63, 0x6c, 0x65, 0x5f, 0x68, 0x74, 0x74, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa2, 0x03, 0x0a, 0x03, 0x4c, 0x6f, 0x67, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x40, 0x0a, 0x0e, 0x6c, 0x6f, 0x67, 0x5f, 0x65,
	0x6e, 0x74, 0x72, 0x79, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0c, 0x6c, 0x6f, 0x67,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x43, 0x0a, 0x0f, 0x63, 0x6f, 0x6c,
	0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0e,
	0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x33,
	0x0a, 0x15, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x14, 0x65,
	0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70,
	0x61, 0x63, 0x65, 0x12, 0x28, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x06, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x4c, 0x6f, 0x67, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x1a, 0x48, 0x0a,
	0x0b, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x23,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e,
	0x4c, 0x6f, 0x67, 0x2e, 0x4c, 0x6f, 0x67, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x43, 0x0a, 0x08, 0x4c, 0x6f, 0x67, 0x4c, 0x61,
	0x62, 0x65, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x72, 0x62, 0x61,
	0x63, 0x5f, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x0b, 0x72, 0x62, 0x61, 0x63, 0x45, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x22, 0x8a, 0x02, 0x0a,
	0x11, 0x49, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x4c, 0x6f, 0x67, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x12, 0x4a, 0x0a, 0x0d, 0x69, 0x6e,
	0x6c, 0x69, 0x6e, 0x65, 0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x23, 0x2e, 0x49, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x4c, 0x6f, 0x67, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x4c, 0x6f, 0x67, 0x73, 0x49, 0x6e, 0x6c, 0x69, 0x6e, 0x65,
	0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x48, 0x00, 0x52, 0x0c, 0x69, 0x6e, 0x6c, 0x69, 0x6e, 0x65,
	0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x69, 0x6e, 0x74, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x69, 0x6e, 0x74, 0x1a, 0x73, 0x0a, 0x10, 0x4c, 0x6f,
	0x67, 0x73, 0x49, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x18,
	0x0a, 0x04, 0x6c, 0x6f, 0x67, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x04, 0x2e, 0x4c,
	0x6f, 0x67, 0x52, 0x04, 0x6c, 0x6f, 0x67, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x66, 0x6f, 0x72, 0x77,
	0x61, 0x72, 0x64, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x66, 0x6f, 0x72,
	0x77, 0x61, 0x72, 0x64, 0x65, 0x72, 0x12, 0x27, 0x0a, 0x0f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x5f, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0e, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x42,
	0x08, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x42, 0x31, 0x5a, 0x2f, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72,
	0x2f, 0x63, 0x68, 0x72, 0x6f, 0x6e, 0x69, 0x63, 0x6c, 0x65, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74,
	0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_chronicle_http_proto_rawDescOnce sync.Once
	file_chronicle_http_proto_rawDescData = file_chronicle_http_proto_rawDesc
)

func file_chronicle_http_proto_rawDescGZIP() []byte {
	file_chronicle_http_proto_rawDescOnce.Do(func() {
		file_chronicle_http_proto_rawDescData = protoimpl.X.CompressGZIP(file_chronicle_http_proto_rawDescData)
	})
	return file_chronicle_http_proto_rawDescData
}

var file_chronicle_http_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_chronicle_http_proto_goTypes = []any{
	(*Log)(nil),               // 0: Log
	(*ImportLogsRequest)(nil), // 1: ImportLogsRequest
	nil,                       // 2: Log.LabelsEntry
	(*Log_LogLabel)(nil),      // 3: Log.LogLabel
	(*ImportLogsRequest_LogsInlineSource)(nil), // 4: ImportLogsRequest.LogsInlineSource
	(*timestamppb.Timestamp)(nil),              // 5: google.protobuf.Timestamp
}
var file_chronicle_http_proto_depIdxs = []int32{
	5, // 0: Log.log_entry_time:type_name -> google.protobuf.Timestamp
	5, // 1: Log.collection_time:type_name -> google.protobuf.Timestamp
	2, // 2: Log.labels:type_name -> Log.LabelsEntry
	4, // 3: ImportLogsRequest.inline_source:type_name -> ImportLogsRequest.LogsInlineSource
	3, // 4: Log.LabelsEntry.value:type_name -> Log.LogLabel
	0, // 5: ImportLogsRequest.LogsInlineSource.logs:type_name -> Log
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_chronicle_http_proto_init() }
func file_chronicle_http_proto_init() {
	if File_chronicle_http_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_chronicle_http_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Log); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_chronicle_http_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*ImportLogsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_chronicle_http_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*Log_LogLabel); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_chronicle_http_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*ImportLogsRequest_LogsInlineSource); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_chronicle_http_proto_msgTypes[1].OneofWrappers = []any{
		(*ImportLogsRequest_InlineSource)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_chronicle_http_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_chronicle_http_proto_goTypes,
		DependencyIndexes: file_chronicle_http_proto_depIdxs,
		MessageInfos:      file_chronicle_http_proto_msgTypes,
	}.Build()
	File_chronicle_http_proto = out.File
	file_chronicle_http_proto_rawDesc = nil
	file_chronicle_http_proto_goTypes = nil
	file_chronicle_http_proto_depIdxs = nil
}