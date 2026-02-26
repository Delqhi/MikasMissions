#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

if ! command -v docker >/dev/null 2>&1; then
  echo "docker is required for compose-smoke"
  exit 1
fi

if ! docker info >/dev/null 2>&1; then
  echo "docker daemon is not available"
  exit 1
fi

compose_file="infra/docker-compose.smoke.yml"
compose_project=${COMPOSE_PROJECT_NAME:-infra}
service_build_timeout_seconds=${COMPOSE_SERVICE_BUILD_TIMEOUT_SECONDS:-900}
docker_build_progress=${DOCKER_BUILD_PROGRESS:-plain}
gateway_health_timeout_seconds=${COMPOSE_GATEWAY_HEALTH_TIMEOUT_SECONDS:-180}
gateway_url=${GATEWAY_URL:-http://127.0.0.1:8080}
smoke_script=${SMOKE_SCRIPT:-./scripts/smoke_local.sh}

services=(
  identity-service
  profile-service
  catalog-service
  recommendation-service
  playback-service
  progress-service
  billing-service
  admin-studio-service
  api-gateway-service
)

run_with_timeout() {
  local timeout_seconds=$1
  shift
  "$@" &
  local pid=$!
  local started_at
  started_at=$(date +%s)
  while kill -0 "$pid" >/dev/null 2>&1; do
    local now
    now=$(date +%s)
    if (( now - started_at >= timeout_seconds )); then
      echo "[ERROR] command timed out after ${timeout_seconds}s: $*"
      kill "$pid" >/dev/null 2>&1 || true
      wait "$pid" >/dev/null 2>&1 || true
      return 124
    fi
    sleep 2
  done
  wait "$pid"
}

cleanup() {
  docker compose -p "$compose_project" -f "$compose_file" down -v >/dev/null 2>&1 || true
}
trap cleanup EXIT

echo "[1/4] build images (sequential)"
for service in "${services[@]}"; do
  dockerfile="apps/${service}/Dockerfile"
  image="${compose_project}-${service}:latest"
  echo "  - ${service}"
  run_with_timeout "$service_build_timeout_seconds" docker build --progress="$docker_build_progress" -f "$dockerfile" -t "$image" .
done

echo "[2/4] start stack"
docker compose -p "$compose_project" -f "$compose_file" up -d --no-build

echo "[3/4] wait for gateway health"
for _ in $(seq 1 "$gateway_health_timeout_seconds"); do
  if curl -sS "$gateway_url/healthz" >/dev/null 2>&1; then
    break
  fi
  sleep 1
done

if ! curl -sS "$gateway_url/healthz" >/dev/null 2>&1; then
  echo "[ERROR] gateway did not become healthy on ${gateway_url}"
  docker compose -p "$compose_project" -f "$compose_file" ps || true
  docker compose -p "$compose_project" -f "$compose_file" logs --tail=120 || true
  exit 1
fi

echo "[4/4] run smoke"
GATEWAY_URL="$gateway_url" "$smoke_script"

echo "[OK] compose smoke passed"
