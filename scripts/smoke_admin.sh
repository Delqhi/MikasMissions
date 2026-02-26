#!/usr/bin/env bash
set -euo pipefail

gateway=${GATEWAY_URL:-http://localhost:8080}
curl_timeout_connect=${SMOKE_CURL_CONNECT_TIMEOUT_SECONDS:-3}
curl_timeout_max=${SMOKE_CURL_MAX_TIMEOUT_SECONDS:-20}
curl_flags=(-sS --connect-timeout "$curl_timeout_connect" --max-time "$curl_timeout_max")

http() {
  curl "${curl_flags[@]}" "$@"
}

assert_contains() {
  local payload=$1
  local marker=$2
  local label=$3
  if ! echo "$payload" | grep -q "$marker"; then
    echo "$label failed: missing marker $marker"
    echo "$payload"
    exit 1
  fi
}

for _ in $(seq 1 40); do
  if http "$gateway/healthz" >/dev/null 2>&1; then
    break
  fi
  sleep 1
done

admin_login=$(http -X POST "$gateway/v1/admin/login" \
  -H 'content-type: application/json' \
  -d '{"email":"admin@mikasmissions.local","password":"AdminPass123!"}')
admin_token=$(echo "$admin_login" | sed -n 's/.*"access_token":"\([^"]*\)".*/\1/p')
if [[ -z "$admin_token" ]]; then
  echo "admin login failed: $admin_login"
  exit 1
fi
auth_header="Authorization: Bearer $admin_token"

workflows_before=$(http -H "$auth_header" "$gateway/v1/admin/workflows")
assert_contains "$workflows_before" '"workflows"' "list workflows"

model_profile=$(http -X PUT "$gateway/v1/admin/model-profiles/nim-default" \
  -H 'content-type: application/json' \
  -H "$auth_header" \
  -d '{"model_profile_id":"nim-default","provider":"nvidia_nim","base_url":"http://127.0.0.1:9000","model_id":"nim-video-v1","timeout_ms":12000,"max_retries":2,"safety_preset":"kids_strict"}')
assert_contains "$model_profile" '"model_profile_id":"nim-default"' "put model profile"

workflow=$(http -X POST "$gateway/v1/admin/workflows" \
  -H 'content-type: application/json' \
  -H "$auth_header" \
  -d '{"name":"Smoke Workflow","description":"Smoke generator workflow","content_suitability":"core","age_band":"6-11","steps":["prompt","generate","qc"],"model_profile_id":"nim-default","safety_profile":"strict"}')
workflow_id=$(echo "$workflow" | sed -n 's/.*"workflow_id":"\([^"]*\)".*/\1/p')
if [[ -z "$workflow_id" ]]; then
  echo "create workflow failed: $workflow"
  exit 1
fi

run=$(http -X POST "$gateway/v1/admin/workflows/$workflow_id/runs" \
  -H 'content-type: application/json' \
  -H "$auth_header" \
  -d '{"input_payload":{"theme":"space"},"priority":"normal","auto_publish":false}')
run_id=$(echo "$run" | sed -n 's/.*"run_id":"\([^"]*\)".*/\1/p')
if [[ -z "$run_id" ]]; then
  echo "create run failed: $run"
  exit 1
fi

run_status=$(http -H "$auth_header" "$gateway/v1/admin/runs/$run_id")
assert_contains "$run_status" '"run_id"' "get run"

run_logs=$(http -H "$auth_header" "$gateway/v1/admin/runs/$run_id/logs")
assert_contains "$run_logs" '"logs"' "get run logs"

echo "[OK] admin smoke passed"
