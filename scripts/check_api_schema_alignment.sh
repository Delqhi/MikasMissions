#!/usr/bin/env bash
set -euo pipefail

missing=0

assert_pattern() {
  local pattern=$1
  local file=$2
  if ! grep -Eq "$pattern" "$file"; then
    echo "[FAIL] missing pattern '$pattern' in $file"
    missing=1
  fi
}

schemas=(
  "libs/contracts-api/schemas/kids.home.v1.response.json"
  "libs/contracts-api/schemas/kids.progress.v1.response.json"
  "libs/contracts-api/schemas/parents.controls.v1.response.json"
  "libs/contracts-api/schemas/playback.sessions.v1.request.json"
)

for schema in "${schemas[@]}"; do
  if [[ ! -f "$schema" ]]; then
    echo "[FAIL] missing canonical API schema: $schema"
    missing=1
  fi
done

assert_pattern 'json:"summary"' "libs/contracts-api/home_rails.go"
assert_pattern 'json:"thumbnail_url"' "libs/contracts-api/home_rails.go"
assert_pattern 'json:"watched_minutes_today"' "libs/contracts-api/kids_progress.go"
assert_pattern 'json:"watched_minutes_7d"' "libs/contracts-api/kids_progress.go"
assert_pattern 'json:"controls"' "libs/contracts-api/safety_controls.go"
assert_pattern 'json:"parent_gate_token"' "libs/contracts-api/playback_sessions.go"
assert_pattern 'json:"entitlement_status"' "libs/contracts-api/playback_sessions.go"

assert_pattern 'summary: string;' "frontend/web/lib/experience_types.ts"
assert_pattern 'thumbnail_url: string;' "frontend/web/lib/experience_types.ts"
assert_pattern 'watched_minutes_7d: number;' "frontend/web/lib/experience_types.ts"
assert_pattern 'controls: ParentalControls;' "frontend/web/lib/experience_types.ts"
assert_pattern 'primary_actions: string\[\];' "frontend/web/lib/experience_types.ts"

if (( missing != 0 )); then
  exit 1
fi

echo "[OK] API schema alignment checks passed"
