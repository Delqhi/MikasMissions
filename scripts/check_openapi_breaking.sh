#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

base_spec="libs/contracts-api/openapi/v1.baseline.yaml"
revision_spec="libs/contracts-api/openapi/v1.yaml"

if [[ ! -f "$base_spec" ]]; then
  echo "[FAIL] missing baseline OpenAPI spec: $base_spec"
  exit 1
fi

if [[ ! -f "$revision_spec" ]]; then
  echo "[FAIL] missing revision OpenAPI spec: $revision_spec"
  exit 1
fi

go run github.com/oasdiff/oasdiff@v1.11.7 breaking "$base_spec" "$revision_spec" --fail-on ERR >/dev/null

echo "[OK] no breaking changes from baseline OpenAPI spec"
