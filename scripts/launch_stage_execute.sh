#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

usage() {
  cat <<'EOF'
Usage:
  ./scripts/launch_stage_execute.sh \
    --stage <1|10|50|100> \
    --read-p95-ms <number> \
    --write-p95-ms <number> \
    --error-5xx-percent <number> \
    [--owner <name>] \
    [--notes <text>]

Behavior:
1. Evaluates stop/go using scripts/launch_go_nogo.sh
2. Writes evidence row into launch rollout report
3. On STOP, triggers rollback by default (set LAUNCH_AUTO_ROLLBACK=false to disable)
EOF
}

is_number() {
  [[ "$1" =~ ^[0-9]+([.][0-9]+)?$ ]]
}

stage=""
read_p95=""
write_p95=""
error_5xx=""
owner=${LAUNCH_OWNER:-unassigned}
notes=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --stage)
      stage=${2:-}
      shift 2
      ;;
    --read-p95-ms)
      read_p95=${2:-}
      shift 2
      ;;
    --write-p95-ms)
      write_p95=${2:-}
      shift 2
      ;;
    --error-5xx-percent)
      error_5xx=${2:-}
      shift 2
      ;;
    --owner)
      owner=${2:-}
      shift 2
      ;;
    --notes)
      notes=${2:-}
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "unknown argument: $1"
      usage
      exit 1
      ;;
  esac
done

if [[ -z "$stage" || -z "$read_p95" || -z "$write_p95" || -z "$error_5xx" ]]; then
  echo "missing required arguments"
  usage
  exit 1
fi

if [[ "$stage" != "1" && "$stage" != "10" && "$stage" != "50" && "$stage" != "100" ]]; then
  echo "invalid stage: $stage (allowed: 1, 10, 50, 100)"
  exit 1
fi

for metric in "$read_p95" "$write_p95" "$error_5xx"; do
  if ! is_number "$metric"; then
    echo "all metric values must be numeric"
    exit 1
  fi
done

rollout_dir=${LAUNCH_ROLLOUT_EVIDENCE_DIR:-docs/runbooks/evidence/launch-rollout-$(date -u +%Y%m%dT%H%M%SZ)}
report_md="$rollout_dir/rollout_report.md"
incident_md="$rollout_dir/incidents.md"
evaluation_log="$rollout_dir/stage-${stage}-evaluation.log"
auto_rollback=${LAUNCH_AUTO_ROLLBACK:-true}
timestamp_cet=$(TZ=Europe/Berlin date +"%Y-%m-%d %H:%M:%S %Z")

mkdir -p "$rollout_dir"

if [[ ! -f "$report_md" ]]; then
  cat >"$report_md" <<'EOF'
# Launch Rollout Report

| Stage | Read p95 (ms) | Write p95 (ms) | 5xx (%) | Result | Time (CET) | Owner | Notes |
|---|---:|---:|---:|---|---|---|---|
EOF
fi

if ./scripts/launch_go_nogo.sh \
  --stage "$stage" \
  --read-p95-ms "$read_p95" \
  --write-p95-ms "$write_p95" \
  --error-5xx-percent "$error_5xx" >"$evaluation_log" 2>&1; then
  result="GO"
  echo "| $stage | $read_p95 | $write_p95 | $error_5xx | $result | $timestamp_cet | $owner | $notes |" >>"$report_md"
  echo "[GO] stage $stage passed; report updated at $report_md"
  exit 0
fi

result="STOP"
echo "| $stage | $read_p95 | $write_p95 | $error_5xx | $result | $timestamp_cet | $owner | $notes |" >>"$report_md"

{
  echo "## Stage ${stage}% STOP - $timestamp_cet"
  echo "- owner: $owner"
  echo "- read_p95_ms: $read_p95"
  echo "- write_p95_ms: $write_p95"
  echo "- error_5xx_percent: $error_5xx"
  echo "- evaluation_log: $evaluation_log"
  echo "- action: rollout halted"
  echo ""
} >>"$incident_md"

if [[ "$auto_rollback" == "true" ]]; then
  if make staging-rollback >>"$evaluation_log" 2>&1; then
    echo "- rollback: executed successfully" >>"$incident_md"
  else
    echo "- rollback: attempted and failed (see $evaluation_log)" >>"$incident_md"
  fi
else
  echo "- rollback: skipped (LAUNCH_AUTO_ROLLBACK=false)" >>"$incident_md"
fi

echo "[STOP] stage $stage failed; report updated at $report_md"
exit 1
