#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

next_config="frontend/web/next.config.mjs"
middleware_file="frontend/web/middleware.ts"
failures=0

assert_contains() {
  local file=$1
  local snippet=$2
  if ! grep -qF "$snippet" "$file"; then
    echo "[FAIL] missing '$snippet' in $file"
    failures=$((failures + 1))
  fi
}

for file in "$next_config" "$middleware_file"; do
  if [[ ! -f "$file" ]]; then
    echo "[FAIL] missing required file: $file"
    failures=$((failures + 1))
  fi
done

assert_contains "$next_config" "Content-Security-Policy"
assert_contains "$next_config" "Cross-Origin-Opener-Policy"
assert_contains "$next_config" "Cross-Origin-Resource-Policy"
assert_contains "$next_config" "Permissions-Policy"
assert_contains "$next_config" "Referrer-Policy"
assert_contains "$next_config" "Strict-Transport-Security"
assert_contains "$next_config" "X-Content-Type-Options"
assert_contains "$next_config" "X-DNS-Prefetch-Control"
assert_contains "$next_config" "X-Frame-Options"
assert_contains "$next_config" "X-Permitted-Cross-Domain-Policies"
assert_contains "$next_config" "default-src 'self'"
assert_contains "$next_config" "frame-ancestors 'none'"
assert_contains "$next_config" "object-src 'none'"
assert_contains "$next_config" "max-age=63072000; includeSubDomains; preload"
assert_contains "$next_config" "same-origin"

assert_contains "$middleware_file" "withSecurityHeaders"
assert_contains "$middleware_file" "Content-Security-Policy"
assert_contains "$middleware_file" "Cross-Origin-Opener-Policy"
assert_contains "$middleware_file" "Strict-Transport-Security"
assert_contains "$middleware_file" "X-Frame-Options"
assert_contains "$middleware_file" "return withSecurityHeaders(response);"

header_apply_count=$(grep -c "return withSecurityHeaders(response);" "$middleware_file" || true)
if (( header_apply_count < 3 )); then
  echo "[FAIL] middleware must apply security headers on all locale flows (expected >=3, got ${header_apply_count})"
  failures=$((failures + 1))
fi

if (( failures > 0 )); then
  exit 1
fi

echo "[OK] web security header gate passed"
