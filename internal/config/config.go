package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	GRPCPort string
	HTTPPort string
	QueueBackend string
	RedisAddr     string        
	RedisPassword string
	RedisDB       int
	RedisTimeout  time.Duration 
	WorkerCount int 
	Env string
}

func Load() *Config {
	return &Config{
		GRPCPort:      envStr("GRPC_PORT", ":50051"),
		HTTPPort:      envStr("HTTP_PORT", ":8080"),
		QueueBackend:  envStr("QUEUE_BACKEND", "redis"),
		RedisAddr:     envStr("REDIS_ADDR", "localhost:6379"),
		RedisPassword: envStr("REDIS_PASSWORD", ""),
		RedisDB:       envInt("REDIS_DB", 0),
		RedisTimeout:  envDuration("REDIS_TIMEOUT", 5*time.Second),
		WorkerCount:   envInt("WORKER_COUNT", 5),
		Env:           envStr("ENV", "development"),
	}
}

func envStr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func envDuration(key string, def time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return def
	}
	return d
}
