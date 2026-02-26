alter table identity.parent_gate_tokens
  add column if not exists expires_at timestamptz not null default (now() + interval '5 minutes'),
  add column if not exists consumed_at timestamptz;

update identity.parent_gate_tokens
set expires_at = coalesce(expires_at, now() + interval '5 minutes')
where expires_at is null;

create index if not exists idx_parent_gate_tokens_valid
on identity.parent_gate_tokens (child_profile_id, expires_at, consumed_at);

