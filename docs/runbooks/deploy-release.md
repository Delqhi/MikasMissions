# Deploy and Release Runbook

## Strategy
Trunk-based with progressive rollouts and feature flags.

## Steps
1. Merge PR after green CI.
2. Validate kube access before deploy: `make kube-validate`.
3. Evaluate launch readiness gate (deadline/fallback policy): `make launch-readiness-gate`.
4. Run launch preflight from repo root: `make launch-preflight`.
5. Build images and update manifests.
6. Ensure runtime secrets are present (`platform-secrets.database-url` for DB-backed services/workers).

```bash
kubectl -n mikasmissions-dev create secret generic platform-secrets \
  --from-literal=database-url='postgres://USER:PASSWORD@HOST:5432/DB?sslmode=require' \
  --dry-run=client -o yaml | kubectl apply -f -
```

7. ArgoCD sync to target environment.
8. Validate health endpoints and key business probes.
9. Run 24h soak in staging: `make staging-soak`.
10. Generate launch packet: `make launch-decision-packet`.
11. Promote traffic gradually with stop/go gates (`docs/runbooks/public_launch_2026-03-07.md`), preferably via `make launch-stage`.

## Rollback
Revert image tag to previous healthy release and re-sync ArgoCD.
