package storage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/vihan/distributed-job-queue/internal/models"
)

const redisJobPrefix = "djq:store:"

type RedisStore struct {
	rdb *redis.Client
}

func NewRedisStore(rdb *redis.Client) *RedisStore {
	return &RedisStore{rdb: rdb}
}

func (s *RedisStore) key(id string) string { return redisJobPrefix + id }

func (s *RedisStore) Save(job *models.Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return s.rdb.Set(context.Background(), s.key(job.ID), data, 0).Err()
}

func (s *RedisStore) Get(jobID string) (*models.Job, error) {
	data, err := s.rdb.Get(context.Background(), s.key(jobID)).Bytes()
	if err == redis.Nil {
		return nil, fmt.Errorf("redisstore: job %q not found", jobID)
	}
	if err != nil {
		return nil, err
	}
	var job models.Job
	return &job, json.Unmarshal(data, &job)
}

func (s *RedisStore) Update(job *models.Job) error {
	return s.Save(job)
}

func (s *RedisStore) List(status models.Status) ([]*models.Job, error) {
	ctx := context.Background()
	var keys []string
	iter := s.rdb.Scan(ctx, 0, redisJobPrefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}

	var out []*models.Job
	for _, k := range keys {
		data, err := s.rdb.Get(ctx, k).Bytes()
		if err != nil {
			continue
		}
		var job models.Job
		if err := json.Unmarshal(data, &job); err != nil {
			continue
		}
		if job.Status == status {
			out = append(out, &job)
		}
	}
	return out, nil
}

func (s *RedisStore) Delete(jobID string) error {
	return s.rdb.Del(context.Background(), s.key(jobID)).Err()
}
