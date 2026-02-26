# Staging 24h Soak Runbook

## Goal
Prove staging stability for at least 24 hours with repeated auth and core flow probes.

## Command

```bash
SOAK_HOURS=24 SOAK_INTERVAL_MINUTES=60 make staging-soak
```

## Defaults
1. `SOAK_HOURS=24`
2. `SOAK_INTERVAL_MINUTES=60`
3. `SOAK_FAIL_FAST=true`
4. Evidence output:
   1. `docs/runbooks/evidence/staging-soak-<timestamp>/summary.md`
   2. Per-iteration logs for `e2e-smoke` and `e2e-auth-smoke`

## Gates
1. No failed iterations.
2. No SEV-1/SEV-2 incident during soak window.
3. If one check fails and `SOAK_FAIL_FAST=true`, soak is failed immediately.

## Failure Procedure
1. Stop progression to rollout stage.
2. Run rollback drill:

```bash
make staging-rollback
```

3. Re-run health and smoke checks:

```bash
make e2e-smoke
make e2e-auth-smoke
```

4. Log incident timeline and corrective actions.
