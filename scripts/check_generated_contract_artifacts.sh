#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

if ! command -v pnpm >/dev/null 2>&1; then
  echo "[FAIL] missing required command: pnpm"
  exit 1
fi

tmp_dir=$(mktemp -d)
trap 'rm -rf "$tmp_dir"' EXIT

tmp_go="$tmp_dir/openapi_types.gen.go"
tmp_ts="$tmp_dir/api-types.ts"

go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1 \
  -generate types \
  -package generated \
  -o "$tmp_go" \
  libs/contracts-api/openapi/v1.yaml >/dev/null

(
  cd frontend/web
  pnpm dlx openapi-typescript@7.13.0 ../../libs/contracts-api/openapi/v1.yaml -o "$tmp_ts" >/dev/null
)

if ! cmp -s "$tmp_go" libs/contracts-api/generated/openapi_types.gen.go; then
  echo "[FAIL] generated Go contract artifact drift detected"
  echo "run: ./scripts/generate_contract_artifacts.sh"
  exit 1
fi

if ! cmp -s "$tmp_ts" frontend/web/lib/generated/api-types.ts; then
  echo "[FAIL] generated TS contract artifact drift detected"
  echo "run: ./scripts/generate_contract_artifacts.sh"
  exit 1
fi

echo "[OK] generated contract artifacts are up-to-date"
