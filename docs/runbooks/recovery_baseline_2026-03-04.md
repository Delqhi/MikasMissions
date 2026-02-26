# Recovery Baseline Runbook (R0)

## Scope
Recovery Freeze baseline for `2026-02-26` to `2026-03-04`.
This runbook defines mandatory checks and evidence for Gate `R0`.

## Mandatory Exit Checks
Run from `/platform`:

```bash
make guard
make test
make build
make contract-check
make e2e-smoke
```

Run from `/frontend/web`:

```bash
npm ci
npm run build
```

## Recovery Invariants
1. `make guard` is green with no bypasses.
2. `make test` is green and deterministic (no race-sensitive flakes).
3. `make build` is green.
4. `make contract-check` is green (`generated` and `breaking` checks).
5. Web production build is green.
6. Local smoke reaches all `/v1/*` critical paths including billing entitlement check.

## Data and Queue Baseline
1. SQL migrations include `events.outbox` and `events.idempotency_keys`.
2. Worker idempotency supports scoped keys with Postgres-backed mode via `DATABASE_URL`.
3. `creator-studio-service` uses persistent outbox in Postgres mode and in-memory outbox for local fallback.
4. `worker-outbox-relay` flushes pending outbox events to the event bus with retry and DLQ handling.
5. `tools/outbox-replay` supports canary-first requeue for failed outbox rows.
6. `DATABASE_URL` empty means local in-memory fallback for developer velocity.
7. In production (`GO_ENV=production` or `PERSISTENCE_MODE=strict`) critical services fail fast without `DATABASE_URL`.

## Evidence Template
For each check, capture:
1. Command.
2. UTC timestamp.
3. Exit code.
4. Short output summary.

Example:

```text
2026-03-04T16:25:10Z | make guard | exit 0 | micro-file/event/API guard all passed
```

## Failure Policy
1. Any red mandatory check blocks feature work.
2. Fixes are prioritized in this order:
   1. Deterministic tests.
   2. Contract compatibility.
   3. Runtime stability.
3. No merge with unresolved hard-fail in R0 mandatory checks.

## Latest Baseline Evidence (Recorded on 25. Februar 2026)
1. `make guard` -> exit `0`.
2. `make test` -> exit `0` after fixing `worker-ingest` idempotency test race.
3. `make build` -> exit `0`.
4. `make contract-check` -> exit `0`.
5. `frontend/web npm run build` -> exit `0`.
6. `make e2e-smoke` -> exit `0` via auto-started local core stack and full `/v1` smoke flow.
