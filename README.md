# Distributed Job Queue

A production-grade distributed job queue written in Go. Supports **gRPC ingestion**, **Redis persistence**, **worker pools**, **priority scheduling**, **delayed jobs**, **exponential-backoff retries**, **dead-letter queues**, and **Prometheus observability**.

---

## Architecture

```
Producer (gRPC client)
        │
        ▼ SubmitJob RPC
┌──────────────────────┐        ┌─────────────┐
│   cmd/server         │◄──────►│    Redis    │
│   (gRPC + /metrics)  │        │  (queue +   │
└──────────┬───────────┘        │   storage)  │
           │ queue.Pop          └─────────────┘
           ▼
┌──────────────────────┐
│   cmd/worker         │
│   (Pool of N goroutines, each running a handler)
└──────────────────────┘
```

---

## Quick Start

### Prerequisites
- Go 1.22+
- Docker (Redis)

### 1. Start Redis

```bash
docker run -d -p 6379:6379 redis:7-alpine
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Run the server

```bash
make run-server
# gRPC listening on :50051
# Prometheus metrics on :8080/metrics
```

### 4. Run a worker pool

```bash
make run-worker
# 5 concurrent workers polling Redis
```

### 5. Submit a job (grpcurl)

```bash
# Install: go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
grpcurl -plaintext -d '{
  "type": "echo",
  "payload": "aGVsbG8gd29ybGQ=",
  "priority": 5,
  "max_retries": 3
}' localhost:50051 jobpb.JobService/SubmitJob
```

### 6. Check health

```bash
grpcurl -plaintext localhost:50051 jobpb.JobService/HealthCheck
```

---

## Configuration (Environment Variables)

| Variable | Default | Description |
|---|---|---|
| `QUEUE_BACKEND` | `redis` | `redis` or `memory` |
| `REDIS_ADDR` | `localhost:6379` | Redis address |
| `REDIS_PASSWORD` | `` | Redis password |
| `REDIS_DB` | `0` | Redis database index |
| `GRPC_PORT` | `:50051` | gRPC listen address |
| `HTTP_PORT` | `:8080` | Prometheus listen address |
| `WORKER_COUNT` | `5` | Number of concurrent workers |
| `ENV` | `development` | `development` (text logs) or `production` (JSON) |

---

## Docker Compose (full stack)

```bash
make docker-up    # starts Redis + server + 2 worker replicas
make docker-down  # stops and removes volumes
```

Scale workers: `docker compose -f deployments/docker/docker-compose.yml up --scale worker=5`

---

## Development

```bash
make test    # unit tests (race detector on)
make bench   # queue push/pop benchmark
make lint    # go vet
make build   # compile both binaries to ./bin/
```

### Re-generate proto (after editing api/proto/job.proto)

```bash
make generate
```

---

## Project Layout

```
├── api/proto/          # gRPC .proto definition
├── cmd/
│   ├── server/         # gRPC + Prometheus server entry-point
│   └── worker/         # Worker pool entry-point
├── internal/
│   ├── models/         # Job struct + Status enum
│   ├── queue/          # Queue interface, in-memory impl, Redis impl
│   ├── storage/        # Storage interface, in-memory impl, Redis impl
│   ├── worker/         # Worker, Pool, HandlerFunc registry
│   ├── producer/       # Validation + job submission logic
│   ├── transport/grpc/ # gRPC server implementation
│   ├── config/         # Env-based config
│   └── utils/          # slog logger + Prometheus metrics
├── pkg/client/         # Public Go SDK (gRPC wrapper)
├── deployments/docker/ # Dockerfile + docker-compose.yml
└── scripts/            # load_test.sh (ghz)
```

---

## Observability

- **Structured logs** — `log/slog` (text in dev, JSON in prod)
- **Prometheus metrics** at `GET :8080/metrics`
  - `djq_jobs_pushed_total`
  - `djq_jobs_processed_total`
  - `djq_jobs_failed_total`
  - `djq_queue_length`

---

## Load Testing

```bash
# Install ghz
go install github.com/bojand/ghz/cmd/ghz@latest

# Run 10 000 jobs at 100 concurrency
make load-test

# Custom
TOTAL=50000 CONCURRENCY=200 bash scripts/load_test.sh
```
