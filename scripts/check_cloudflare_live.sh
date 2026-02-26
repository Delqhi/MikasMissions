#!/usr/bin/env bash
set -euo pipefail

url="${1:-${CLOUDFLARE_WORKERS_DEV_URL:-}}"
if [[ -z "$url" ]]; then
  echo "[FAIL] missing target URL (arg1 or CLOUDFLARE_WORKERS_DEV_URL)"
  exit 1
fi

headers_file="$(mktemp)"
status_code="$(curl -sS -o /dev/null -D "$headers_file" -w "%{http_code}" --max-time 25 "$url")"

echo "[metric] live_url=${url}"
echo "[metric] status_code=${status_code}"

case "$status_code" in
  200|301|302|307|308) ;;
  *)
    echo "[FAIL] unexpected HTTP status from live URL: ${status_code}"
    exit 1
    ;;
esac

assert_header() {
  local name="$1"
  if ! grep -qi "^${name}:" "$headers_file"; then
    echo "[FAIL] missing live response header: ${name}"
    exit 1
  fi
}

assert_header "Content-Security-Policy"
assert_header "Strict-Transport-Security"
assert_header "X-Frame-Options"
assert_header "X-Content-Type-Options"
assert_header "Permissions-Policy"
assert_header "Referrer-Policy"

page_status() {
  local path="$1"
  curl -sS -o /dev/null -w "%{http_code}" --max-time 25 "${url%/}${path}"
}

home_status="$(page_status "/de")"
parents_status="$(page_status "/de/parents")"
echo "[metric] de_home_status=${home_status}"
echo "[metric] de_parents_status=${parents_status}"

if [[ "$home_status" != "200" ]]; then
  echo "[FAIL] /de must return 200 (actual: ${home_status})"
  exit 1
fi

if [[ "$parents_status" != "200" ]]; then
  echo "[FAIL] /de/parents must return 200 (actual: ${parents_status})"
  exit 1
fi

echo "[OK] live Cloudflare smoke passed"
