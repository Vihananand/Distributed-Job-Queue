package queue

import (
	"container/heap"
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/vihan/distributed-job-queue/internal/models"
)


type item struct {
	job   *models.Job
	index int
}

type pHeap []*item

func (h pHeap) Len() int            { return len(h) }
func (h pHeap) Less(i, j int) bool  { return h[i].job.Priority > h[j].job.Priority }
func (h pHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i]; h[i].index = i; h[j].index = j }
func (h *pHeap) Push(x interface{}) { n := len(*h); it := x.(*item); it.index = n; *h = append(*h, it) }
func (h *pHeap) Pop() interface{}   { old := *h; n := len(old); it := old[n-1]; old[n-1] = nil; *h = old[:n-1]; return it }

type DLQ struct {
	mu   sync.RWMutex
	jobs []*models.Job
}

func (d *DLQ) Push(job *models.Job) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.jobs = append(d.jobs, job)
}

func (d *DLQ) List() []*models.Job {
	d.mu.RLock()
	defer d.mu.RUnlock()
	out := make([]*models.Job, len(d.jobs))
	copy(out, d.jobs)
	return out
}

type MemoryQueue struct {
	mu         sync.Mutex
	h          pHeap
	processing map[string]*models.Job
	signal     chan struct{} 
	DLQ        *DLQ
}

func NewMemoryQueue() *MemoryQueue {
	mq := &MemoryQueue{
		processing: make(map[string]*models.Job),
		signal:     make(chan struct{}, 1),
		DLQ:        &DLQ{},
	}
	heap.Init(&mq.h)
	return mq
}

func (mq *MemoryQueue) Push(_ context.Context, job *models.Job) error {
	mq.mu.Lock()
	heap.Push(&mq.h, &item{job: job})
	mq.mu.Unlock()
	select {
	case mq.signal <- struct{}{}:
	default:
	}
	return nil
}

func (mq *MemoryQueue) Pop(ctx context.Context) (*models.Job, error) {
	for {
		mq.mu.Lock()
		now := time.Now()
		for i := 0; i < mq.h.Len(); i++ {
			candidate := mq.h[i]
			if !now.Before(candidate.job.RunAt) {
				heap.Remove(&mq.h, candidate.index)
				job := candidate.job
				job.Status = models.StatusRunning
				job.UpdatedAt = now
				mq.processing[job.ID] = job
				mq.mu.Unlock()
				return job, nil
			}
		}
		mq.mu.Unlock()

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-mq.signal:
		case <-time.After(250 * time.Millisecond):
		}
	}
}

func (mq *MemoryQueue) Ack(_ context.Context, jobID string) error {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	if _, ok := mq.processing[jobID]; !ok {
		return fmt.Errorf("memoryqueue: job %q not in processing state", jobID)
	}
	delete(mq.processing, jobID)
	return nil
}

func (mq *MemoryQueue) Fail(ctx context.Context, jobID string, jobErr error) error {
	mq.mu.Lock()
	job, ok := mq.processing[jobID]
	if !ok {
		mq.mu.Unlock()
		return fmt.Errorf("memoryqueue: job %q not in processing state", jobID)
	}
	delete(mq.processing, jobID)
	mq.mu.Unlock()

	job.RetryCount++
	job.Error = jobErr.Error()
	job.UpdatedAt = time.Now()

	if job.RetryCount > job.MaxRetries {
		job.Status = models.StatusDead
		mq.DLQ.Push(job)
		return nil
	}

	backoff := time.Duration(math.Pow(2, float64(job.RetryCount))) * time.Second
	if backoff > 5*time.Minute {
		backoff = 5 * time.Minute
	}
	job.RunAt = time.Now().Add(backoff)
	job.Status = models.StatusQueued
	return mq.Push(ctx, job)
}

func (mq *MemoryQueue) Len(_ context.Context) (int, error) {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	return mq.h.Len(), nil
}
