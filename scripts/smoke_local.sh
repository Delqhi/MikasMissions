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
  if echo "$payload" | grep -q "upstream unavailable"; then
    echo "$label failed: upstream unavailable"
    echo "$payload"
    exit 1
  fi
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

signup_response=$(http -X POST "$gateway/v1/parents/signup" \
  -H 'content-type: application/json' \
  -d '{"email":"parent@example.com","password":"1234567890x","country":"DE","language":"de","marketing":false,"accepted_terms":true}')

login_response=$(http -X POST "$gateway/v1/parents/login" \
  -H 'content-type: application/json' \
  -d '{"email":"parent@example.com","password":"1234567890x"}')

parent_id=$(echo "$login_response" | sed -n 's/.*"parent_user_id":"\([^"]*\)".*/\1/p')
auth_token=$(echo "$login_response" | sed -n 's/.*"access_token":"\([^"]*\)".*/\1/p')
if [[ -z "$parent_id" || -z "$auth_token" ]]; then
  echo "login failed: $login_response"
  exit 1
fi
auth_header="Authorization: Bearer $auth_token"

http -X POST "$gateway/v1/parents/consent/verify" \
  -H 'content-type: application/json' \
  -d "{\"parent_user_id\":\"$parent_id\",\"method\":\"card\",\"challenge\":\"ok\"}" >/dev/null

profile_response=$(http -X POST "$gateway/v1/children/profiles" \
  -H 'content-type: application/json' \
  -H "$auth_header" \
  -d "{\"parent_user_id\":\"$parent_id\",\"display_name\":\"Mika\",\"age_band\":\"6-11\",\"avatar\":\"robot\"}")

child_id=$(echo "$profile_response" | sed -n 's/.*"child_profile_id":"\([^"]*\)".*/\1/p')
if [[ -z "$child_id" ]]; then
  echo "profile failed: $profile_response"
  exit 1
fi

challenge_response=$(http -X POST "$gateway/v1/parents/gates/challenge" \
  -H 'content-type: application/json' \
  -H "$auth_header" \
  -d "{\"parent_user_id\":\"$parent_id\",\"child_profile_id\":\"$child_id\",\"method\":\"pin\"}")
challenge_id=$(echo "$challenge_response" | sed -n 's/.*"challenge_id":"\([^"]*\)".*/\1/p')
if [[ -z "$challenge_id" ]]; then
  echo "challenge failed: $challenge_response"
  exit 1
fi

verify_response=$(http -X POST "$gateway/v1/parents/gates/verify" \
  -H 'content-type: application/json' \
  -H "$auth_header" \
  -d "{\"parent_user_id\":\"$parent_id\",\"child_profile_id\":\"$child_id\",\"challenge_id\":\"$challenge_id\",\"response\":\"ok\"}")
gate_token=$(echo "$verify_response" | sed -n 's/.*"gate_token":"\([^"]*\)".*/\1/p')
if [[ -z "$gate_token" ]]; then
  echo "gate verify failed: $verify_response"
  exit 1
fi

home_rails_response=$(http -H "$auth_header" "$gateway/v1/home/rails?child_profile_id=$child_id")
assert_contains "$home_rails_response" '"rails"' "home rails"
echo "$home_rails_response"

kids_home_response=$(http -H "$auth_header" "$gateway/v1/kids/home?child_profile_id=$child_id&mode=core")
assert_contains "$kids_home_response" '"mode"' "kids home"
echo "$kids_home_response"

catalog_response=$(http -H "$auth_header" "$gateway/v1/catalog/episodes/ep-demo-1")
assert_contains "$catalog_response" '"episode_id"' "catalog episode"
echo "$catalog_response"

entitlements_response=$(http -H "$auth_header" "$gateway/v1/billing/entitlements?parent_user_id=$parent_id")
assert_contains "$entitlements_response" '"active":true' "billing"
if ! echo "$entitlements_response" | grep -q '"active":true'; then
  echo "billing failed: $entitlements_response"
  exit 1
fi
playback_response=$(http -X POST "$gateway/v1/playback/sessions" \
  -H 'content-type: application/json' \
  -H "$auth_header" \
  -d "{\"child_profile_id\":\"$child_id\",\"episode_id\":\"ep-demo-1\",\"device_type\":\"external-link\",\"parent_gate_token\":\"$gate_token\",\"session_limit_minutes\":45,\"session_minutes_used\":10,\"entitlement_status\":\"active\",\"autoplay_requested\":false}")
assert_contains "$playback_response" '"playback_session_id"' "playback session"
echo "$playback_response"

watch_event_response=$(http -X POST "$gateway/v1/progress/watch-events" \
  -H 'content-type: application/json' \
  -H "$auth_header" \
  -d "{\"child_profile_id\":\"$child_id\",\"episode_id\":\"ep-demo-1\",\"watch_ms\":1000,\"event_time\":\"2026-03-02T12:00:00Z\"}")
assert_contains "$watch_event_response" '"accepted":true' "watch event"
echo "$watch_event_response"

kids_progress_response=$(http -H "$auth_header" "$gateway/v1/kids/progress/$child_id")
assert_contains "$kids_progress_response" '"child_profile_id"' "kids progress"
echo "$kids_progress_response"

echo "\n[OK] smoke flow passed"
