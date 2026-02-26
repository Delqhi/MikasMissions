#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

cp libs/contracts-api/openapi/v1.yaml libs/contracts-api/openapi/v1.baseline.yaml
echo "[OK] updated OpenAPI baseline"
