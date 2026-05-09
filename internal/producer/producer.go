package producer

import (
	"context"
	"fmt"
	"time"

	"github.com/vihan/distributed-job-queue/internal/models"
	"github.com/vihan/distributed-job-queue/internal/queue"
	"github.com/vihan/distributed-job-queue/internal/storage"
	"github.com/vihan/distributed-job-queue/internal/utils"
)

type SubmitRequest struct {
	Type         string
	Payload      []byte
	Priority     int
	DelaySeconds int64
	MaxRetries   int
}

type Producer struct {
	q       queue.Queue
	store   storage.Storage
	metrics *utils.Metrics
}

func New(q queue.Queue, store storage.Storage, m *utils.Metrics) *Producer {
	return &Producer{q: q, store: store, metrics: m}
}

func (p *Producer) Submit(ctx context.Context, req SubmitRequest) (*models.Job, error) {
	if err := validate(req); err != nil {
		return nil, fmt.Errorf("producer: invalid request: %w", err)
	}

	maxRetries := req.MaxRetries
	if maxRetries <= 0 {
		maxRetries = 3
	}

	job := models.NewJob(
		req.Type,
		req.Payload,
		req.Priority,
		maxRetries,
		time.Duration(req.DelaySeconds)*time.Second,
	)

	if err := p.store.Save(job); err != nil {
		return nil, fmt.Errorf("producer: store save: %w", err)
	}
	if err := p.q.Push(ctx, job); err != nil {
		return nil, fmt.Errorf("producer: queue push: %w", err)
	}
	p.metrics.JobsPushed.Inc()
	return job, nil
}

func validate(req SubmitRequest) error {
	if req.Type == "" {
		return fmt.Errorf("type is required")
	}
	if len(req.Payload) == 0 {
		return fmt.Errorf("payload must not be empty")
	}
	if req.Priority < 0 {
		return fmt.Errorf("priority must be >= 0")
	}
	if req.DelaySeconds < 0 {
		return fmt.Errorf("delay_seconds must be >= 0")
	}
	return nil
}
