# ADR 0001: Go + Supabase Foundation

## Status
Accepted

## Context
The platform needs rapid MVP delivery with strong scale-up options and strict child-safety controls.

## Decision
Use Go for backend services and Supabase Postgres/Auth/Storage for initial platform foundation.

## Consequences
- Fast implementation with simple operational overhead in early phases.
- Clear migration path to specialized components at higher scale.
- RLS and SQL controls enforce data boundaries early.
