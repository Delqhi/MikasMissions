#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

hours=${SOAK_HOURS:-24}
interval_minutes=${SOAK_INTERVAL_MINUTES:-60}
fail_fast=${SOAK_FAIL_FAST:-true}
evidence_dir=${SOAK_EVIDENCE_DIR:-docs/runbooks/evidence/staging-soak-$(date -u +%Y%m%dT%H%M%SZ)}

if ! [[ "$hours" =~ ^[0-9]+$ ]]; then
  echo "SOAK_HOURS must be an integer"
  exit 1
fi

if ! [[ "$interval_minutes" =~ ^[0-9]+$ ]]; then
  echo "SOAK_INTERVAL_MINUTES must be an integer"
  exit 1
fi

if [[ "$interval_minutes" -le 0 ]]; then
  echo "SOAK_INTERVAL_MINUTES must be greater than 0"
  exit 1
fi

iterations=$((hours * 60 / interval_minutes))
if [[ "$iterations" -le 0 ]]; then
  echo "calculated soak iterations is 0; increase SOAK_HOURS or lower SOAK_INTERVAL_MINUTES"
  exit 1
fi

mkdir -p "$evidence_dir"
summary="$evidence_dir/summary.md"

passed=0
failed=0

run_probe() {
  local name=$1
  local cmd=$2
  local log=$3
  if eval "$cmd" >"$log" 2>&1; then
    echo "PASS" >>"$log.status"
    return 0
  fi
  echo "FAIL" >>"$log.status"
  return 1
}

echo "# Staging Soak Summary" >"$summary"
echo "" >>"$summary"
echo "- start_utc: $(date -u +"%Y-%m-%dT%H:%M:%SZ")" >>"$summary"
echo "- hours: $hours" >>"$summary"
echo "- interval_minutes: $interval_minutes" >>"$summary"
echo "- iterations: $iterations" >>"$summary"
echo "" >>"$summary"
echo "| Iteration | UTC Time | e2e-smoke | e2e-auth-smoke |" >>"$summary"
echo "|---|---|---|---|" >>"$summary"

for i in $(seq 1 "$iterations"); do
  stamp=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
  smoke_log="$evidence_dir/iter-${i}-e2e-smoke.log"
  auth_log="$evidence_dir/iter-${i}-e2e-auth-smoke.log"

  smoke_status="PASS"
  auth_status="PASS"

  if ! run_probe "e2e-smoke" "make e2e-smoke" "$smoke_log"; then
    smoke_status="FAIL"
    failed=$((failed + 1))
  else
    passed=$((passed + 1))
  fi

  if ! run_probe "e2e-auth-smoke" "make e2e-auth-smoke" "$auth_log"; then
    auth_status="FAIL"
    failed=$((failed + 1))
  else
    passed=$((passed + 1))
  fi

  echo "| $i | $stamp | $smoke_status | $auth_status |" >>"$summary"

  if [[ "$failed" -gt 0 && "$fail_fast" == "true" ]]; then
    echo "" >>"$summary"
    echo "- status: FAILED (fail-fast)" >>"$summary"
    echo "- first_failed_iteration: $i" >>"$summary"
    echo "[FAIL] soak failed at iteration $i"
    exit 1
  fi

  if [[ "$i" -lt "$iterations" ]]; then
    sleep "$((interval_minutes * 60))"
  fi
done

echo "" >>"$summary"
echo "- passed_checks: $passed" >>"$summary"
echo "- failed_checks: $failed" >>"$summary"

if [[ "$failed" -gt 0 ]]; then
  echo "- status: FAILED" >>"$summary"
  echo "[FAIL] soak completed with failures"
  exit 1
fi

echo "- status: PASSED" >>"$summary"
echo "[OK] soak completed without failures"
