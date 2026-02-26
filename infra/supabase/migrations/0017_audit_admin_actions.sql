create table if not exists audit.admin_actions (
  id bigserial primary key,
  admin_user_id text not null,
  action text not null,
  resource_type text not null,
  resource_id text not null,
  payload jsonb not null default '{}'::jsonb,
  created_at timestamptz not null default now()
);

create index if not exists idx_audit_admin_actions_created
on audit.admin_actions (created_at desc);
