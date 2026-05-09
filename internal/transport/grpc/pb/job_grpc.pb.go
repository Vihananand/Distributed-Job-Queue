package jobpb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const _ = grpc.SupportPackageIsVersion9

const (
	JobService_SubmitJob_FullMethodName    = "/jobpb.JobService/SubmitJob"
	JobService_GetJob_FullMethodName       = "/jobpb.JobService/GetJob"
	JobService_ListDeadJobs_FullMethodName = "/jobpb.JobService/ListDeadJobs"
	JobService_HealthCheck_FullMethodName  = "/jobpb.JobService/HealthCheck"
)

type JobServiceClient interface {
	SubmitJob(ctx context.Context, in *SubmitJobRequest, opts ...grpc.CallOption) (*SubmitJobResponse, error)
	GetJob(ctx context.Context, in *GetJobRequest, opts ...grpc.CallOption) (*GetJobResponse, error)
	ListDeadJobs(ctx context.Context, in *ListDeadJobsRequest, opts ...grpc.CallOption) (*ListDeadJobsResponse, error)
	HealthCheck(ctx context.Context, in *HealthCheckRequest, opts ...grpc.CallOption) (*HealthCheckResponse, error)
}

type jobServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewJobServiceClient(cc grpc.ClientConnInterface) JobServiceClient {
	return &jobServiceClient{cc}
}

func (c *jobServiceClient) SubmitJob(ctx context.Context, in *SubmitJobRequest, opts ...grpc.CallOption) (*SubmitJobResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SubmitJobResponse)
	err := c.cc.Invoke(ctx, JobService_SubmitJob_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jobServiceClient) GetJob(ctx context.Context, in *GetJobRequest, opts ...grpc.CallOption) (*GetJobResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetJobResponse)
	err := c.cc.Invoke(ctx, JobService_GetJob_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jobServiceClient) ListDeadJobs(ctx context.Context, in *ListDeadJobsRequest, opts ...grpc.CallOption) (*ListDeadJobsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListDeadJobsResponse)
	err := c.cc.Invoke(ctx, JobService_ListDeadJobs_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jobServiceClient) HealthCheck(ctx context.Context, in *HealthCheckRequest, opts ...grpc.CallOption) (*HealthCheckResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(HealthCheckResponse)
	err := c.cc.Invoke(ctx, JobService_HealthCheck_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type JobServiceServer interface {
	SubmitJob(context.Context, *SubmitJobRequest) (*SubmitJobResponse, error)
	GetJob(context.Context, *GetJobRequest) (*GetJobResponse, error)
	ListDeadJobs(context.Context, *ListDeadJobsRequest) (*ListDeadJobsResponse, error)
	HealthCheck(context.Context, *HealthCheckRequest) (*HealthCheckResponse, error)
	mustEmbedUnimplementedJobServiceServer()
}

type UnimplementedJobServiceServer struct{}

func (UnimplementedJobServiceServer) SubmitJob(context.Context, *SubmitJobRequest) (*SubmitJobResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method SubmitJob not implemented")
}
func (UnimplementedJobServiceServer) GetJob(context.Context, *GetJobRequest) (*GetJobResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method GetJob not implemented")
}
func (UnimplementedJobServiceServer) ListDeadJobs(context.Context, *ListDeadJobsRequest) (*ListDeadJobsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method ListDeadJobs not implemented")
}
func (UnimplementedJobServiceServer) HealthCheck(context.Context, *HealthCheckRequest) (*HealthCheckResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method HealthCheck not implemented")
}
func (UnimplementedJobServiceServer) mustEmbedUnimplementedJobServiceServer() {}
func (UnimplementedJobServiceServer) testEmbeddedByValue()                    {}

type UnsafeJobServiceServer interface {
	mustEmbedUnimplementedJobServiceServer()
}

func RegisterJobServiceServer(s grpc.ServiceRegistrar, srv JobServiceServer) {
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&JobService_ServiceDesc, srv)
}

func _JobService_SubmitJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubmitJobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobServiceServer).SubmitJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: JobService_SubmitJob_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobServiceServer).SubmitJob(ctx, req.(*SubmitJobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _JobService_GetJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetJobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobServiceServer).GetJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: JobService_GetJob_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobServiceServer).GetJob(ctx, req.(*GetJobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _JobService_ListDeadJobs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListDeadJobsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobServiceServer).ListDeadJobs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: JobService_ListDeadJobs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobServiceServer).ListDeadJobs(ctx, req.(*ListDeadJobsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _JobService_HealthCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HealthCheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobServiceServer).HealthCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: JobService_HealthCheck_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobServiceServer).HealthCheck(ctx, req.(*HealthCheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var JobService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "jobpb.JobService",
	HandlerType: (*JobServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SubmitJob",
			Handler:    _JobService_SubmitJob_Handler,
		},
		{
			MethodName: "GetJob",
			Handler:    _JobService_GetJob_Handler,
		},
		{
			MethodName: "ListDeadJobs",
			Handler:    _JobService_ListDeadJobs_Handler,
		},
		{
			MethodName: "HealthCheck",
			Handler:    _JobService_HealthCheck_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "job.proto",
}
