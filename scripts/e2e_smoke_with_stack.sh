#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

gateway=${GATEWAY_URL:-http://localhost:8080}
gocache=${GOCACHE:-$(pwd)/.cache/go-build}
run_id=$$
health_timeout_seconds=${E2E_HEALTH_TIMEOUT_SECONDS:-60}
health_poll_interval_seconds=${E2E_HEALTH_POLL_INTERVAL_SECONDS:-1}
export AUTH_JWT_SECRET=${AUTH_JWT_SECRET:-local-dev-secret}
smoke_script=${SMOKE_SCRIPT:-./scripts/smoke_local.sh}

SERVICE_PIDS=()
SERVICE_BINS=()
SERVICE_NAMES=()
STARTED_STACK=0

clear_ports() {
  for port in 8080 8081 8082 8083 8084 8085 8086 8089 8090; do
    if lsof -t -iTCP:"$port" -sTCP:LISTEN >/dev/null 2>&1; then
      lsof -t -iTCP:"$port" -sTCP:LISTEN | xargs -r kill >/dev/null 2>&1 || true
    fi
  done
}

cleanup() {
  if [[ "$STARTED_STACK" -eq 1 ]]; then
    for pid in "${SERVICE_PIDS[@]}"; do
      pkill -TERM -P "$pid" >/dev/null 2>&1 || true
      if kill -0 "$pid" >/dev/null 2>&1; then
        kill "$pid" >/dev/null 2>&1 || true
      fi
    done
    sleep 0.2
    clear_ports
    for bin in "${SERVICE_BINS[@]}"; do
      rm -f "$bin"
    done
    wait >/dev/null 2>&1 || true
  fi
}

service_log_path() {
  local name=$1
  echo "/tmp/mikasmissions-${name}.log"
}

dump_service_diagnostics() {
  if [[ "$STARTED_STACK" -ne 1 ]]; then
    return
  fi

  echo
  echo "[diag] service process status"
  if [[ "${#SERVICE_PIDS[@]}" -gt 0 ]]; then
    ps -o pid,ppid,stat,etime,command -p "${SERVICE_PIDS[@]}" 2>/dev/null || true
  else
    echo "[diag] no tracked service pids"
  fi

  echo
  echo "[diag] listeners on ports 8080-8090"
  lsof -nP -iTCP:8080-8090 -sTCP:LISTEN 2>/dev/null || echo "[diag] no listeners found"

  echo
  echo "[diag] tail service logs"
  local i name log
  for i in "${!SERVICE_NAMES[@]}"; do
    name="${SERVICE_NAMES[$i]}"
    log=$(service_log_path "$name")
    echo "----- ${name} (${log}) -----"
    if [[ ! -f "$log" ]]; then
      echo "(missing)"
      continue
    fi
    if [[ ! -s "$log" ]]; then
      echo "(empty)"
      continue
    fi
    tail -n 80 "$log" || true
  done
}

wait_for_gateway() {
  wait_for_url "$gateway/healthz" "gateway"
}

start_service() {
  local name=$1
  local cmd=$2
  local bin=/tmp/mikasmissions-${name}-${run_id}
  local log
  log=$(service_log_path "$name")
  GOCACHE="$gocache" go build -o "$bin" "$cmd"
  : > "$log"
  SERVICE_BINS+=("$bin")
  SERVICE_NAMES+=("$name")
  "$bin" >"$log" 2>&1 &
  SERVICE_PIDS+=($!)
}

wait_for_url() {
  local url=$1
  local label=${2:-$url}
  local attempts=$((health_timeout_seconds / health_poll_interval_seconds))
  if (( attempts < 1 )); then
    attempts=1
  fi
  for _ in $(seq 1 "$attempts"); do
    if curl -fsS "$url" >/dev/null 2>&1; then
      return 0
    fi
    sleep "$health_poll_interval_seconds"
  done
  echo "[ERROR] timed out waiting for ${label} at ${url} after ${health_timeout_seconds}s"
  dump_service_diagnostics
  return 1
}

if ! curl -fsS "$gateway/healthz" >/dev/null 2>&1; then
  STARTED_STACK=1
  trap cleanup EXIT INT TERM
  clear_ports
  start_service identity ./apps/identity-service/cmd
  start_service profile ./apps/profile-service/cmd
  start_service catalog ./apps/catalog-service/cmd
  start_service recommendation ./apps/recommendation-service/cmd
  start_service playback ./apps/playback-service/cmd
  start_service progress ./apps/progress-service/cmd
  start_service billing ./apps/billing-service/cmd
  start_service admin-studio ./apps/admin-studio-service/cmd
  start_service gateway ./apps/api-gateway-service/cmd
  wait_for_url http://127.0.0.1:8081/healthz identity-service
  wait_for_url http://127.0.0.1:8082/healthz profile-service
  wait_for_url http://127.0.0.1:8083/healthz catalog-service
  wait_for_url http://127.0.0.1:8084/healthz recommendation-service
  wait_for_url http://127.0.0.1:8085/healthz playback-service
  wait_for_url http://127.0.0.1:8086/healthz progress-service
  wait_for_url http://127.0.0.1:8089/healthz billing-service
  wait_for_url http://127.0.0.1:8090/healthz admin-studio-service
  wait_for_gateway
fi

"$smoke_script"
