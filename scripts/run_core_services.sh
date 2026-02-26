#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."
export GOCACHE=${GOCACHE:-$(pwd)/.cache/go-build}

GOCACHE="$GOCACHE" go run ./apps/identity-service/cmd &
GOCACHE="$GOCACHE" go run ./apps/profile-service/cmd &
GOCACHE="$GOCACHE" go run ./apps/catalog-service/cmd &
GOCACHE="$GOCACHE" go run ./apps/recommendation-service/cmd &
GOCACHE="$GOCACHE" go run ./apps/playback-service/cmd &
GOCACHE="$GOCACHE" go run ./apps/progress-service/cmd &
GOCACHE="$GOCACHE" go run ./apps/billing-service/cmd &
GOCACHE="$GOCACHE" go run ./apps/admin-studio-service/cmd &
GOCACHE="$GOCACHE" go run ./apps/api-gateway-service/cmd &

wait
