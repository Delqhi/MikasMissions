#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

GOCACHE=${GOCACHE:-$(pwd)/.cache/go-build}
export GOCACHE

go test ./workers/worker-gen-orchestrator/... ./workers/worker-gen-nim/... ./workers/worker-gen-qc/... ./libs/generatorprovider/...

echo "[OK] generator worker smoke passed"
