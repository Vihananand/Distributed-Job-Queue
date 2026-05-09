package storage

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/vihan/distributed-job-queue/internal/models"
)

type MemoryStore struct {
	mu   sync.RWMutex
	jobs map[string]*models.Job
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{jobs: make(map[string]*models.Job)}
}

func (s *MemoryStore) Save(job *models.Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	cp := copyJob(job)
	s.jobs[cp.ID] = cp
	return nil
}

func (s *MemoryStore) Get(jobID string) (*models.Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	job, ok := s.jobs[jobID]
	if !ok {
		return nil, fmt.Errorf("memorystore: job %q not found", jobID)
	}
	return copyJob(job), nil
}

func (s *MemoryStore) Update(job *models.Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.jobs[job.ID]; !ok {
		return fmt.Errorf("memorystore: job %q not found for update", job.ID)
	}
	s.jobs[job.ID] = copyJob(job)
	return nil
}

func (s *MemoryStore) List(status models.Status) ([]*models.Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []*models.Job
	for _, j := range s.jobs {
		if j.Status == status {
			out = append(out, copyJob(j))
		}
	}
	return out, nil
}

func (s *MemoryStore) Delete(jobID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.jobs, jobID)
	return nil
}

func copyJob(j *models.Job) *models.Job {
	data, _ := json.Marshal(j)
	var cp models.Job
	_ = json.Unmarshal(data, &cp)
	return &cp
}
