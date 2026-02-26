#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

soak_hours=${SOAK_HOURS:-1}
soak_interval_minutes=${SOAK_INTERVAL_MINUTES:-60}
skip_external_cluster_policy=${TODAY_SKIP_EXTERNAL_CLUSTER_POLICY:-true}

echo "[1/5] launch preflight"
make launch-preflight

echo "[2/5] launch readiness gate"
if [[ "$skip_external_cluster_policy" == "true" ]]; then
  echo "[INFO] TODAY_SKIP_EXTERNAL_CLUSTER_POLICY=true -> external-cluster policy is bypassed for local readiness only"
  LAUNCH_REQUIRE_EXTERNAL_CLUSTER=false make launch-readiness-gate
else
  make launch-readiness-gate
fi

echo "[3/5] staging deploy dry-run"
make staging-deploy-dry-run

echo "[4/5] staging soak"
SOAK_HOURS="$soak_hours" SOAK_INTERVAL_MINUTES="$soak_interval_minutes" make staging-soak

echo "[5/5] decision packet"
LAUNCH_SKIP_PREFLIGHT=true make launch-decision-packet

echo
echo "[OK] today-ready essentials passed"
