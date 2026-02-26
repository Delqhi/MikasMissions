#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."
export GOCACHE=${GOCACHE:-$(pwd)/.cache/go-build}

export NATS_URL=${NATS_URL:-nats://127.0.0.1:4222}
creator_url=${CREATOR_URL:-http://127.0.0.1:8087}

probe_output=$(mktemp)

for _ in $(seq 1 40); do
  if (echo > /dev/tcp/127.0.0.1/4222) >/dev/null 2>&1; then
    break
  fi
  sleep 1
done

GOCACHE="$GOCACHE" go run ./tools/event-probe -nats-url "$NATS_URL" -topic "episode.published.v1" -timeout 25s > "$probe_output" 2>/tmp/mm_probe.err &
probe_pid=$!

for _ in $(seq 1 40); do
  if curl -sS "$creator_url/healthz" >/dev/null 2>&1; then
    break
  fi
  sleep 1
done

curl -sS -X POST "$creator_url/v1/creator/assets/upload" \
  -H 'content-type: application/json' \
  -d '{"source_url":"https://cdn.local/new_episode.mp4","uploader_id":"creator-1"}' >/tmp/mm_upload.out

if ! wait "$probe_pid"; then
  cat /tmp/mm_probe.err
  echo "upload response:"
  cat /tmp/mm_upload.out
  exit 1
fi

echo "episode published event:"
cat "$probe_output"
rm -f "$probe_output"
echo "[OK] phase-b smoke passed"
