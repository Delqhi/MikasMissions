# MikasMissions Platform

Event-driven micro-modular foundation for the MikasMissions kids streaming platform.

## Quick Start

```bash
make guard
make test
make build
```

## Cloudflare Deploy (Workers + OpenNext)

Local deploy from `frontend/web`:

```bash
pnpm install --frozen-lockfile
NEXT_PUBLIC_API_BASE_URL=https://api.example.com NEXT_PUBLIC_USE_API_FALLBACKS=false pnpm run deploy
```

Required CI secrets for `.github/workflows/cloudflare-deploy.yml`:
- `CLOUDFLARE_API_TOKEN`
- `CLOUDFLARE_ACCOUNT_ID`
- `NEXT_PUBLIC_API_BASE_URL`
- `CLOUDFLARE_WORKERS_DEV_URL`

## R0 Local Exit Check

```bash
make guard
make test
make build
make contract-check
make e2e-smoke
cd frontend/web
pnpm install --frozen-lockfile
pnpm run build
```

## Contract Source Of Truth

The canonical public API spec lives in:

- `libs/contracts-api/openapi/v1.yaml`

Generate and validate contract artifacts:

```bash
make contract
make contract-check
```

## Persistence Mode (Identity + Profile + Progress + Billing + Queue Guards)

By default, services run with in-memory repositories.
To enable Postgres-backed repositories for `identity-service`, `profile-service`, `progress-service`, and `billing-service`:

```bash
export DATABASE_URL=postgres://postgres:postgres@127.0.0.1:5432/temporal?sslmode=disable
make migrate-db
make run-identity
make run-profile
make run-progress
make run-billing
```

The same `DATABASE_URL` also enables persistent idempotency keys for workers via `events.idempotency_keys`.
`creator-studio-service` writes events via persistent outbox (`events.outbox`) when `DATABASE_URL` is set.
In this mode, run `worker-outbox-relay` to publish queued outbox events to NATS.
Use `tools/outbox-replay` to inspect and safely requeue failed outbox rows.

`playback-service` enforces entitlements server-side through `BILLING_URL` (default: `http://127.0.0.1:8089`).

For live environments, disable in-memory fallbacks:

```bash
export GO_ENV=production
# or: export PERSISTENCE_MODE=strict
```

With strict mode, services/workers that need persistent storage fail fast if `DATABASE_URL` is missing or unavailable.

## Run Services (Local)

In separate terminals:

```bash
make run-gateway
make run-identity
make run-profile
make run-catalog
make run-recommendation
make run-playback
make run-progress
```

Then run:

```bash
make e2e-smoke
```

`smoke_local.sh` now targets the public gateway on `http://localhost:8080` by default.

## Container Smoke Profile

```bash
make compose-smoke
```

## Web Fallback Mode (Dev Only)

The web app uses strict API mode by default. To enable mock fallback data explicitly:

```bash
cd frontend/web
NEXT_PUBLIC_USE_API_FALLBACKS=true pnpm run dev
```

## Launch Preflight (Public Go-Live)

```bash
make launch-preflight
```

Validate Kubernetes access before deploy:

```bash
make kube-validate
```

Evaluate launch readiness gate (deadline + staging-only fallback):

```bash
make launch-readiness-gate
```

By default, launch readiness requires an external (non-local) cluster context.
`kind`/`minikube`/`k3d` contexts are treated as staging-only readiness.
For local drills only, you can bypass this policy:

```bash
LAUNCH_REQUIRE_EXTERNAL_CLUSTER=false make launch-readiness-gate
```

Generate launch decision artifacts:

```bash
make launch-decision-packet
```

By default, the decision packet runs preflight with a 30-minute timeout (`LAUNCH_PREFLIGHT_TIMEOUT_SECONDS=1800`).
Override if needed:

```bash
LAUNCH_PREFLIGHT_TIMEOUT_SECONDS=3600 make launch-decision-packet
```

Run the essential "works today" local readiness flow in one command:

```bash
make today-ready
```

Defaults:
1. Uses local readiness mode (`TODAY_SKIP_EXTERNAL_CLUSTER_POLICY=true`).
2. Runs a 1-hour soak (`SOAK_HOURS=1`, `SOAK_INTERVAL_MINUTES=60`).
3. Generates a decision packet after checks.

Run 24h staging soak (hourly smoke/auth probes):

```bash
SOAK_HOURS=24 SOAK_INTERVAL_MINUTES=60 make staging-soak
```

For rollout stage decisions:

```bash
make launch-go-nogo ARGS="--stage 10 --read-p95-ms 285 --write-p95-ms 430 --error-5xx-percent 0.22"
```

Rollout stage execution with evidence logging and automatic rollback on STOP:

```bash
make launch-stage ARGS="--stage 10 --read-p95-ms 285 --write-p95-ms 430 --error-5xx-percent 0.22 --owner sre-oncall --notes rollout-window-2"
```

Staging deploy/rollback:

```bash
export DATABASE_URL='postgres://USER:PASSWORD@HOST:5432/DB?sslmode=require'
export AUTH_JWT_SECRET='replace-with-strong-secret'
make staging-deploy
```

```bash
make staging-rollback
```

If no cluster context is configured yet:

```bash
make staging-deploy-dry-run
```

## Gateway Auth Modes

- `AUTH_MODE=permissive` (default): protected routes still work without token (legacy/local mode).
- `AUTH_MODE=enforce`: protected routes require bearer JWT with role mapping (`parent|child|service`).

Token verification sources:

- `AUTH_JWT_SECRET` for HS256 tokens.
- `AUTH_JWKS_URL` or `SUPABASE_JWKS_URL` for RS256/JWKS validation.
- optional strict claim checks via `AUTH_JWT_ISSUER` and `AUTH_JWT_AUDIENCE`.

## Service-Side Role Enforcement

`identity-service`, `profile-service`, `progress-service`, and `playback-service` enforce role headers:

- required headers in `AUTH_MODE=enforce`: `X-Auth-Role` and optional `X-Auth-Parent-User-ID`
- accepted roles per route follow gateway policy (`parent|child|service`)
- in `AUTH_MODE=permissive` missing headers are allowed for local compatibility

## Phase-B Event Pipeline Demo

Run the creator service plus `ingest -> transcode -> policy -> publish` workers:

```bash
./scripts/run_phase_b_stack.sh
```

In another terminal:

```bash
./scripts/smoke_phase_b.sh
```

The stack script starts NATS via Docker (preferred) or local `nats-server` binary.

## Outbox Replay Operations

Preview failed rows:

```bash
DATABASE_URL=postgres://postgres:postgres@127.0.0.1:5432/temporal?sslmode=disable \
  make outbox-replay ARGS="-mode=list-failed -limit=20"
```

Dry-run requeue candidates:

```bash
DATABASE_URL=postgres://postgres:postgres@127.0.0.1:5432/temporal?sslmode=disable \
  make outbox-replay ARGS="-mode=requeue-failed -limit=10 -dry-run=true"
```

Execute requeue:

```bash
DATABASE_URL=postgres://postgres:postgres@127.0.0.1:5432/temporal?sslmode=disable \
  make outbox-replay ARGS="-mode=requeue-failed -limit=10 -dry-run=false -reset-attempts=true"
```

## Implemented Foundation

- Modular Go services for v1 public APIs
- Worker binaries for core media/recommendation pipeline
- Event contracts (`libs/contracts-events/schemas`) and contract tests
- Supabase SQL bootstrap with RLS foundations
- Kubernetes, Terraform, ArgoCD scaffolds
- CI guardrails for micro-file constraints
