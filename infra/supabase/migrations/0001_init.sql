-- Core schemas
create schema if not exists identity;
create schema if not exists profiles;
create schema if not exists catalog;
create schema if not exists progress;
create schema if not exists moderation;
create schema if not exists billing;

-- Identity
create table if not exists identity.parents (
  id uuid primary key,
  email text not null unique,
  country text not null,
  language text not null,
  created_at timestamptz not null default now()
);

create table if not exists identity.consents (
  id uuid primary key,
  parent_id uuid not null references identity.parents(id),
  method text not null,
  verified boolean not null default false,
  created_at timestamptz not null default now()
);

-- Profiles
create table if not exists profiles.child_profiles (
  id uuid primary key,
  parent_id uuid not null references identity.parents(id),
  display_name text not null,
  age_band text not null check (age_band in ('3-5', '6-11', '12-16')),
  avatar text not null,
  created_at timestamptz not null default now()
);

-- Catalog
create table if not exists catalog.episodes (
  id text primary key,
  show_id text not null,
  title text not null,
  summary text not null,
  age_band text not null,
  duration_ms bigint not null,
  learning_tags text[] not null default '{}',
  playback_ready boolean not null default false,
  thumbnail_url text not null default '',
  published_at timestamptz
);

-- Progress
create table if not exists progress.watch_events (
  id bigserial primary key,
  child_profile_id uuid not null references profiles.child_profiles(id),
  episode_id text not null references catalog.episodes(id),
  watch_ms bigint not null,
  event_time timestamptz not null,
  created_at timestamptz not null default now()
);

-- Billing
create table if not exists billing.subscriptions (
  id uuid primary key,
  parent_id uuid not null references identity.parents(id),
  plan text not null,
  active boolean not null default true,
  created_at timestamptz not null default now()
);

-- RLS foundation
alter table profiles.child_profiles enable row level security;

create policy child_profiles_parent_read
on profiles.child_profiles
for select
using (parent_id::text = auth.uid()::text);

create policy child_profiles_parent_write
on profiles.child_profiles
for insert
with check (parent_id::text = auth.uid()::text);
