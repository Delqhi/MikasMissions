#!/usr/bin/env bash
set -euo pipefail

max_lines=250
violations=0

if command -v rg >/dev/null 2>&1; then
  go_files_cmd=(rg --files -g '*.go' -g '!bin/**' -g '!vendor/**' -g '!libs/contracts-api/generated/**')
else
  go_files_cmd=(find . -name '*.go' \
    -not -path './bin/*' \
    -not -path './vendor/*' \
    -not -path './libs/contracts-api/generated/*')
fi

while IFS= read -r file; do
  line_count=$(wc -l < "$file" | tr -d ' ')
  if (( line_count > max_lines )); then
    echo "[FAIL] $file has $line_count lines (max $max_lines)"
    violations=1
  fi
done < <("${go_files_cmd[@]}" | sort)

if (( violations != 0 )); then
  exit 1
fi

echo "[OK] micro-file line guard passed"
