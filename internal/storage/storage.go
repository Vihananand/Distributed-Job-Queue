package storage

import (
	"github.com/vihan/distributed-job-queue/internal/models"
)

type Storage interface {
	Save(job *models.Job) error
	Get(jobID string) (*models.Job, error)
	Update(job *models.Job) error
	List(status models.Status) ([]*models.Job, error)
	Delete(jobID string) error
}
