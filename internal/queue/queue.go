package queue

import (
	"context"

	"github.com/vihan/distributed-job-queue/internal/models"
)

type Queue interface {
	Push(ctx context.Context, job *models.Job) error
	Pop(ctx context.Context) (*models.Job, error)
	Ack(ctx context.Context, jobID string) error
	Fail(ctx context.Context, jobID string, jobErr error) error
	Len(ctx context.Context) (int, error)
}
