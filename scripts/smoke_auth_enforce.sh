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

signup_parent() {
  local email=$1
  http -X POST "$gateway/v1/parents/signup" \
    -H 'content-type: application/json' \
    -d "{\"email\":\"$email\",\"password\":\"1234567890x\",\"country\":\"DE\",\"language\":\"de\",\"marketing\":false,\"accepted_terms\":true}" >/dev/null
  local login
  login=$(http -X POST "$gateway/v1/parents/login" \
    -H 'content-type: application/json' \
    -d "{\"email\":\"$email\",\"password\":\"1234567890x\"}")
  local parent_id
  local token
  parent_id=$(echo "$login" | sed -n 's/.*"parent_user_id":"\([^"]*\)".*/\1/p')
  token=$(echo "$login" | sed -n 's/.*"access_token":"\([^"]*\)".*/\1/p')
  if [[ -z "$parent_id" || -z "$token" ]]; then
    echo "login failed for $email: $login"
    exit 1
  fi
  http -X POST "$gateway/v1/parents/consent/verify" \
    -H 'content-type: application/json' \
    -d "{\"parent_user_id\":\"$parent_id\",\"method\":\"card\",\"challenge\":\"ok\"}" >/dev/null
  echo "$parent_id|$token"
}

parent_one=$(signup_parent "parent-one@example.com")
parent_one_id=${parent_one%%|*}
parent_one_token=${parent_one##*|}
parent_one_auth="Authorization: Bearer $parent_one_token"

profile_response=$(http -X POST "$gateway/v1/children/profiles" \
  -H 'content-type: application/json' \
  -H "$parent_one_auth" \
  -d "{\"parent_user_id\":\"$parent_one_id\",\"display_name\":\"Mika\",\"age_band\":\"6-11\",\"avatar\":\"robot\"}")
child_id=$(echo "$profile_response" | sed -n 's/.*"child_profile_id":"\([^"]*\)".*/\1/p')
if [[ -z "$child_id" ]]; then
  echo "profile failed: $profile_response"
  exit 1
fi

challenge_response=$(http -X POST "$gateway/v1/parents/gates/challenge" \
  -H 'content-type: application/json' \
  -H "$parent_one_auth" \
  -d "{\"parent_user_id\":\"$parent_one_id\",\"child_profile_id\":\"$child_id\",\"method\":\"pin\"}")
challenge_id=$(echo "$challenge_response" | sed -n 's/.*"challenge_id":"\([^"]*\)".*/\1/p')
verify_response=$(http -X POST "$gateway/v1/parents/gates/verify" \
  -H 'content-type: application/json' \
  -H "$parent_one_auth" \
  -d "{\"parent_user_id\":\"$parent_one_id\",\"child_profile_id\":\"$child_id\",\"challenge_id\":\"$challenge_id\",\"response\":\"ok\"}")
gate_token=$(echo "$verify_response" | sed -n 's/.*"gate_token":"\([^"]*\)".*/\1/p')
if [[ -z "$gate_token" ]]; then
  echo "gate verify failed: $verify_response"
  exit 1
fi

playback_ok=$(http -X POST "$gateway/v1/playback/sessions" \
  -H 'content-type: application/json' \
  -H "$parent_one_auth" \
  -d "{\"child_profile_id\":\"$child_id\",\"episode_id\":\"ep-demo-1\",\"device_type\":\"external-link\",\"parent_gate_token\":\"$gate_token\",\"session_limit_minutes\":45,\"session_minutes_used\":10,\"autoplay_requested\":false}")
assert_contains "$playback_ok" '"playback_session_id"' "playback first gate use"

playback_reuse=$(http -X POST "$gateway/v1/playback/sessions" \
  -H 'content-type: application/json' \
  -H "$parent_one_auth" \
  -d "{\"child_profile_id\":\"$child_id\",\"episode_id\":\"ep-demo-1\",\"device_type\":\"external-link\",\"parent_gate_token\":\"$gate_token\",\"session_limit_minutes\":45,\"session_minutes_used\":10,\"autoplay_requested\":false}")
assert_contains "$playback_reuse" '"parent_gate_required"' "playback gate token reuse"

parent_two=$(signup_parent "parent-two@example.com")
parent_two_token=${parent_two##*|}
parent_two_auth="Authorization: Bearer $parent_two_token"

foreign_progress=$(http -X POST "$gateway/v1/progress/watch-events" \
  -H 'content-type: application/json' \
  -H "$parent_two_auth" \
  -d "{\"child_profile_id\":\"$child_id\",\"episode_id\":\"ep-demo-1\",\"watch_ms\":1000,\"event_time\":\"2026-03-02T12:00:00Z\"}")
assert_contains "$foreign_progress" '"child_profile_forbidden"' "cross-parent protection"

echo "[OK] auth enforce smoke passed"
