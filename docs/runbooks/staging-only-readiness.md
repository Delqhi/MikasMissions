# Staging-Only Readiness Runbook

## Goal
Apply the decision-complete policy: if valid external kube access is missing past the cluster deadline, public go-live is blocked and release status is staging-only.

## Command

```bash
make launch-readiness-gate
```

## Behavior
1. `READY_FOR_STAGING_DEPLOY`
   1. `kube-validate` passed.
   2. External cluster policy passed (no local dev context).
   3. Continue with `make staging-deploy`.
2. `BLOCKED_PENDING_CLUSTER_ACCESS`
   1. `kube-validate` failed before deadline, or kube is local-only (`kind|minikube|k3d|localhost`) before deadline.
   2. Fix to external cluster access and re-run gate.
3. `STAGING_ONLY_READINESS`
   1. `kube-validate` failed after deadline, or only local cluster policy is met after deadline.
   2. Public rollout is blocked until cluster access is restored.

## Config
1. Deadline override (CET with offset):
   1. `LAUNCH_CLUSTER_DEADLINE_CET=2026-02-26T18:00:00+01:00`
2. Evidence path override:
   1. `LAUNCH_READINESS_EVIDENCE_DIR=docs/runbooks/evidence`
3. External cluster policy:
   1. default: `LAUNCH_REQUIRE_EXTERNAL_CLUSTER=true`
   2. local drill override: `LAUNCH_REQUIRE_EXTERNAL_CLUSTER=false`

## Evidence
1. Writes one file per evaluation:
   1. `docs/runbooks/evidence/launch-readiness-<timestamp>.md`
2. kube details log:
   1. `/tmp/launch-readiness-kube.log`
