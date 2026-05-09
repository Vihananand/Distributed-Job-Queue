package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vihan/distributed-job-queue/internal/models"
)

const (
	pendingKey    = "djq:pending"    
	processingKey = "djq:processing" 
	delayedKey    = "djq:delayed"    
	jobHashPrefix = "djq:job:"       
	dlqKey        = "djq:dlq"       
)

type RedisQueue struct {
	rdb *redis.Client
}
func NewRedisQueue(rdb *redis.Client) *RedisQueue {
	return &RedisQueue{rdb: rdb}
}

func jobKey(id string) string { return jobHashPrefix + id }

func (q *RedisQueue) Push(ctx context.Context, job *models.Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("redisqueue push marshal: %w", err)
	}
	pipe := q.rdb.Pipeline()
	pipe.Set(ctx, jobKey(job.ID), data, 0)

	if time.Now().Before(job.RunAt) {
		pipe.ZAdd(ctx, delayedKey, redis.Z{Score: float64(job.RunAt.Unix()), Member: job.ID})
	} else {
		pipe.LPush(ctx, pendingKey, job.ID)
	}
	_, err = pipe.Exec(ctx)
	return err
}

func (q *RedisQueue) Pop(ctx context.Context) (*models.Job, error) {
	for {
		if err := q.promoteDelayed(ctx); err != nil {
			return nil, err
		}
		res, err := q.rdb.BRPopLPush(ctx, pendingKey, processingKey, time.Second).Result()
		if err == redis.Nil {
			continue 
		}
		if err != nil {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			return nil, fmt.Errorf("redisqueue pop: %w", err)
		}

		job, err := q.loadJob(ctx, res)
		if err != nil {
			return nil, err
		}
		job.Status = models.StatusRunning
		job.UpdatedAt = time.Now()
		if err := q.saveJob(ctx, job); err != nil {
			return nil, err
		}
		return job, nil
	}
}

func (q *RedisQueue) Ack(ctx context.Context, jobID string) error {
	job, err := q.loadJob(ctx, jobID)
	if err != nil {
		return err
	}
	job.Status = models.StatusSuccess
	job.UpdatedAt = time.Now()

	pipe := q.rdb.Pipeline()
	pipe.LRem(ctx, processingKey, 1, jobID)
	data, _ := json.Marshal(job)
	pipe.Set(ctx, jobKey(jobID), data, 0)
	_, err = pipe.Exec(ctx)
	return err
}

func (q *RedisQueue) Fail(ctx context.Context, jobID string, jobErr error) error {
	job, err := q.loadJob(ctx, jobID)
	if err != nil {
		return err
	}
	job.RetryCount++
	job.Error = jobErr.Error()
	job.UpdatedAt = time.Now()

	if err := q.rdb.LRem(ctx, processingKey, 1, jobID).Err(); err != nil {
		return err
	}

	if job.RetryCount > job.MaxRetries {
		job.Status = models.StatusDead
		if err := q.saveJob(ctx, job); err != nil {
			return err
		}
		return q.rdb.LPush(ctx, dlqKey, jobID).Err()
	}

	backoff := time.Duration(math.Pow(2, float64(job.RetryCount))) * time.Second
	if backoff > 5*time.Minute {
		backoff = 5 * time.Minute
	}
	job.RunAt = time.Now().Add(backoff)
	job.Status = models.StatusQueued
	return q.Push(ctx, job)
}

func (q *RedisQueue) Len(ctx context.Context) (int, error) {
	n, err := q.rdb.LLen(ctx, pendingKey).Result()
	return int(n), err
}

func (q *RedisQueue) loadJob(ctx context.Context, id string) (*models.Job, error) {
	data, err := q.rdb.Get(ctx, jobKey(id)).Bytes()
	if err != nil {
		return nil, fmt.Errorf("redisqueue load job %q: %w", id, err)
	}
	var job models.Job
	if err := json.Unmarshal(data, &job); err != nil {
		return nil, fmt.Errorf("redisqueue unmarshal job %q: %w", id, err)
	}
	return &job, nil
}

func (q *RedisQueue) saveJob(ctx context.Context, job *models.Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return q.rdb.Set(ctx, jobKey(job.ID), data, 0).Err()
}

func (q *RedisQueue) promoteDelayed(ctx context.Context) error {
	now := float64(time.Now().Unix())
	ids, err := q.rdb.ZRangeByScore(ctx, delayedKey, &redis.ZRangeBy{
		Min: "-inf",
		Max: fmt.Sprintf("%f", now),
	}).Result()
	if err != nil || len(ids) == 0 {
		return err
	}
	pipe := q.rdb.Pipeline()
	for _, id := range ids {
		pipe.ZRem(ctx, delayedKey, id)
		pipe.LPush(ctx, pendingKey, id)
	}
	_, err = pipe.Exec(ctx)
	return err
}

func (q *RedisQueue) ListDLQ(ctx context.Context) ([]*models.Job, error) {
	ids, err := q.rdb.LRange(ctx, dlqKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	jobs := make([]*models.Job, 0, len(ids))
	for _, id := range ids {
		job, err := q.loadJob(ctx, id)
		if err != nil {
			continue
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}
