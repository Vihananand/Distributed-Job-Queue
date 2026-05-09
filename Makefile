PROTO_DIR    := api/proto
PB_OUT       := internal/transport/grpc/pb
PROTOC       := $(LOCALAPPDATA)/protoc/bin/protoc
GOBIN        := $(shell go env GOPATH)/bin

.PHONY: generate build test bench run-server run-server-mem run-worker docker-up docker-down tidy lint load-test

## generate: Re-run protoc code generation
generate:
	$(PROTOC) \
	  --proto_path=$(PROTO_DIR) \
	  --go_out=$(PB_OUT) --go_opt=paths=source_relative \
	  --go-grpc_out=$(PB_OUT) --go-grpc_opt=paths=source_relative \
	  --plugin=protoc-gen-go=$(GOBIN)/protoc-gen-go \
	  --plugin=protoc-gen-go-grpc=$(GOBIN)/protoc-gen-go-grpc \
	  $(PROTO_DIR)/job.proto

## tidy: Sync go.mod / go.sum
tidy:
	go mod tidy

## build: Build both server and worker binaries
build:
	go build -o bin/server ./cmd/server
	go build -o bin/worker ./cmd/worker

## test: Run all unit tests
test:
	go test -race -timeout 30s ./...

## bench: Run queue benchmarks
bench:
	go test -bench=. -benchmem ./internal/queue/...

## run-server: Run the gRPC server locally (Redis backend)
run-server:
	go env -w GOFLAGS="" && \
	  cmd /c "set QUEUE_BACKEND=redis&& set REDIS_ADDR=localhost:6379&& go run ./cmd/server"

## run-server-mem: Run the gRPC server with in-memory queue (no Redis needed)
run-server-mem:
	cmd /c "set QUEUE_BACKEND=memory&& go run ./cmd/server"

## run-worker: Run a worker pool locally
run-worker:
	cmd /c "set QUEUE_BACKEND=redis&& set REDIS_ADDR=localhost:6379&& set WORKER_COUNT=5&& go run ./cmd/worker"

## docker-up: Start full stack via Docker Compose
docker-up:
	docker compose -f deployments/docker/docker-compose.yml up --build -d

## docker-down: Stop Docker Compose stack
docker-down:
	docker compose -f deployments/docker/docker-compose.yml down -v

## load-test: Run gRPC load test (requires ghz: go install github.com/bojand/ghz/cmd/ghz@latest)
GRPC_ADDR  ?= localhost:50051
TOTAL      ?= 10000
CONCURRENCY ?= 100

load-test:
	ghz --insecure \
	  --proto $(PROTO_DIR)/job.proto \
	  --call jobpb.JobService.SubmitJob \
	  --data '{"type":"echo","payload":"aGVsbG8gd29ybGQ=","priority":0,"delay_seconds":0,"max_retries":3}' \
	  --total $(TOTAL) \
	  --concurrency $(CONCURRENCY) \
	  --timeout 30s \
	  $(GRPC_ADDR)

## lint: Run go vet
lint:
	go vet ./...
