#!/usr/bin/env bash
set -euo pipefail

kubeconfig=${KUBECONFIG:-$HOME/.kube/config}
namespace=${K8S_NAMESPACE:-mikasmissions-dev}
kubectl_timeout=${KUBE_REQUEST_TIMEOUT:-8s}
kubectl_command_timeout=${KUBE_COMMAND_TIMEOUT_SECONDS:-12}

require_cmd() {
  local cmd=$1
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "missing required command: $cmd"
    exit 1
  fi
}

base64_flag() {
  if printf 'QQ==' | base64 --decode >/dev/null 2>&1; then
    echo "--decode"
    return
  fi
  echo "-D"
}

read_field() {
  local jsonpath=$1
  kubectl --kubeconfig "$kubeconfig" config view --minify -o "jsonpath=$jsonpath" 2>/dev/null || true
}

validate_base64_field() {
  local key=$1
  local jsonpath=$2
  local value
  value=$(read_field "$jsonpath")
  if [[ -z "$value" ]]; then
    return
  fi
  if ! printf '%s' "$value" | base64 "$(base64_flag)" >/dev/null 2>&1; then
    echo "invalid base64 for field: $key in $kubeconfig"
    echo "repair kube credentials and retry (or set KUBECONFIG to a valid file)."
    exit 1
  fi
}

read_server_host() {
  local server
  local host
  server=$(read_field '{.clusters[0].cluster.server}')
  host=${server#https://}
  host=${host#http://}
  host=${host%%/*}
  host=${host%%:*}
  printf '%s' "$host"
}

run_with_timeout() {
  local timeout_seconds=$1
  shift

  "$@" &
  local pid=$!
  local elapsed=0

  while kill -0 "$pid" >/dev/null 2>&1; do
    if (( elapsed >= timeout_seconds )); then
      kill "$pid" >/dev/null 2>&1 || true
      sleep 1
      if kill -0 "$pid" >/dev/null 2>&1; then
        kill -9 "$pid" >/dev/null 2>&1 || true
      fi
      wait "$pid" >/dev/null 2>&1 || true
      return 124
    fi
    sleep 1
    elapsed=$((elapsed + 1))
  done

  wait "$pid"
}

require_cmd kubectl
require_cmd base64

if [[ ! -f "$kubeconfig" ]]; then
  echo "kube config not found: $kubeconfig"
  exit 1
fi

if ! kubectl --kubeconfig "$kubeconfig" config current-context >/dev/null 2>&1; then
  echo "kubectl context is not configured or kube config is invalid."
  exit 1
fi

validate_base64_field "certificate-authority-data" '{.clusters[0].cluster.certificate-authority-data}'
validate_base64_field "client-certificate-data" '{.users[0].user.client-certificate-data}'
validate_base64_field "client-key-data" '{.users[0].user.client-key-data}'

server_host=$(read_server_host)
if [[ -z "$server_host" ]]; then
  echo "kube server endpoint could not be resolved from current context."
  exit 1
fi
if [[ "$server_host" == "kubernetes.default.svc.cluster.local" ]]; then
  echo "kube config points to in-cluster endpoint: $server_host"
  echo "use an external kubeconfig from your cluster provider (EKS/GKE/AKS/k3s/kind)."
  exit 1
fi

if ! run_with_timeout "$kubectl_command_timeout" kubectl --kubeconfig "$kubeconfig" --request-timeout="$kubectl_timeout" get ns >/dev/null 2>&1; then
  echo "kubectl is configured but cluster is not reachable."
  exit 1
fi

if run_with_timeout "$kubectl_command_timeout" kubectl --kubeconfig "$kubeconfig" --request-timeout="$kubectl_timeout" get ns "$namespace" >/dev/null 2>&1; then
  echo "[OK] kube config valid and namespace reachable: $namespace"
else
  echo "[OK] kube config valid and cluster reachable; namespace $namespace will be created on deploy."
fi
