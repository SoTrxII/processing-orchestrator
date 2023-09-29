// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.24.2
// source: proto/processor.proto

package processing_orchestrator

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type WatchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *WatchRequest) Reset() {
	*x = WatchRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_processor_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WatchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WatchRequest) ProtoMessage() {}

func (x *WatchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_processor_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WatchRequest.ProtoReflect.Descriptor instead.
func (*WatchRequest) Descriptor() ([]byte, []int) {
	return file_proto_processor_proto_rawDescGZIP(), []int{0}
}

func (x *WatchRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type ProcessingStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id               string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Error            string   `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	StepsList        []string `protobuf:"bytes,3,rep,name=stepsList,proto3" json:"stepsList,omitempty"`
	CurrentStepIndex uint32   `protobuf:"varint,4,opt,name=currentStepIndex,proto3" json:"currentStepIndex,omitempty"`
	ItemList         []string `protobuf:"bytes,5,rep,name=itemList,proto3" json:"itemList,omitempty"`
	CurrentItemIndex uint32   `protobuf:"varint,6,opt,name=currentItemIndex,proto3" json:"currentItemIndex,omitempty"`
	// All processing is done
	Done bool `protobuf:"varint,7,opt,name=done,proto3" json:"done,omitempty"`
	// Progress as a user friendly string
	Progress string `protobuf:"bytes,8,opt,name=progress,proto3" json:"progress,omitempty"`
}

func (x *ProcessingStatus) Reset() {
	*x = ProcessingStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_processor_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessingStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessingStatus) ProtoMessage() {}

func (x *ProcessingStatus) ProtoReflect() protoreflect.Message {
	mi := &file_proto_processor_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessingStatus.ProtoReflect.Descriptor instead.
func (*ProcessingStatus) Descriptor() ([]byte, []int) {
	return file_proto_processor_proto_rawDescGZIP(), []int{1}
}

func (x *ProcessingStatus) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ProcessingStatus) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

func (x *ProcessingStatus) GetStepsList() []string {
	if x != nil {
		return x.StepsList
	}
	return nil
}

func (x *ProcessingStatus) GetCurrentStepIndex() uint32 {
	if x != nil {
		return x.CurrentStepIndex
	}
	return 0
}

func (x *ProcessingStatus) GetItemList() []string {
	if x != nil {
		return x.ItemList
	}
	return nil
}

func (x *ProcessingStatus) GetCurrentItemIndex() uint32 {
	if x != nil {
		return x.CurrentItemIndex
	}
	return 0
}

func (x *ProcessingStatus) GetDone() bool {
	if x != nil {
		return x.Done
	}
	return false
}

func (x *ProcessingStatus) GetProgress() string {
	if x != nil {
		return x.Progress
	}
	return ""
}

type ProcessRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DiscordAudioKeys   []string `protobuf:"bytes,1,rep,name=discordAudioKeys,proto3" json:"discordAudioKeys,omitempty"`
	BackgroundAudioKey string   `protobuf:"bytes,2,opt,name=backgroundAudioKey,proto3" json:"backgroundAudioKey,omitempty"`
}

func (x *ProcessRequest) Reset() {
	*x = ProcessRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_processor_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessRequest) ProtoMessage() {}

func (x *ProcessRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_processor_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessRequest.ProtoReflect.Descriptor instead.
func (*ProcessRequest) Descriptor() ([]byte, []int) {
	return file_proto_processor_proto_rawDescGZIP(), []int{2}
}

func (x *ProcessRequest) GetDiscordAudioKeys() []string {
	if x != nil {
		return x.DiscordAudioKeys
	}
	return nil
}

func (x *ProcessRequest) GetBackgroundAudioKey() string {
	if x != nil {
		return x.BackgroundAudioKey
	}
	return ""
}

type ProcessResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *ProcessResponse) Reset() {
	*x = ProcessResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_processor_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessResponse) ProtoMessage() {}

func (x *ProcessResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_processor_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessResponse.ProtoReflect.Descriptor instead.
func (*ProcessResponse) Descriptor() ([]byte, []int) {
	return file_proto_processor_proto_rawDescGZIP(), []int{3}
}

func (x *ProcessResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

var File_proto_processor_proto protoreflect.FileDescriptor

var file_proto_processor_proto_rawDesc = []byte{
	0x0a, 0x15, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x6f,
	0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x17, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73,
	0x69, 0x6e, 0x67, 0x5f, 0x6f, 0x72, 0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72,
	0x22, 0x1e, 0x0a, 0x0c, 0x57, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x22, 0xfa, 0x01, 0x0a, 0x10, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x69, 0x6e, 0x67, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x1c, 0x0a, 0x09, 0x73,
	0x74, 0x65, 0x70, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09,
	0x73, 0x74, 0x65, 0x70, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x2a, 0x0a, 0x10, 0x63, 0x75, 0x72,
	0x72, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x65, 0x70, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x10, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x65, 0x70,
	0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x74, 0x65, 0x6d, 0x4c, 0x69, 0x73,
	0x74, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x69, 0x74, 0x65, 0x6d, 0x4c, 0x69, 0x73,
	0x74, 0x12, 0x2a, 0x0a, 0x10, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x49, 0x74, 0x65, 0x6d,
	0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x10, 0x63, 0x75, 0x72,
	0x72, 0x65, 0x6e, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x12, 0x0a,
	0x04, 0x64, 0x6f, 0x6e, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x64, 0x6f, 0x6e,
	0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73, 0x73, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73, 0x73, 0x22, 0x6c, 0x0a,
	0x0e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x2a, 0x0a, 0x10, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x41, 0x75, 0x64, 0x69, 0x6f, 0x4b,
	0x65, 0x79, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x10, 0x64, 0x69, 0x73, 0x63, 0x6f,
	0x72, 0x64, 0x41, 0x75, 0x64, 0x69, 0x6f, 0x4b, 0x65, 0x79, 0x73, 0x12, 0x2e, 0x0a, 0x12, 0x62,
	0x61, 0x63, 0x6b, 0x67, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x41, 0x75, 0x64, 0x69, 0x6f, 0x4b, 0x65,
	0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x62, 0x61, 0x63, 0x6b, 0x67, 0x72, 0x6f,
	0x75, 0x6e, 0x64, 0x41, 0x75, 0x64, 0x69, 0x6f, 0x4b, 0x65, 0x79, 0x22, 0x21, 0x0a, 0x0f, 0x50,
	0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x32, 0xc4,
	0x01, 0x0a, 0x09, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x12, 0x5a, 0x0a, 0x05,
	0x53, 0x74, 0x61, 0x72, 0x74, 0x12, 0x27, 0x2e, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x69,
	0x6e, 0x67, 0x5f, 0x6f, 0x72, 0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e,
	0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x28,
	0x2e, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x69, 0x6e, 0x67, 0x5f, 0x6f, 0x72, 0x63, 0x68,
	0x65, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x5b, 0x0a, 0x05, 0x57, 0x61, 0x74, 0x63,
	0x68, 0x12, 0x25, 0x2e, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x69, 0x6e, 0x67, 0x5f, 0x6f,
	0x72, 0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x57, 0x61, 0x74, 0x63,
	0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x29, 0x2e, 0x70, 0x72, 0x6f, 0x63, 0x65,
	0x73, 0x73, 0x69, 0x6e, 0x67, 0x5f, 0x6f, 0x72, 0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61, 0x74,
	0x6f, 0x72, 0x2e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x69, 0x6e, 0x67, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x30, 0x01, 0x42, 0x1b, 0x5a, 0x19, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x63, 0x65,
	0x73, 0x73, 0x69, 0x6e, 0x67, 0x2d, 0x6f, 0x72, 0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61, 0x74,
	0x6f, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_processor_proto_rawDescOnce sync.Once
	file_proto_processor_proto_rawDescData = file_proto_processor_proto_rawDesc
)

func file_proto_processor_proto_rawDescGZIP() []byte {
	file_proto_processor_proto_rawDescOnce.Do(func() {
		file_proto_processor_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_processor_proto_rawDescData)
	})
	return file_proto_processor_proto_rawDescData
}

var file_proto_processor_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_processor_proto_goTypes = []interface{}{
	(*WatchRequest)(nil),     // 0: processing_orchestrator.WatchRequest
	(*ProcessingStatus)(nil), // 1: processing_orchestrator.ProcessingStatus
	(*ProcessRequest)(nil),   // 2: processing_orchestrator.ProcessRequest
	(*ProcessResponse)(nil),  // 3: processing_orchestrator.ProcessResponse
}
var file_proto_processor_proto_depIdxs = []int32{
	2, // 0: processing_orchestrator.Processor.Start:input_type -> processing_orchestrator.ProcessRequest
	0, // 1: processing_orchestrator.Processor.Watch:input_type -> processing_orchestrator.WatchRequest
	3, // 2: processing_orchestrator.Processor.Start:output_type -> processing_orchestrator.ProcessResponse
	1, // 3: processing_orchestrator.Processor.Watch:output_type -> processing_orchestrator.ProcessingStatus
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_processor_proto_init() }
func file_proto_processor_proto_init() {
	if File_proto_processor_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_processor_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WatchRequest); i {
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
		file_proto_processor_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProcessingStatus); i {
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
		file_proto_processor_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProcessRequest); i {
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
		file_proto_processor_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProcessResponse); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_processor_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_processor_proto_goTypes,
		DependencyIndexes: file_proto_processor_proto_depIdxs,
		MessageInfos:      file_proto_processor_proto_msgTypes,
	}.Build()
	File_proto_processor_proto = out.File
	file_proto_processor_proto_rawDesc = nil
	file_proto_processor_proto_goTypes = nil
	file_proto_processor_proto_depIdxs = nil
}
