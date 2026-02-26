create table if not exists identity.admin_users (
  id uuid primary key default gen_random_uuid(),
  email text not null unique,
  password_hash text not null,
  last_login_at timestamptz,
  created_at timestamptz not null default now()
);

create index if not exists idx_identity_admin_users_email_lower
on identity.admin_users (lower(email));
