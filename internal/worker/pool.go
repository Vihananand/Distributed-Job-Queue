package worker

import (
	"context"
	"log/slog"
	"sync"

	"github.com/vihan/distributed-job-queue/internal/queue"
	"github.com/vihan/distributed-job-queue/internal/utils"
)

type Pool struct {
	size     int
	q        queue.Queue
	registry *Registry
	metrics  *utils.Metrics
	wg       sync.WaitGroup
}

func NewPool(size int, q queue.Queue, r *Registry, m *utils.Metrics) *Pool {
	return &Pool{size: size, q: q, registry: r, metrics: m}
}

func (p *Pool) Start(ctx context.Context) {
	slog.Info("worker pool starting", "size", p.size)
	p.wg.Add(p.size)
	for i := range p.size {
		w := New(i+1, p.q, p.registry, p.metrics)
		go func(w *Worker) {
			defer p.wg.Done()
			w.Run(ctx)
		}(w)
	}
	<-ctx.Done()
	slog.Info("worker pool draining…")
	p.wg.Wait()
	slog.Info("worker pool stopped")
}
