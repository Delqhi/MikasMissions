create table if not exists creator.generated_assets (
  id uuid primary key default gen_random_uuid(),
  run_id uuid not null references creator.workflow_runs(id) on delete cascade,
  asset_id text not null unique,
  source_url text not null,
  qc_status text not null,
  publish_event_id text,
  created_at timestamptz not null default now()
);

create index if not exists idx_creator_generated_assets_run
on creator.generated_assets (run_id, created_at desc);
