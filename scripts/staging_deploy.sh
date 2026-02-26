#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

namespace=${K8S_NAMESPACE:-mikasmissions-dev}
kustomize_path=${KUSTOMIZE_PATH:-infra/kubernetes/base}
dry_run=${DRY_RUN:-false}

required_env() {
  local var_name=$1
  if [[ -z "${!var_name:-}" ]]; then
    echo "missing required env var: $var_name"
    exit 1
  fi
}

require_cmd() {
  local cmd=$1
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "missing required command: $cmd"
    exit 1
  fi
}

is_true() {
  case "$1" in
    1|true|TRUE|yes|YES) return 0 ;;
    *) return 1 ;;
  esac
}

rollout_wait() {
  local deployment=$1
  kubectl -n "$namespace" rollout status "deployment/$deployment" --timeout=180s
}

require_cmd kubectl

if is_true "$dry_run"; then
  database_url=${DATABASE_URL:-postgres://placeholder:placeholder@localhost:5432/mikasmissions?sslmode=disable}
  auth_jwt_secret=${AUTH_JWT_SECRET:-placeholder-dev-secret}
  kubeconfig_tmp=$(mktemp)
  cat >"$kubeconfig_tmp" <<'EOF'
apiVersion: v1
kind: Config
clusters: []
contexts: []
current-context: ""
users: []
EOF
  cleanup_tmp() {
    rm -f "$kubeconfig_tmp"
  }
  trap cleanup_tmp EXIT
else
  ./scripts/check_kube_config.sh
  required_env DATABASE_URL
  required_env AUTH_JWT_SECRET
  database_url=${DATABASE_URL}
  auth_jwt_secret=${AUTH_JWT_SECRET}
fi

echo "[1/6] Applying namespace and runtime config"
if is_true "$dry_run"; then
  test -f infra/kubernetes/base/namespace.yaml
  test -f infra/kubernetes/base/configmap-platform.yaml
else
  kubectl apply -f infra/kubernetes/base/namespace.yaml
  kubectl -n "$namespace" apply -f infra/kubernetes/base/configmap-platform.yaml
fi

echo "[2/6] Upserting platform secrets"
if is_true "$dry_run"; then
  kubectl --kubeconfig="$kubeconfig_tmp" -n "$namespace" create secret generic platform-secrets \
    --from-literal=database-url="$database_url" \
    --from-literal=auth-jwt-secret="$auth_jwt_secret" \
    --from-literal=auth-jwks-url="${AUTH_JWKS_URL:-}" \
    --from-literal=admin-bootstrap-email="${ADMIN_BOOTSTRAP_EMAIL:-}" \
    --from-literal=admin-bootstrap-password="${ADMIN_BOOTSTRAP_PASSWORD:-}" \
    --dry-run=client -o yaml >/dev/null
else
  kubectl -n "$namespace" create secret generic platform-secrets \
    --from-literal=database-url="$database_url" \
    --from-literal=auth-jwt-secret="$auth_jwt_secret" \
    --from-literal=auth-jwks-url="${AUTH_JWKS_URL:-}" \
    --from-literal=admin-bootstrap-email="${ADMIN_BOOTSTRAP_EMAIL:-}" \
    --from-literal=admin-bootstrap-password="${ADMIN_BOOTSTRAP_PASSWORD:-}" \
    --dry-run=client -o yaml | kubectl apply -f -
fi

echo "[3/6] Applying kustomize base"
if is_true "$dry_run"; then
  kubectl kustomize "$kustomize_path" >/dev/null
else
  kubectl -n "$namespace" apply -k "$kustomize_path"
fi

if is_true "$dry_run"; then
  echo "[4/6] Skipping rollout waits in dry-run mode"
  echo "[5/6] Skipping worker waits in dry-run mode"
  echo "[6/6] Dry-run completed (manifests and secrets validate locally)"
  echo
  echo "[OK] staging deploy dry-run passed for namespace $namespace"
  exit 0
fi

echo "[4/6] Waiting for core service rollouts"
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
  rollout_wait "$dep"
done

echo "[5/6] Waiting for worker rollouts"
for dep in \
  worker-ingest \
  worker-transcode \
  worker-policy \
  worker-publish \
  worker-outbox-relay \
  worker-gen-orchestrator \
  worker-gen-nim \
  worker-gen-qc; do
  rollout_wait "$dep"
done

echo "[6/6] Cluster readiness snapshot"
kubectl -n "$namespace" get deploy
kubectl -n "$namespace" get pods
kubectl -n "$namespace" get ingress

echo
echo "[OK] staging deploy completed for namespace $namespace"
