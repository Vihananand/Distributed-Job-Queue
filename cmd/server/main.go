package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/vihan/distributed-job-queue/internal/config"
	"github.com/vihan/distributed-job-queue/internal/producer"
	inqueue "github.com/vihan/distributed-job-queue/internal/queue"
	"github.com/vihan/distributed-job-queue/internal/storage"
	grpcsrv "github.com/vihan/distributed-job-queue/internal/transport/grpc"
	pb "github.com/vihan/distributed-job-queue/internal/transport/grpc/pb"
	"github.com/vihan/distributed-job-queue/internal/utils"
)

func main() {
	cfg := config.Load()
	utils.InitLogger(cfg.Env)

	metrics := utils.NewMetrics()

	// Queue + Storage
	var (
		q     inqueue.Queue
		store storage.Storage
	)
	switch cfg.QueueBackend {
	case "memory":
		slog.Info("using in-memory queue and store")
		q = inqueue.NewMemoryQueue()
		store = storage.NewMemoryStore()
	default:
		slog.Info("using Redis queue and store", "addr", cfg.RedisAddr)
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
		store = storage.NewRedisStore(rdb)
	}

	// Producer + gRPC server
	prod := producer.New(q, store, metrics)
	srv := grpcsrv.New(prod, store, q)

	grpcServer := grpc.NewServer()
	pb.RegisterJobServiceServer(grpcServer, srv)
	reflection.Register(grpcServer) // enables grpcurl introspection

	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		slog.Error("failed to listen", "port", cfg.GRPCPort, "err", err)
		os.Exit(1)
	}

	//Prometheus HTTP server
	httpMux := http.NewServeMux()
	httpMux.Handle("/metrics", promhttp.HandlerFor(metrics.Registry, promhttp.HandlerOpts{}))
	httpMux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, `{"status":"ok"}`)
	})
	httpServer := &http.Server{Addr: cfg.HTTPPort, Handler: httpMux}

	go func() {
		slog.Info("Prometheus metrics", "addr", cfg.HTTPPort+"/metrics")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("http server error", "err", err)
		}
	}()

	//Serve
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("gRPC server listening", "addr", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("gRPC serve error", "err", err)
		}
	}()

	<-ctx.Done()
	slog.Info("shutdown signal received")
	grpcServer.GracefulStop()
	_ = httpServer.Shutdown(context.Background())
	slog.Info("server stopped")
}
