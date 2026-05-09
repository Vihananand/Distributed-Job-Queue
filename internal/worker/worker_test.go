package worker_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/vihan/distributed-job-queue/internal/models"
	"github.com/vihan/distributed-job-queue/internal/queue"
	"github.com/vihan/distributed-job-queue/internal/utils"
	"github.com/vihan/distributed-job-queue/internal/worker"
)

func TestWorker_ProcessesJob(t *testing.T) {
	q := queue.NewMemoryQueue()
	reg := worker.NewRegistry()

	var called atomic.Bool
	reg.Register("test", func(_ context.Context, _ *models.Job) error {
		called.Store(true)
		return nil
	})

	metrics := utils.NewMetrics()
	w := worker.New(1, q, reg, metrics)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	job := models.NewJob("test", []byte(`{}`), 0, 0, 0)
	_ = q.Push(ctx, job)

	go w.Run(ctx)

	time.Sleep(500 * time.Millisecond)
	if !called.Load() {
		t.Error("handler was not called")
	}
}

func TestWorker_FailsJobToRetry(t *testing.T) {
	q := queue.NewMemoryQueue()
	reg := worker.NewRegistry()

	var attempts atomic.Int32
	reg.Register("flaky", func(_ context.Context, _ *models.Job) error {
		attempts.Add(1)
		if attempts.Load() < 3 {
			return errors.New("temporary error")
		}
		return nil
	})

	metrics := utils.NewMetrics()
	pool := worker.NewPool(1, q, reg, metrics)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	job := models.NewJob("flaky", []byte(`{}`), 0, 5, 0)
	_ = q.Push(ctx, job)

	go pool.Start(ctx)
	time.Sleep(8 * time.Second)

	if attempts.Load() < 3 {
		t.Errorf("expected at least 3 attempts, got %d", attempts.Load())
	}
}

func TestRegistry_UnknownType(t *testing.T) {
	reg := worker.NewRegistry()
	_, err := reg.Get("nonexistent")
	if err == nil {
		t.Error("expected error for unregistered job type")
	}
}
