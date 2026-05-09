// Package worker provides the handler registry, single worker, and worker pool.
package worker

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/vihan/distributed-job-queue/internal/models"
	"github.com/vihan/distributed-job-queue/internal/queue"
	"github.com/vihan/distributed-job-queue/internal/utils"
)

type HandlerFunc func(ctx context.Context, job *models.Job) error
type Registry struct {
	mu       sync.RWMutex
	handlers map[string]HandlerFunc
}

func NewRegistry() *Registry {
	return &Registry{handlers: make(map[string]HandlerFunc)}
}

func (r *Registry) Register(jobType string, fn HandlerFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[jobType] = fn
}

func (r *Registry) Get(jobType string) (HandlerFunc, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	h, ok := r.handlers[jobType]
	if !ok {
		return nil, fmt.Errorf("worker: no handler registered for type %q", jobType)
	}
	return h, nil
}

type Worker struct {
	id       int
	q        queue.Queue
	registry *Registry
	metrics  *utils.Metrics
	log      *slog.Logger
}

func New(id int, q queue.Queue, r *Registry, m *utils.Metrics) *Worker {
	return &Worker{
		id:       id,
		q:        q,
		registry: r,
		metrics:  m,
		log:      slog.Default().With("worker_id", id),
	}
}

func (w *Worker) Run(ctx context.Context) {
	w.log.Info("worker started")
	for {
		job, err := w.q.Pop(ctx)
		if err != nil {
			if ctx.Err() != nil {
				w.log.Info("worker shutting down")
				return
			}
			w.log.Error("pop error", "err", err)
			continue
		}
		w.process(ctx, job)
	}
}

func (w *Worker) process(ctx context.Context, job *models.Job) {
	log := w.log.With("job_id", job.ID, "job_type", job.Type)
	log.Info("processing job")

	handler, err := w.registry.Get(job.Type)
	if err != nil {
		log.Warn("no handler, failing job", "err", err)
		_ = w.q.Fail(ctx, job.ID, err)
		w.metrics.JobsFailed.Inc()
		return
	}

	if err := handler(ctx, job); err != nil {
		log.Warn("job failed", "err", err, "retry", job.RetryCount)
		_ = w.q.Fail(ctx, job.ID, err)
		w.metrics.JobsFailed.Inc()
		return
	}

	if err := w.q.Ack(ctx, job.ID); err != nil {
		log.Error("ack failed", "err", err)
		return
	}
	log.Info("job completed")
	w.metrics.JobsProcessed.Inc()
}
