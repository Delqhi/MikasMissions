#!/usr/bin/env bash
set -euo pipefail

missing=0

for schema in libs/contracts-events/schemas/*.v1.json; do
  base=$(basename "$schema" .json)
  test_file="libs/contracts-events/${base//./_}_contract_test.go"
  if [[ ! -f "$test_file" ]]; then
    echo "[FAIL] missing contract test for schema: $schema"
    missing=1
  fi
done

if (( missing != 0 )); then
  exit 1
fi

echo "[OK] event schema contract coverage passed"
