package jobpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SubmitJobRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Type          string                 `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`                                      // handler key, e.g. "email", "resize"
	Payload       []byte                 `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`                                // arbitrary JSON bytes
	Priority      int32                  `protobuf:"varint,3,opt,name=priority,proto3" json:"priority,omitempty"`                             // higher = served first (default 0)
	DelaySeconds  int64                  `protobuf:"varint,4,opt,name=delay_seconds,json=delaySeconds,proto3" json:"delay_seconds,omitempty"` // seconds from now before job becomes eligible
	MaxRetries    int32                  `protobuf:"varint,5,opt,name=max_retries,json=maxRetries,proto3" json:"max_retries,omitempty"`       // 0 = no retries (default 3)
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SubmitJobRequest) Reset() {
	*x = SubmitJobRequest{}
	mi := &file_job_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SubmitJobRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitJobRequest) ProtoMessage() {}

func (x *SubmitJobRequest) ProtoReflect() protoreflect.Message {
	mi := &file_job_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*SubmitJobRequest) Descriptor() ([]byte, []int) {
	return file_job_proto_rawDescGZIP(), []int{0}
}

func (x *SubmitJobRequest) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *SubmitJobRequest) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *SubmitJobRequest) GetPriority() int32 {
	if x != nil {
		return x.Priority
	}
	return 0
}

func (x *SubmitJobRequest) GetDelaySeconds() int64 {
	if x != nil {
		return x.DelaySeconds
	}
	return 0
}

func (x *SubmitJobRequest) GetMaxRetries() int32 {
	if x != nil {
		return x.MaxRetries
	}
	return 0
}

type SubmitJobResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	JobId         string                 `protobuf:"bytes,1,opt,name=job_id,json=jobId,proto3" json:"job_id,omitempty"`
	Status        string                 `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"` // always "queued"
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SubmitJobResponse) Reset() {
	*x = SubmitJobResponse{}
	mi := &file_job_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SubmitJobResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitJobResponse) ProtoMessage() {}

func (x *SubmitJobResponse) ProtoReflect() protoreflect.Message {
	mi := &file_job_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*SubmitJobResponse) Descriptor() ([]byte, []int) {
	return file_job_proto_rawDescGZIP(), []int{1}
}

func (x *SubmitJobResponse) GetJobId() string {
	if x != nil {
		return x.JobId
	}
	return ""
}

func (x *SubmitJobResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type GetJobRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	JobId         string                 `protobuf:"bytes,1,opt,name=job_id,json=jobId,proto3" json:"job_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetJobRequest) Reset() {
	*x = GetJobRequest{}
	mi := &file_job_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetJobRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetJobRequest) ProtoMessage() {}

func (x *GetJobRequest) ProtoReflect() protoreflect.Message {
	mi := &file_job_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*GetJobRequest) Descriptor() ([]byte, []int) {
	return file_job_proto_rawDescGZIP(), []int{2}
}

func (x *GetJobRequest) GetJobId() string {
	if x != nil {
		return x.JobId
	}
	return ""
}

type GetJobResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	JobId         string                 `protobuf:"bytes,1,opt,name=job_id,json=jobId,proto3" json:"job_id,omitempty"`
	Type          string                 `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	Status        string                 `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
	RetryCount    int32                  `protobuf:"varint,4,opt,name=retry_count,json=retryCount,proto3" json:"retry_count,omitempty"`
	Error         string                 `protobuf:"bytes,5,opt,name=error,proto3" json:"error,omitempty"`
	CreatedAt     int64                  `protobuf:"varint,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"` // Unix seconds
	UpdatedAt     int64                  `protobuf:"varint,7,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetJobResponse) Reset() {
	*x = GetJobResponse{}
	mi := &file_job_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetJobResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetJobResponse) ProtoMessage() {}

func (x *GetJobResponse) ProtoReflect() protoreflect.Message {
	mi := &file_job_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*GetJobResponse) Descriptor() ([]byte, []int) {
	return file_job_proto_rawDescGZIP(), []int{3}
}

func (x *GetJobResponse) GetJobId() string {
	if x != nil {
		return x.JobId
	}
	return ""
}

func (x *GetJobResponse) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *GetJobResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *GetJobResponse) GetRetryCount() int32 {
	if x != nil {
		return x.RetryCount
	}
	return 0
}

func (x *GetJobResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

func (x *GetJobResponse) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *GetJobResponse) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

type ListDeadJobsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListDeadJobsRequest) Reset() {
	*x = ListDeadJobsRequest{}
	mi := &file_job_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListDeadJobsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListDeadJobsRequest) ProtoMessage() {}

func (x *ListDeadJobsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_job_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ListDeadJobsRequest) Descriptor() ([]byte, []int) {
	return file_job_proto_rawDescGZIP(), []int{4}
}

type ListDeadJobsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Jobs          []*GetJobResponse      `protobuf:"bytes,1,rep,name=jobs,proto3" json:"jobs,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListDeadJobsResponse) Reset() {
	*x = ListDeadJobsResponse{}
	mi := &file_job_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListDeadJobsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListDeadJobsResponse) ProtoMessage() {}

func (x *ListDeadJobsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_job_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ListDeadJobsResponse) Descriptor() ([]byte, []int) {
	return file_job_proto_rawDescGZIP(), []int{5}
}

func (x *ListDeadJobsResponse) GetJobs() []*GetJobResponse {
	if x != nil {
		return x.Jobs
	}
	return nil
}

type HealthCheckRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HealthCheckRequest) Reset() {
	*x = HealthCheckRequest{}
	mi := &file_job_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HealthCheckRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthCheckRequest) ProtoMessage() {}

func (x *HealthCheckRequest) ProtoReflect() protoreflect.Message {
	mi := &file_job_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*HealthCheckRequest) Descriptor() ([]byte, []int) {
	return file_job_proto_rawDescGZIP(), []int{6}
}

type HealthCheckResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Status        string                 `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"` // "ok"
	QueueLength   int32                  `protobuf:"varint,2,opt,name=queue_length,json=queueLength,proto3" json:"queue_length,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HealthCheckResponse) Reset() {
	*x = HealthCheckResponse{}
	mi := &file_job_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HealthCheckResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthCheckResponse) ProtoMessage() {}

func (x *HealthCheckResponse) ProtoReflect() protoreflect.Message {
	mi := &file_job_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*HealthCheckResponse) Descriptor() ([]byte, []int) {
	return file_job_proto_rawDescGZIP(), []int{7}
}

func (x *HealthCheckResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *HealthCheckResponse) GetQueueLength() int32 {
	if x != nil {
		return x.QueueLength
	}
	return 0
}

var File_job_proto protoreflect.FileDescriptor

const file_job_proto_rawDesc = "" +
	"\n" +
	"\tjob.proto\x12\x05jobpb\"\xa2\x01\n" +
	"\x10SubmitJobRequest\x12\x12\n" +
	"\x04type\x18\x01 \x01(\tR\x04type\x12\x18\n" +
	"\apayload\x18\x02 \x01(\fR\apayload\x12\x1a\n" +
	"\bpriority\x18\x03 \x01(\x05R\bpriority\x12#\n" +
	"\rdelay_seconds\x18\x04 \x01(\x03R\fdelaySeconds\x12\x1f\n" +
	"\vmax_retries\x18\x05 \x01(\x05R\n" +
	"maxRetries\"B\n" +
	"\x11SubmitJobResponse\x12\x15\n" +
	"\x06job_id\x18\x01 \x01(\tR\x05jobId\x12\x16\n" +
	"\x06status\x18\x02 \x01(\tR\x06status\"&\n" +
	"\rGetJobRequest\x12\x15\n" +
	"\x06job_id\x18\x01 \x01(\tR\x05jobId\"\xc8\x01\n" +
	"\x0eGetJobResponse\x12\x15\n" +
	"\x06job_id\x18\x01 \x01(\tR\x05jobId\x12\x12\n" +
	"\x04type\x18\x02 \x01(\tR\x04type\x12\x16\n" +
	"\x06status\x18\x03 \x01(\tR\x06status\x12\x1f\n" +
	"\vretry_count\x18\x04 \x01(\x05R\n" +
	"retryCount\x12\x14\n" +
	"\x05error\x18\x05 \x01(\tR\x05error\x12\x1d\n" +
	"\n" +
	"created_at\x18\x06 \x01(\x03R\tcreatedAt\x12\x1d\n" +
	"\n" +
	"updated_at\x18\a \x01(\x03R\tupdatedAt\"\x15\n" +
	"\x13ListDeadJobsRequest\"A\n" +
	"\x14ListDeadJobsResponse\x12)\n" +
	"\x04jobs\x18\x01 \x03(\v2\x15.jobpb.GetJobResponseR\x04jobs\"\x14\n" +
	"\x12HealthCheckRequest\"P\n" +
	"\x13HealthCheckResponse\x12\x16\n" +
	"\x06status\x18\x01 \x01(\tR\x06status\x12!\n" +
	"\fqueue_length\x18\x02 \x01(\x05R\vqueueLength2\x92\x02\n" +
	"\n" +
	"JobService\x12>\n" +
	"\tSubmitJob\x12\x17.jobpb.SubmitJobRequest\x1a\x18.jobpb.SubmitJobResponse\x125\n" +
	"\x06GetJob\x12\x14.jobpb.GetJobRequest\x1a\x15.jobpb.GetJobResponse\x12G\n" +
	"\fListDeadJobs\x12\x1a.jobpb.ListDeadJobsRequest\x1a\x1b.jobpb.ListDeadJobsResponse\x12D\n" +
	"\vHealthCheck\x12\x19.jobpb.HealthCheckRequest\x1a\x1a.jobpb.HealthCheckResponseBIZGgithub.com/vihan/distributed-job-queue/internal/transport/grpc/pb;jobpbb\x06proto3"

var (
	file_job_proto_rawDescOnce sync.Once
	file_job_proto_rawDescData []byte
)

func file_job_proto_rawDescGZIP() []byte {
	file_job_proto_rawDescOnce.Do(func() {
		file_job_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_job_proto_rawDesc), len(file_job_proto_rawDesc)))
	})
	return file_job_proto_rawDescData
}

var file_job_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_job_proto_goTypes = []any{
	(*SubmitJobRequest)(nil),     // 0: jobpb.SubmitJobRequest
	(*SubmitJobResponse)(nil),    // 1: jobpb.SubmitJobResponse
	(*GetJobRequest)(nil),        // 2: jobpb.GetJobRequest
	(*GetJobResponse)(nil),       // 3: jobpb.GetJobResponse
	(*ListDeadJobsRequest)(nil),  // 4: jobpb.ListDeadJobsRequest
	(*ListDeadJobsResponse)(nil), // 5: jobpb.ListDeadJobsResponse
	(*HealthCheckRequest)(nil),   // 6: jobpb.HealthCheckRequest
	(*HealthCheckResponse)(nil),  // 7: jobpb.HealthCheckResponse
}
var file_job_proto_depIdxs = []int32{
	3, // 0: jobpb.ListDeadJobsResponse.jobs:type_name -> jobpb.GetJobResponse
	0, // 1: jobpb.JobService.SubmitJob:input_type -> jobpb.SubmitJobRequest
	2, // 2: jobpb.JobService.GetJob:input_type -> jobpb.GetJobRequest
	4, // 3: jobpb.JobService.ListDeadJobs:input_type -> jobpb.ListDeadJobsRequest
	6, // 4: jobpb.JobService.HealthCheck:input_type -> jobpb.HealthCheckRequest
	1, // 5: jobpb.JobService.SubmitJob:output_type -> jobpb.SubmitJobResponse
	3, // 6: jobpb.JobService.GetJob:output_type -> jobpb.GetJobResponse
	5, // 7: jobpb.JobService.ListDeadJobs:output_type -> jobpb.ListDeadJobsResponse
	7, // 8: jobpb.JobService.HealthCheck:output_type -> jobpb.HealthCheckResponse
	5, // [5:9] is the sub-list for method output_type
	1, // [1:5] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_job_proto_init() }
func file_job_proto_init() {
	if File_job_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_job_proto_rawDesc), len(file_job_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_job_proto_goTypes,
		DependencyIndexes: file_job_proto_depIdxs,
		MessageInfos:      file_job_proto_msgTypes,
	}.Build()
	File_job_proto = out.File
	file_job_proto_goTypes = nil
	file_job_proto_depIdxs = nil
}
