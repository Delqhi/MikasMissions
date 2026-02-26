# ADR 0003: Micro-File Guardrails

## Status
Accepted

## Context
Large files and mixed concerns quickly create monolith behavior in codebases.

## Decision
Enforce max 250 lines per Go file in CI and require event schema + contract tests.

## Consequences
- Better readability and maintainability.
- Slight overhead in splitting logic into focused files.
- Predictable review quality across services.
