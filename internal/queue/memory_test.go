package queue_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/vihan/distributed-job-queue/internal/models"
	"github.com/vihan/distributed-job-queue/internal/queue"
)

func newJob(t string, priority int) *models.Job {
	return models.NewJob(t, []byte(`{"hello":"world"}`), priority, 3, 0)
}

func TestMemoryQueue_PushPop(t *testing.T) {
	q := queue.NewMemoryQueue()
	ctx := context.Background()

	job := newJob("echo", 0)
	if err := q.Push(ctx, job); err != nil {
		t.Fatalf("push: %v", err)
	}

	popCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	got, err := q.Pop(popCtx)
	if err != nil {
		t.Fatalf("pop: %v", err)
	}
	if got.ID != job.ID {
		t.Errorf("got job %q, want %q", got.ID, job.ID)
	}
	if got.Status != models.StatusRunning {
		t.Errorf("expected StatusRunning, got %s", got.Status)
	}
}

func TestMemoryQueue_Priority(t *testing.T) {
	q := queue.NewMemoryQueue()
	ctx := context.Background()

	low := newJob("low", 1)
	high := newJob("high", 10)

	_ = q.Push(ctx, low)
	_ = q.Push(ctx, high)

	popCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	first, _ := q.Pop(popCtx)
	if first.ID != high.ID {
		t.Errorf("expected high-priority job first, got %q", first.Type)
	}
}

func TestMemoryQueue_Ack(t *testing.T) {
	q := queue.NewMemoryQueue()
	ctx := context.Background()

	job := newJob("echo", 0)
	_ = q.Push(ctx, job)

	popCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	got, _ := q.Pop(popCtx)
	if err := q.Ack(ctx, got.ID); err != nil {
		t.Errorf("ack: %v", err)
	}
}

func TestMemoryQueue_RetryOnFail(t *testing.T) {
	q := queue.NewMemoryQueue()
	ctx := context.Background()

	job := newJob("echo", 0)
	job.MaxRetries = 2
	_ = q.Push(ctx, job)

	popCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	got, _ := q.Pop(popCtx)
	if err := q.Fail(ctx, got.ID, errors.New("boom")); err != nil {
		t.Errorf("fail: %v", err)
	}

	if q.DLQ.List() != nil {
		if len(q.DLQ.List()) > 0 {
			t.Error("job should not be dead-lettered after first failure")
		}
	}
}

func TestMemoryQueue_DeadLetter(t *testing.T) {
	q := queue.NewMemoryQueue()
	ctx := context.Background()

	job := newJob("echo", 0)
	job.MaxRetries = 0
	_ = q.Push(ctx, job)

	popCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	got, _ := q.Pop(popCtx)
	_ = q.Fail(ctx, got.ID, errors.New("permanent failure"))

	dead := q.DLQ.List()
	if len(dead) != 1 {
		t.Errorf("expected 1 dead job, got %d", len(dead))
	}
	if dead[0].Status != models.StatusDead {
		t.Errorf("expected StatusDead, got %s", dead[0].Status)
	}
}

func TestMemoryQueue_DelayedJob(t *testing.T) {
	q := queue.NewMemoryQueue()
	ctx := context.Background()

	job := models.NewJob("delayed", []byte(`{}`), 0, 0, 200*time.Millisecond)
	_ = q.Push(ctx, job)

	quickCtx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
	defer cancel()

	_, err := q.Pop(quickCtx)
	if err == nil {
		t.Error("delayed job should not be available immediately")
	}

	time.Sleep(300 * time.Millisecond)
	lateCtx, cancel2 := context.WithTimeout(ctx, time.Second)
	defer cancel2()

	got, err := q.Pop(lateCtx)
	if err != nil {
		t.Fatalf("delayed job not available after delay: %v", err)
	}
	if got.ID != job.ID {
		t.Errorf("unexpected job ID")
	}
}

func BenchmarkMemoryQueue_PushPop(b *testing.B) {
	q := queue.NewMemoryQueue()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		job := newJob("bench", 0)
		_ = q.Push(ctx, job)
		popCtx, cancel := context.WithTimeout(ctx, time.Second)
		got, _ := q.Pop(popCtx)
		_ = q.Ack(ctx, got.ID)
		cancel()
	}
}
