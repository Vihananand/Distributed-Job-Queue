#!/usr/bin/env bash
# scripts/load_test.sh — load test the gRPC server via hey (HTTP) or ghz (gRPC)
#
# Usage:
#   ./scripts/load_test.sh           # 10 000 jobs, 100 concurrent via grpcurl+ghz
#   TOTAL=50000 CONCURRENCY=200 ./scripts/load_test.sh

set -euo pipefail

GRPC_ADDR="${GRPC_ADDR:-localhost:50051}"
TOTAL="${TOTAL:-10000}"
CONCURRENCY="${CONCURRENCY:-100}"

echo "=== Distributed Job Queue — Load Test ==="
echo "  Server   : $GRPC_ADDR"
echo "  Total    : $TOTAL requests"
echo "  Workers  : $CONCURRENCY concurrent"
echo ""

# Check ghz is available (go install github.com/bojand/ghz/cmd/ghz@latest)
if ! command -v ghz &>/dev/null; then
  echo "[!] ghz not found. Install with:"
  echo "    go install github.com/bojand/ghz/cmd/ghz@latest"
  exit 1
fi

ghz \
  --insecure \
  --proto api/proto/job.proto \
  --call jobpb.JobService.SubmitJob \
  --data '{"type":"echo","payload":"aGVsbG8gd29ybGQ=","priority":0,"delay_seconds":0,"max_retries":3}' \
  --total "$TOTAL" \
  --concurrency "$CONCURRENCY" \
  --timeout 30s \
  "$GRPC_ADDR"
