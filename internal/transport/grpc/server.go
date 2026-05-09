// Package grpcsrv implements the gRPC JobService server.
package grpcsrv

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vihan/distributed-job-queue/internal/producer"
	pb "github.com/vihan/distributed-job-queue/internal/transport/grpc/pb"
	"github.com/vihan/distributed-job-queue/internal/queue"
	"github.com/vihan/distributed-job-queue/internal/storage"
)

type Server struct {
	pb.UnimplementedJobServiceServer
	producer *producer.Producer
	store    storage.Storage
	q        queue.Queue
}

func New(p *producer.Producer, store storage.Storage, q queue.Queue) *Server {
	return &Server{producer: p, store: store, q: q}
}

func (s *Server) SubmitJob(ctx context.Context, req *pb.SubmitJobRequest) (*pb.SubmitJobResponse, error) {
	job, err := s.producer.Submit(ctx, producer.SubmitRequest{
		Type:         req.Type,
		Payload:      req.Payload,
		Priority:     int(req.Priority),
		DelaySeconds: req.DelaySeconds,
		MaxRetries:   int(req.MaxRetries),
	})
	if err != nil {
		slog.Warn("SubmitJob failed", "err", err)
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	return &pb.SubmitJobResponse{
		JobId:  job.ID,
		Status: string(job.Status),
	}, nil
}

func (s *Server) GetJob(ctx context.Context, req *pb.GetJobRequest) (*pb.GetJobResponse, error) {
	if req.JobId == "" {
		return nil, status.Error(codes.InvalidArgument, "job_id is required")
	}
	job, err := s.store.Get(req.JobId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "job not found: %v", err)
	}
	return &pb.GetJobResponse{
		JobId:      job.ID,
		Type:       job.Type,
		Status:     string(job.Status),
		RetryCount: int32(job.RetryCount),
		Error:      job.Error,
		CreatedAt:  job.CreatedAt.Unix(),
		UpdatedAt:  job.UpdatedAt.Unix(),
	}, nil
}

func (s *Server) ListDeadJobs(ctx context.Context, _ *pb.ListDeadJobsRequest) (*pb.ListDeadJobsResponse, error) {
	jobs, err := s.store.List("dead")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list dead jobs: %v", err)
	}
	var pbJobs []*pb.GetJobResponse
	for _, job := range jobs {
		pbJobs = append(pbJobs, &pb.GetJobResponse{
			JobId:      job.ID,
			Type:       job.Type,
			Status:     string(job.Status),
			RetryCount: int32(job.RetryCount),
			Error:      job.Error,
			CreatedAt:  job.CreatedAt.Unix(),
			UpdatedAt:  job.UpdatedAt.Unix(),
		})
	}
	return &pb.ListDeadJobsResponse{Jobs: pbJobs}, nil
}

func (s *Server) HealthCheck(ctx context.Context, _ *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	n, err := s.q.Len(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "queue len: %v", err)
	}
	return &pb.HealthCheckResponse{
		Status:      "ok",
		QueueLength: int32(n),
	}, nil
}
