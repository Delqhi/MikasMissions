create table if not exists identity.parent_controls (
  child_profile_id text primary key,
  autoplay boolean not null,
  chat_enabled boolean not null,
  external_links boolean not null,
  session_limit_minutes integer not null check (session_limit_minutes between 5 and 180),
  bedtime_window text not null,
  safety_mode text not null check (safety_mode in ('strict', 'balanced')),
  updated_at timestamptz not null default now()
);

create table if not exists identity.parent_gate_tokens (
  child_profile_id text primary key,
  gate_token text not null,
  updated_at timestamptz not null default now()
);

create table if not exists identity.parent_gate_challenges (
  challenge_id text primary key,
  parent_user_id uuid not null references identity.parents(id),
  child_profile_id text not null,
  method text not null,
  expires_at timestamptz not null,
  used boolean not null default false,
  created_at timestamptz not null default now()
);

create index if not exists idx_parent_gate_challenges_lookup
on identity.parent_gate_challenges (challenge_id, parent_user_id, child_profile_id);
