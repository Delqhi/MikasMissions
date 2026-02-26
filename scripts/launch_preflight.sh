#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

api_base_url=${NEXT_PUBLIC_API_BASE_URL:-http://127.0.0.1:8080}

step() {
  printf '\n[%s] %s\n' "$1" "$2"
}

require_cmd() {
  local cmd=$1
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "missing required command: $cmd"
    exit 1
  fi
}

check_runtime_defaults() {
  local config=infra/kubernetes/base/configmap-platform.yaml
  local web=infra/kubernetes/base/services/web-frontend.yaml

  grep -q 'GO_ENV: production' "$config" || {
    echo "runtime config missing GO_ENV=production"
    exit 1
  }
  grep -q 'PERSISTENCE_MODE: strict' "$config" || {
    echo "runtime config missing PERSISTENCE_MODE=strict"
    exit 1
  }
  grep -q 'AUTH_MODE: enforce' "$config" || {
    echo "runtime config missing AUTH_MODE=enforce"
    exit 1
  }
  grep -q 'name: NEXT_PUBLIC_USE_API_FALLBACKS' "$web" || {
    echo "web deployment missing NEXT_PUBLIC_USE_API_FALLBACKS env var"
    exit 1
  }
  grep -q 'value: "false"' "$web" || {
    echo "web deployment must set NEXT_PUBLIC_USE_API_FALLBACKS=false"
    exit 1
  }
}

require_cmd go
require_cmd make
require_cmd npm
require_cmd curl

step "1/11" "Guardrails"
make guard

step "2/11" "Unit and integration tests"
make test

step "3/11" "Build all Go binaries"
make build

step "4/11" "Contract checks"
make contract-check

step "5/11" "End-to-end smoke"
make e2e-smoke

step "6/11" "Auth enforce smoke"
make e2e-auth-smoke

step "7/11" "Admin smoke"
make e2e-admin-smoke

step "8/11" "Generator worker smoke"
make e2e-generator-smoke

step "9/11" "A11y smoke"
make a11y-smoke

step "10/11" "Web production build (fail-closed)"
(
  cd frontend/web
  if [[ ! -d node_modules ]]; then
    npm ci
  fi
  NEXT_PUBLIC_API_BASE_URL="$api_base_url" NEXT_PUBLIC_USE_API_FALLBACKS=false npm run build
)

step "11/11" "Runtime defaults static verification"
check_runtime_defaults

echo
echo "[OK] launch preflight passed"
