#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage:
  ./scripts/launch_go_nogo.sh \
    --stage <1|10|50|100> \
    --read-p95-ms <number> \
    --write-p95-ms <number> \
    --error-5xx-percent <number>

Thresholds (launch v1 defaults):
  read p95 < 400ms
  write p95 < 700ms
  5xx error rate < 1.0%
EOF
}

is_number() {
  [[ "$1" =~ ^[0-9]+([.][0-9]+)?$ ]]
}

stage=""
read_p95=""
write_p95=""
error_5xx=""

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

max_read_ms=${MAX_READ_P95_MS:-400}
max_write_ms=${MAX_WRITE_P95_MS:-700}
max_error_pct=${MAX_5XX_PERCENT:-1.0}

fail=0

awk "BEGIN { exit !($read_p95 < $max_read_ms) }" || {
  echo "[STOP] stage ${stage}%: read p95 ${read_p95}ms >= ${max_read_ms}ms"
  fail=1
}

awk "BEGIN { exit !($write_p95 < $max_write_ms) }" || {
  echo "[STOP] stage ${stage}%: write p95 ${write_p95}ms >= ${max_write_ms}ms"
  fail=1
}

awk "BEGIN { exit !($error_5xx < $max_error_pct) }" || {
  echo "[STOP] stage ${stage}%: 5xx ${error_5xx}% >= ${max_error_pct}%"
  fail=1
}

if [[ "$fail" -eq 1 ]]; then
  exit 1
fi

echo "[GO] stage ${stage}%: read p95=${read_p95}ms write p95=${write_p95}ms 5xx=${error_5xx}%"
