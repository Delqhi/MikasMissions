alter table identity.parents
  add column if not exists password_hash text,
  add column if not exists last_login_at timestamptz;

update identity.parents
set password_hash = coalesce(password_hash, '')
where password_hash is null;

create index if not exists idx_identity_parents_email_lower
on identity.parents (lower(email));

create index if not exists idx_profiles_child_profiles_parent
on profiles.child_profiles (parent_id);
