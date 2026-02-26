create schema if not exists creator;
create schema if not exists audit;

create table if not exists creator.workflow_templates (
  id uuid primary key default gen_random_uuid(),
  name text not null,
  description text not null,
  content_suitability text not null,
  age_band text not null check (age_band in ('3-5', '6-11', '12-16')),
  steps jsonb not null,
  model_profile_id text not null,
  safety_profile text not null,
  version int not null default 1,
  created_by text,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table if not exists creator.workflow_template_versions (
  id bigserial primary key,
  workflow_id uuid not null references creator.workflow_templates(id) on delete cascade,
  version int not null,
  snapshot jsonb not null,
  created_by text,
  created_at timestamptz not null default now()
);

create index if not exists idx_creator_workflow_templates_name
on creator.workflow_templates (name);
