#!/usr/bin/env bash
set -euo pipefail

lock_file="libs/contracts-api/schemas/API_SCHEMAS.lock"

if [[ ! -f "$lock_file" ]]; then
  echo "[FAIL] missing schema lock file: $lock_file"
  exit 1
fi

if command -v sha256sum >/dev/null 2>&1; then
  current=$(cat libs/contracts-api/schemas/*.json | sha256sum | awk '{print $1}')
else
  current=$(cat libs/contracts-api/schemas/*.json | shasum -a 256 | awk '{print $1}')
fi
expected=$(cat "$lock_file")

if [[ "$current" != "$expected" ]]; then
  echo "[FAIL] API schema lock mismatch"
  echo "expected: $expected"
  echo "current : $current"
  echo "update with: ./scripts/update_api_schema_lock.sh"
  exit 1
fi

echo "[OK] API schema lock matches"
