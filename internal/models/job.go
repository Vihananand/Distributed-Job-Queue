package models

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusQueued  Status = "queued"
	StatusRunning Status = "running"
	StatusSuccess Status = "success"
	StatusFailed  Status = "failed"
	StatusDead    Status = "dead" 
)

type Job struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`      
	Payload    []byte    `json:"payload"`    
	Status     Status    `json:"status"`
	Priority   int       `json:"priority"`   
	RetryCount int       `json:"retry_count"`
	MaxRetries int       `json:"max_retries"`
	RunAt      time.Time `json:"run_at"`    
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Error      string    `json:"error,omitempty"`
}

func NewJob(jobType string, payload []byte, priority, maxRetries int, delay time.Duration) *Job {
	now := time.Now()
	runAt := now
	if delay > 0 {
		runAt = now.Add(delay)
	}
	return &Job{
		ID:         uuid.New().String(),
		Type:       jobType,
		Payload:    payload,
		Status:     StatusQueued,
		Priority:   priority,
		MaxRetries: maxRetries,
		RunAt:      runAt,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
