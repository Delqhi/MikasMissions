#!/usr/bin/env bash
set -euo pipefail

lock_file="libs/contracts-api/schemas/API_SCHEMAS.lock"

if command -v sha256sum >/dev/null 2>&1; then
  cat libs/contracts-api/schemas/*.json | sha256sum | awk '{print $1}' > "$lock_file"
else
  cat libs/contracts-api/schemas/*.json | shasum -a 256 | awk '{print $1}' > "$lock_file"
fi

echo "updated $lock_file"
