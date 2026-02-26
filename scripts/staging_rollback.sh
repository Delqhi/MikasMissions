#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

namespace=${K8S_NAMESPACE:-mikasmissions-dev}

require_cmd() {
  local cmd=$1
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "missing required command: $cmd"
    exit 1
  fi
}

rollout_undo_wait() {
  local deployment=$1
  kubectl -n "$namespace" rollout undo "deployment/$deployment"
  kubectl -n "$namespace" rollout status "deployment/$deployment" --timeout=180s
}

require_cmd kubectl

./scripts/check_kube_config.sh

echo "[1/2] Rolling back core deployments"
for dep in \
  api-gateway-service \
  identity-service \
  profile-service \
  catalog-service \
  recommendation-service \
  playback-service \
  progress-service \
  admin-studio-service \
  billing-service \
  web-frontend; do
  rollout_undo_wait "$dep"
done

echo "[2/2] Rolling back worker deployments"
for dep in \
  worker-ingest \
  worker-transcode \
  worker-policy \
  worker-publish \
  worker-outbox-relay \
  worker-gen-orchestrator \
  worker-gen-nim \
  worker-gen-qc; do
  rollout_undo_wait "$dep"
done

kubectl -n "$namespace" get deploy

echo
echo "[OK] rollback completed for namespace $namespace"
