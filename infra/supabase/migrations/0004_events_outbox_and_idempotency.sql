create schema if not exists events;

create table if not exists events.outbox (
  id bigserial primary key,
  event_id text not null unique,
  topic text not null,
  payload jsonb not null,
  status text not null default 'pending' check (status in ('pending', 'published', 'failed')),
  attempts integer not null default 0,
  available_at timestamptz not null default now(),
  last_error text,
  published_at timestamptz,
  created_at timestamptz not null default now()
);

create index if not exists idx_events_outbox_pending
on events.outbox (status, available_at, id);

create table if not exists events.idempotency_keys (
  consumer_scope text not null,
  event_id text not null,
  expires_at timestamptz not null,
  created_at timestamptz not null default now(),
  primary key (consumer_scope, event_id)
);

create index if not exists idx_events_idempotency_expiry
on events.idempotency_keys (expires_at);
