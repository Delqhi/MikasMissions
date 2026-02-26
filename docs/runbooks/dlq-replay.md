# DLQ and Outbox Replay Runbook

## Goal
Recover failed outbox deliveries safely, with auditability and canary-first rollout.

## Preconditions
1. Root cause is identified and fixed.
2. Fix is deployed to all relevant producers/consumers.
3. Replay scope is approved by engineering + safety owner.
4. `worker-outbox-relay` is running with valid `DATABASE_URL` and `NATS_URL`.

## Inputs
1. `DATABASE_URL` with access to `events.outbox`.
2. `NATS_URL` for event canary observation.
3. Target scope:
   1. single `event_id`, or
   2. bounded batch (`limit`, optional `topic`).

## Procedure
1. List failed rows and capture snapshot evidence.

```bash
DATABASE_URL=... make outbox-replay ARGS="-mode=list-failed -limit=50"
```

2. Optional topic filter for scoped replay:

```bash
DATABASE_URL=... make outbox-replay ARGS="-mode=list-failed -topic=episode.published.v1 -limit=50"
```

3. Dry-run candidate selection (no writes):

```bash
DATABASE_URL=... make outbox-replay ARGS="-mode=requeue-failed -topic=episode.published.v1 -limit=5 -dry-run=true"
```

4. Start canary observation on the expected downstream topic:

```bash
NATS_URL=nats://127.0.0.1:4222 go run ./tools/event-probe -topic "episode.published.v1" -timeout 25s
```

5. Execute canary replay:

```bash
DATABASE_URL=... make outbox-replay ARGS="-mode=requeue-failed -topic=episode.published.v1 -limit=1 -dry-run=false -reset-attempts=true"
```

6. Verify canary outcomes:
   1. `worker-outbox-relay` logs show successful flush.
   2. downstream consumer lag is stable.
   3. no spike in `.dlq.v1` topics.

7. If canary is healthy, replay in bounded batches until queue is drained.

```bash
DATABASE_URL=... make outbox-replay ARGS="-mode=requeue-failed -limit=25 -dry-run=false -reset-attempts=true"
```

8. Re-run listing to confirm closure:

```bash
DATABASE_URL=... make outbox-replay ARGS="-mode=list-failed -limit=50"
```

## Single Event Replay
1. Dry-run check:

```bash
DATABASE_URL=... make outbox-replay ARGS="-mode=requeue-event -event-id=<event-id> -dry-run=true"
```

2. Execute:

```bash
DATABASE_URL=... make outbox-replay ARGS="-mode=requeue-event -event-id=<event-id> -dry-run=false -reset-attempts=true"
```

## Safety Rules
1. Never replay unvalidated payloads.
2. Never replay without canary when root cause was schema/contract related.
3. Use bounded batches (`limit`) and re-evaluate between batches.
4. Keep an evidence log (command, timestamp, operator, output summary).

## Failure Handling
1. If replayed rows fail again immediately, stop and reopen incident.
2. If DLQ volume grows after replay, pause replay and roll back rollout.
3. If `worker-outbox-relay` is down, do not perform replay writes until relay health is restored.
