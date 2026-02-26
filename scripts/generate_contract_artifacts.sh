#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

if ! command -v pnpm >/dev/null 2>&1; then
  echo "[FAIL] missing required command: pnpm"
  exit 1
fi

go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1 \
  -generate types \
  -package generated \
  -o libs/contracts-api/generated/openapi_types.gen.go \
  libs/contracts-api/openapi/v1.yaml

(
  cd frontend/web
  pnpm dlx openapi-typescript@7.13.0 ../../libs/contracts-api/openapi/v1.yaml -o lib/generated/api-types.ts >/dev/null
)

echo "[OK] contract artifacts generated"
