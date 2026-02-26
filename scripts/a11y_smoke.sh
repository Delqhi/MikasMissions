#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

required_files=(
  "frontend/web/app/globals.css"
  "frontend/web/app/parents/page.tsx"
  "frontend/web/app/parents/onboarding/page.tsx"
  "frontend/web/app/page.tsx"
)

for file in "${required_files[@]}"; do
  if [[ ! -f "$file" ]]; then
    echo "[FAIL] missing required UI file: $file"
    exit 1
  fi
done

if ! rg -q "focus-visible" frontend/web/app/globals.css; then
  echo "[FAIL] missing focus-visible styles"
  exit 1
fi

if ! rg -q "<main" frontend/web/app/parents/page.tsx frontend/web/app/parents/onboarding/page.tsx frontend/web/app/page.tsx; then
  echo "[FAIL] expected semantic main landmark in core pages"
  exit 1
fi

echo "[OK] a11y smoke passed"
