package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/redis/go-redis/v9"

	"github.com/vihan/distributed-job-queue/internal/config"
	"github.com/vihan/distributed-job-queue/internal/models"
	inqueue "github.com/vihan/distributed-job-queue/internal/queue"
	"github.com/vihan/distributed-job-queue/internal/utils"
	"github.com/vihan/distributed-job-queue/internal/worker"
)

func main() {
	cfg := config.Load()
	utils.InitLogger(cfg.Env)

	metrics := utils.NewMetrics()

	// Queue
	var q inqueue.Queue
	switch cfg.QueueBackend {
	case "memory":
		slog.Info("worker: using in-memory queue")
		q = inqueue.NewMemoryQueue()
	default:
		slog.Info("worker: connecting to Redis", "addr", cfg.RedisAddr)
		rdb := redis.NewClient(&redis.Options{
			Addr:     cfg.RedisAddr,
			Password: cfg.RedisPassword,
			DB:       cfg.RedisDB,
		})
		if err := rdb.Ping(context.Background()).Err(); err != nil {
			slog.Error("Redis ping failed", "err", err)
			os.Exit(1)
		}
		q = inqueue.NewRedisQueue(rdb)
	}

	// Handler Registry
	registry := worker.NewRegistry()

	// Register the built-in "echo" handler (logs payload, always succeeds).
	registry.Register("echo", func(ctx context.Context, job *models.Job) error {
		slog.Info("echo handler", "job_id", job.ID, "payload", string(job.Payload))
		return nil
	})

	// Worker Pool
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool := worker.NewPool(cfg.WorkerCount, q, registry, metrics)
	pool.Start(ctx)
}
