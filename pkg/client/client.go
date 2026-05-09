package client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/vihan/distributed-job-queue/internal/transport/grpc/pb"
)

type Client struct {
	conn *grpc.ClientConn
	svc  pb.JobServiceClient
}

func New(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("client dial %q: %w", addr, err)
	}
	return &Client{conn: conn, svc: pb.NewJobServiceClient(conn)}, nil
}

func (c *Client) Close() error { return c.conn.Close() }

func (c *Client) Submit(ctx context.Context, jobType string, payload []byte, opts ...Option) (string, error) {
	req := &pb.SubmitJobRequest{
		Type:       jobType,
		Payload:    payload,
		MaxRetries: 3,
	}
	for _, o := range opts {
		o(req)
	}
	resp, err := c.svc.SubmitJob(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.JobId, nil
}

func (c *Client) GetJob(ctx context.Context, jobID string) (*pb.GetJobResponse, error) {
	return c.svc.GetJob(ctx, &pb.GetJobRequest{JobId: jobID})
}

func (c *Client) Health(ctx context.Context) (*pb.HealthCheckResponse, error) {
	return c.svc.HealthCheck(ctx, &pb.HealthCheckRequest{})
}

type Option func(*pb.SubmitJobRequest)

func WithPriority(p int32) Option        { return func(r *pb.SubmitJobRequest) { r.Priority = p } }
func WithDelay(seconds int64) Option     { return func(r *pb.SubmitJobRequest) { r.DelaySeconds = seconds } }
func WithMaxRetries(n int32) Option      { return func(r *pb.SubmitJobRequest) { r.MaxRetries = n } }
