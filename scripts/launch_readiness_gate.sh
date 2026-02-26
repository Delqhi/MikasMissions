#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

# Decision-complete default from plan:
# If kube access is still missing after 2026-02-26 18:00 CET,
# launch is blocked and status switches to staging-only readiness.
deadline_cet=${LAUNCH_CLUSTER_DEADLINE_CET:-2026-02-26T18:00:00+01:00}
evidence_dir=${LAUNCH_READINESS_EVIDENCE_DIR:-docs/runbooks/evidence}
require_external_cluster=${LAUNCH_REQUIRE_EXTERNAL_CLUSTER:-true}
kubeconfig=${KUBECONFIG:-$HOME/.kube/config}
stamp=$(date -u +%Y%m%dT%H%M%SZ)
report="$evidence_dir/launch-readiness-$stamp.md"

mkdir -p "$evidence_dir"

is_local_context() {
  local ctx=$1
  if [[ "$ctx" =~ ^kind- ]]; then
    return 0
  fi
  if [[ "$ctx" == "minikube" ]]; then
    return 0
  fi
  if [[ "$ctx" =~ ^k3d- ]]; then
    return 0
  fi
  return 1
}

is_local_server_host() {
  local host=$1
  if [[ "$host" == "localhost" ]]; then
    return 0
  fi
  if [[ "$host" == "0.0.0.0" ]]; then
    return 0
  fi
  if [[ "$host" == "::1" ]]; then
    return 0
  fi
  if [[ "$host" == "kubernetes.default.svc.cluster.local" ]]; then
    return 0
  fi
  if [[ "$host" =~ ^127\. ]]; then
    return 0
  fi
  return 1
}

normalize_offset_colon() {
  local ts=$1
  # 2026-02-26T18:00:00+01:00 -> 2026-02-26T18:00:00+0100
  if [[ "$ts" =~ ^(.+[0-9]{2}):([0-9]{2})$ ]]; then
    echo "${BASH_REMATCH[1]}${BASH_REMATCH[2]}"
    return
  fi
  echo "$ts"
}

current_epoch=$(date -u +%s)
deadline_epoch=$(date -j -f "%Y-%m-%dT%H:%M:%S%z" "$deadline_cet" +%s 2>/dev/null || true)
if [[ -z "$deadline_epoch" ]]; then
  deadline_cet_normalized=$(normalize_offset_colon "$deadline_cet")
  deadline_epoch=$(date -j -f "%Y-%m-%dT%H:%M:%S%z" "$deadline_cet_normalized" +%s 2>/dev/null || true)
fi

if [[ -z "$deadline_epoch" ]]; then
  echo "invalid LAUNCH_CLUSTER_DEADLINE_CET: $deadline_cet"
  exit 1
fi

kube_status="FAIL"
external_policy_status="SKIPPED"
current_context=""
host=""

if make kube-validate >/tmp/launch-readiness-kube.log 2>&1; then
  kube_status="PASS"
  if [[ "$require_external_cluster" == "true" ]]; then
    external_policy_status="FAIL"
    current_context=$(kubectl --kubeconfig "$kubeconfig" config current-context 2>/dev/null || true)
    server=$(kubectl --kubeconfig "$kubeconfig" config view --minify -o 'jsonpath={.clusters[0].cluster.server}' 2>/dev/null || true)
    host=${server#https://}
    host=${host#http://}
    host=${host%%/*}
    host=${host%%:*}

    if [[ -z "$current_context" || -z "$host" ]]; then
      {
        echo "external cluster policy check failed: unable to resolve context/host"
        echo "current_context='$current_context'"
        echo "server='$server'"
      } >>/tmp/launch-readiness-kube.log
    elif is_local_context "$current_context" || is_local_server_host "$host"; then
      {
        echo "external cluster policy check failed: local/dev cluster context is not allowed for public launch"
        echo "current_context='$current_context'"
        echo "server_host='$host'"
      } >>/tmp/launch-readiness-kube.log
    else
      external_policy_status="PASS"
    fi
  fi
fi

if [[ "$kube_status" == "PASS" && "$external_policy_status" != "FAIL" ]]; then
  cat >"$report" <<EOF
# Launch Readiness Gate

- status: READY_FOR_STAGING_DEPLOY
- generated_utc: $(date -u +"%Y-%m-%dT%H:%M:%SZ")
- deadline_cet: $deadline_cet
- kube_validate: $kube_status
- external_cluster_policy: $external_policy_status
EOF
  if [[ "$external_policy_status" == "PASS" ]]; then
    {
      echo "- current_context: $current_context"
      echo "- server_host: $host"
    } >>"$report"
  fi
  echo "[OK] launch readiness gate passed: $report"
  exit 0
fi

if (( current_epoch < deadline_epoch )); then
  cat >"$report" <<EOF
# Launch Readiness Gate

- status: BLOCKED_PENDING_CLUSTER_ACCESS
- generated_utc: $(date -u +"%Y-%m-%dT%H:%M:%SZ")
- deadline_cet: $deadline_cet
- kube_validate: $kube_status
- external_cluster_policy: $external_policy_status
- policy: launch still possible before deadline after external cluster access is fixed
- details_log: /tmp/launch-readiness-kube.log
EOF
  echo "[WARN] launch readiness currently blocked (before deadline): $report"
  exit 1
fi

cat >"$report" <<EOF
# Launch Readiness Gate

- status: STAGING_ONLY_READINESS
- generated_utc: $(date -u +"%Y-%m-%dT%H:%M:%SZ")
- deadline_cet: $deadline_cet
- kube_validate: $kube_status
- external_cluster_policy: $external_policy_status
- policy: public go-live blocked until external kube access is restored
- details_log: /tmp/launch-readiness-kube.log
EOF
echo "[STOP] deadline passed without valid kube access; staging-only readiness activated: $report"
exit 2
