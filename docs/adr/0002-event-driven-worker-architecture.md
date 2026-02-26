# ADR 0002: Event-Driven Worker Architecture

## Status
Accepted

## Context
Media processing and recommendation pipelines require asynchronous retries and independent scaling.

## Decision
Use event topics with idempotent workers and outbox-driven publish semantics.

## Consequences
- Failures are isolated and retryable.
- Worker throughput can scale independently from API services.
- Requires strict contract versioning discipline.
