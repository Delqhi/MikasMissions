# Cluster Unblock Runbook (Kube Config)

## Goal
Restore a valid Kubernetes context so `staging-deploy` and `staging-rollback` are executable.
For public go-live readiness, this must be an external cluster context (not `kind|minikube|k3d`).

## Steps
1. Validate local kube configuration:

```bash
make kube-validate
```

2. If validation fails with `invalid base64`:
   1. Replace `~/.kube/config` with valid credentials from your cluster provider.
   2. Re-run `make kube-validate`.

3. Confirm cluster and namespace visibility:

```bash
kubectl config current-context
kubectl get ns
kubectl get ns mikasmissions-dev
```

4. Dry-run deploy validation:

```bash
make staging-deploy-dry-run
```

## Required Environment Variables for Real Deploy
1. `DATABASE_URL`
2. `AUTH_JWT_SECRET`
3. Optional:
   1. `AUTH_JWKS_URL`
   2. `ADMIN_BOOTSTRAP_EMAIL`
   3. `ADMIN_BOOTSTRAP_PASSWORD`

## Exit Criteria
1. `make kube-validate` passes.
2. `make staging-deploy-dry-run` passes.
3. Required env vars are available for `make staging-deploy`.
4. `make launch-readiness-gate` reports `READY_FOR_STAGING_DEPLOY` with external policy pass.
