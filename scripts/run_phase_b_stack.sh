#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."
export GOCACHE=${GOCACHE:-$(pwd)/.cache/go-build}

nats_pid=""
if docker info >/dev/null 2>&1; then
  docker compose -f infra/docker-compose.local.yml up -d nats
elif command -v nats-server >/dev/null 2>&1; then
  nats-server -js -p 4222 >/tmp/mm_nats_server.log 2>&1 &
  nats_pid=$!
else
  echo "NATS startup unavailable: run Docker daemon or install nats-server"
  exit 1
fi

export NATS_URL=${NATS_URL:-nats://127.0.0.1:4222}

ready=0
for _ in $(seq 1 90); do
  if (echo > /dev/tcp/127.0.0.1/4222) >/dev/null 2>&1; then
    ready=1
    break
  fi
  sleep 1
done
if [[ "$ready" -ne 1 ]]; then
  echo "nats did not become ready on :4222"
  exit 1
fi

GOCACHE="$GOCACHE" go run ./apps/creator-studio-service/cmd &
pids=($!)
GOCACHE="$GOCACHE" go run ./workers/worker-ingest/cmd &
pids+=($!)
GOCACHE="$GOCACHE" go run ./workers/worker-transcode/cmd &
pids+=($!)
GOCACHE="$GOCACHE" go run ./workers/worker-policy/cmd &
pids+=($!)
GOCACHE="$GOCACHE" go run ./workers/worker-publish/cmd &
pids+=($!)
if [[ -n "${DATABASE_URL:-}" ]]; then
  GOCACHE="$GOCACHE" go run ./workers/worker-outbox-relay/cmd &
  pids+=($!)
fi

cleanup() {
  for pid in "${pids[@]}"; do
    kill "$pid" >/dev/null 2>&1 || true
  done
  if [[ -n "${nats_pid}" ]]; then
    kill "${nats_pid}" >/dev/null 2>&1 || true
  fi
}
trap cleanup EXIT

wait
