#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

budget_file="frontend/web/config/web-vitals-budgets.json"
ci_workflow=".github/workflows/ci.yml"
preflight_script="scripts/launch_preflight.sh"

if [[ ! -f "$budget_file" ]]; then
  echo "[FAIL] missing budget policy file: $budget_file"
  exit 1
fi

node - "$budget_file" <<'NODE'
const fs = require("node:fs");

const path = process.argv[2];
const raw = fs.readFileSync(path, "utf8");
const budget = JSON.parse(raw);
const failures = [];

function assertNum(name) {
  const value = budget[name];
  if (typeof value !== "number" || Number.isNaN(value)) {
    failures.push(`missing numeric field '${name}'`);
    return NaN;
  }
  return value;
}

const lcpMax = assertNum("lcp_ms_max");
const inpMax = assertNum("inp_ms_max");
const clsMax = assertNum("cls_max");
const perfMin = assertNum("lighthouse_performance_score_min");
const a11yMin = assertNum("lighthouse_accessibility_score_min");

if (Number.isFinite(lcpMax) && lcpMax > 2500) {
  failures.push(`lcp_ms_max must be <= 2500 (actual: ${lcpMax})`);
}
if (Number.isFinite(inpMax) && inpMax > 200) {
  failures.push(`inp_ms_max must be <= 200 (actual: ${inpMax})`);
}
if (Number.isFinite(clsMax) && clsMax > 0.1) {
  failures.push(`cls_max must be <= 0.1 (actual: ${clsMax})`);
}
if (Number.isFinite(perfMin) && perfMin < 90) {
  failures.push(`lighthouse_performance_score_min must be >= 90 (actual: ${perfMin})`);
}
if (Number.isFinite(a11yMin) && a11yMin < 95) {
  failures.push(`lighthouse_accessibility_score_min must be >= 95 (actual: ${a11yMin})`);
}

console.log(`[metric] lcp_ms_max=${Number.isFinite(lcpMax) ? lcpMax : "missing"}`);
console.log(`[metric] inp_ms_max=${Number.isFinite(inpMax) ? inpMax : "missing"}`);
console.log(`[metric] cls_max=${Number.isFinite(clsMax) ? clsMax : "missing"}`);
console.log(`[metric] lighthouse_performance_score_min=${Number.isFinite(perfMin) ? perfMin : "missing"}`);
console.log(`[metric] lighthouse_accessibility_score_min=${Number.isFinite(a11yMin) ? a11yMin : "missing"}`);

if (failures.length > 0) {
  for (const failure of failures) {
    console.error(`[FAIL] ${failure}`);
  }
  process.exit(1);
}
NODE

if ! grep -q "check_web_vitals_budgets.sh" "$ci_workflow"; then
  echo "[FAIL] $ci_workflow must include check_web_vitals_budgets.sh"
  exit 1
fi

if ! grep -q "check_web_vitals_budgets.sh" "$preflight_script"; then
  echo "[FAIL] $preflight_script must include check_web_vitals_budgets.sh"
  exit 1
fi

echo "[OK] web vitals budget gate passed"
