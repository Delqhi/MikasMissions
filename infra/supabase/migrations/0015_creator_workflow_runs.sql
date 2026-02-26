create table if not exists creator.workflow_runs (
  id uuid primary key default gen_random_uuid(),
  workflow_id uuid not null references creator.workflow_templates(id) on delete cascade,
  status text not null,
  input_payload jsonb not null default '{}'::jsonb,
  priority text not null default 'normal',
  auto_publish boolean not null default false,
  created_by text,
  last_error text,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table if not exists creator.workflow_run_steps (
  id bigserial primary key,
  run_id uuid not null references creator.workflow_runs(id) on delete cascade,
  step text not null,
  status text not null,
  message text not null,
  created_at timestamptz not null default now()
);

create index if not exists idx_creator_workflow_runs_status
on creator.workflow_runs (status, updated_at desc);

create index if not exists idx_creator_workflow_run_steps_run
on creator.workflow_run_steps (run_id, id);
