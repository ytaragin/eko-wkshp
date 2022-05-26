// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.6.1
// source: tasks/tasks.proto

package tasks

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

type TaskMessage_TaskStatus int32

const (
	TaskMessage_CREATED    TaskMessage_TaskStatus = 0
	TaskMessage_INPROGRESS TaskMessage_TaskStatus = 1
	TaskMessage_COMPLETE   TaskMessage_TaskStatus = 2
	TaskMessage_ERROR      TaskMessage_TaskStatus = 10
	TaskMessage_UNKNOWN    TaskMessage_TaskStatus = 99
)

// Enum value maps for TaskMessage_TaskStatus.
var (
	TaskMessage_TaskStatus_name = map[int32]string{
		0:  "CREATED",
		1:  "INPROGRESS",
		2:  "COMPLETE",
		10: "ERROR",
		99: "UNKNOWN",
	}
	TaskMessage_TaskStatus_value = map[string]int32{
		"CREATED":    0,
		"INPROGRESS": 1,
		"COMPLETE":   2,
		"ERROR":      10,
		"UNKNOWN":    99,
	}
)

func (x TaskMessage_TaskStatus) Enum() *TaskMessage_TaskStatus {
	p := new(TaskMessage_TaskStatus)
	*p = x
	return p
}

func (x TaskMessage_TaskStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TaskMessage_TaskStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_tasks_tasks_proto_enumTypes[0].Descriptor()
}

func (TaskMessage_TaskStatus) Type() protoreflect.EnumType {
	return &file_tasks_tasks_proto_enumTypes[0]
}

func (x TaskMessage_TaskStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TaskMessage_TaskStatus.Descriptor instead.
func (TaskMessage_TaskStatus) EnumDescriptor() ([]byte, []int) {
	return file_tasks_tasks_proto_rawDescGZIP(), []int{1, 0}
}

type CreateTaskRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CreateTaskRequest) Reset() {
	*x = CreateTaskRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tasks_tasks_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateTaskRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateTaskRequest) ProtoMessage() {}

func (x *CreateTaskRequest) ProtoReflect() protoreflect.Message {
	mi := &file_tasks_tasks_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateTaskRequest.ProtoReflect.Descriptor instead.
func (*CreateTaskRequest) Descriptor() ([]byte, []int) {
	return file_tasks_tasks_proto_rawDescGZIP(), []int{0}
}

// The response message containing the greetings
type TaskMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Taskid string                 `protobuf:"bytes,1,opt,name=taskid,proto3" json:"taskid,omitempty"`
	Status TaskMessage_TaskStatus `protobuf:"varint,2,opt,name=status,proto3,enum=tasks.TaskMessage_TaskStatus" json:"status,omitempty"`
}

func (x *TaskMessage) Reset() {
	*x = TaskMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tasks_tasks_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskMessage) ProtoMessage() {}

func (x *TaskMessage) ProtoReflect() protoreflect.Message {
	mi := &file_tasks_tasks_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskMessage.ProtoReflect.Descriptor instead.
func (*TaskMessage) Descriptor() ([]byte, []int) {
	return file_tasks_tasks_proto_rawDescGZIP(), []int{1}
}

func (x *TaskMessage) GetTaskid() string {
	if x != nil {
		return x.Taskid
	}
	return ""
}

func (x *TaskMessage) GetStatus() TaskMessage_TaskStatus {
	if x != nil {
		return x.Status
	}
	return TaskMessage_CREATED
}

var File_tasks_tasks_proto protoreflect.FileDescriptor

var file_tasks_tasks_proto_rawDesc = []byte{
	0x0a, 0x11, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x2f, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x22, 0x13, 0x0a, 0x11, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22,
	0xad, 0x01, 0x0a, 0x0b, 0x54, 0x61, 0x73, 0x6b, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x16, 0x0a, 0x06, 0x74, 0x61, 0x73, 0x6b, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x74, 0x61, 0x73, 0x6b, 0x69, 0x64, 0x12, 0x35, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1d, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x2e,
	0x54, 0x61, 0x73, 0x6b, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x54, 0x61, 0x73, 0x6b,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x4f,
	0x0a, 0x0a, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0b, 0x0a, 0x07,
	0x43, 0x52, 0x45, 0x41, 0x54, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x49, 0x4e, 0x50,
	0x52, 0x4f, 0x47, 0x52, 0x45, 0x53, 0x53, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x43, 0x4f, 0x4d,
	0x50, 0x4c, 0x45, 0x54, 0x45, 0x10, 0x02, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x52, 0x52, 0x4f, 0x52,
	0x10, 0x0a, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x63, 0x32,
	0x7d, 0x0a, 0x05, 0x54, 0x61, 0x73, 0x6b, 0x73, 0x12, 0x3c, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x18, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x2e, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x12, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x22, 0x00, 0x12, 0x36, 0x0a, 0x0a, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x54, 0x61, 0x73, 0x6b, 0x12, 0x12, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x2e, 0x54, 0x61, 0x73,
	0x6b, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x12, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x73,
	0x2e, 0x54, 0x61, 0x73, 0x6b, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x00, 0x42, 0x0a,
	0x5a, 0x08, 0x2e, 0x2e, 0x2f, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_tasks_tasks_proto_rawDescOnce sync.Once
	file_tasks_tasks_proto_rawDescData = file_tasks_tasks_proto_rawDesc
)

func file_tasks_tasks_proto_rawDescGZIP() []byte {
	file_tasks_tasks_proto_rawDescOnce.Do(func() {
		file_tasks_tasks_proto_rawDescData = protoimpl.X.CompressGZIP(file_tasks_tasks_proto_rawDescData)
	})
	return file_tasks_tasks_proto_rawDescData
}

var file_tasks_tasks_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_tasks_tasks_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_tasks_tasks_proto_goTypes = []interface{}{
	(TaskMessage_TaskStatus)(0), // 0: tasks.TaskMessage.TaskStatus
	(*CreateTaskRequest)(nil),   // 1: tasks.CreateTaskRequest
	(*TaskMessage)(nil),         // 2: tasks.TaskMessage
}
var file_tasks_tasks_proto_depIdxs = []int32{
	0, // 0: tasks.TaskMessage.status:type_name -> tasks.TaskMessage.TaskStatus
	1, // 1: tasks.Tasks.CreateTask:input_type -> tasks.CreateTaskRequest
	2, // 2: tasks.Tasks.UpdateTask:input_type -> tasks.TaskMessage
	2, // 3: tasks.Tasks.CreateTask:output_type -> tasks.TaskMessage
	2, // 4: tasks.Tasks.UpdateTask:output_type -> tasks.TaskMessage
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_tasks_tasks_proto_init() }
func file_tasks_tasks_proto_init() {
	if File_tasks_tasks_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_tasks_tasks_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateTaskRequest); i {
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
		file_tasks_tasks_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TaskMessage); i {
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
			RawDescriptor: file_tasks_tasks_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_tasks_tasks_proto_goTypes,
		DependencyIndexes: file_tasks_tasks_proto_depIdxs,
		EnumInfos:         file_tasks_tasks_proto_enumTypes,
		MessageInfos:      file_tasks_tasks_proto_msgTypes,
	}.Build()
	File_tasks_tasks_proto = out.File
	file_tasks_tasks_proto_rawDesc = nil
	file_tasks_tasks_proto_goTypes = nil
	file_tasks_tasks_proto_depIdxs = nil
}