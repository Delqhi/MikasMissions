#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

packet_dir=${LAUNCH_PACKET_DIR:-docs/runbooks/evidence/launch-$(date -u +%Y%m%dT%H%M%SZ)}
skip_preflight=${LAUNCH_SKIP_PREFLIGHT:-false}
preflight_timeout_seconds=${LAUNCH_PREFLIGHT_TIMEOUT_SECONDS:-1800}
mkdir -p "$packet_dir"

preflight_log="$packet_dir/launch_preflight.log"
packet_md="$packet_dir/decision_packet.md"

if [[ "$skip_preflight" == "true" ]]; then
  preflight_status="SKIPPED"
  echo "[WARN] preflight skipped via LAUNCH_SKIP_PREFLIGHT=true" >"$preflight_log"
else
  if [[ "$preflight_timeout_seconds" =~ ^[0-9]+$ ]] && [[ "$preflight_timeout_seconds" -gt 0 ]]; then
    if timeout "$preflight_timeout_seconds" make launch-preflight >"$preflight_log" 2>&1; then
      preflight_status="PASSED"
    else
      preflight_status="FAILED"
    fi
  else
    if make launch-preflight >"$preflight_log" 2>&1; then
      preflight_status="PASSED"
    else
      preflight_status="FAILED"
    fi
  fi
fi

cat >"$packet_md" <<EOF
# Launch Decision Packet

- generated_utc: $(date -u +"%Y-%m-%dT%H:%M:%SZ")
- git_head: $(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
- preflight_status: $preflight_status
- preflight_log: $preflight_log

## Mandatory Evidence

1. Attach output from:
   1. \`make e2e-smoke\`
   2. \`make e2e-auth-smoke\`
   3. \`make e2e-admin-smoke\`
   4. \`make e2e-generator-smoke\`
   5. \`make a11y-smoke\`
2. Attach 24h soak summary:
   1. \`docs/runbooks/evidence/staging-soak-*/summary.md\`
3. Attach rollout decisions for stages:
   1. 1%
   2. 10%
   3. 50%
   4. 100%

## Rollout Metrics Table

| Stage | Read p95 (ms) | Write p95 (ms) | 5xx (%) | launch-go-nogo result | Time (CET) | Owner |
|---|---:|---:|---:|---|---|---|
| 1 |  |  |  |  |  |  |
| 10 |  |  |  |  |  |  |
| 50 |  |  |  |  |  |  |
| 100 |  |  |  |  |  |  |

## Stop Criteria

1. \`5xx < 1.0%\`
2. \`p95 read < 400ms\`
3. \`p95 write < 700ms\`
4. On violation: immediate \`make staging-rollback\` and incident start.

## Final Decision

- go_no_go: GO | NO-GO
- decided_by:
- decided_at_cet:
- notes:
EOF

if [[ "$preflight_status" == "FAILED" ]]; then
  echo "[FAIL] launch preflight failed; see $preflight_log"
  exit 1
fi

echo "[OK] decision packet created: $packet_md"
