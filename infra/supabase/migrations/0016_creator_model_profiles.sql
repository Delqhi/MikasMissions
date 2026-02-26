create table if not exists creator.model_profiles (
  id text primary key,
  provider text not null,
  base_url text not null,
  model_id text not null,
  timeout_ms int not null,
  max_retries int not null,
  safety_preset text not null,
  updated_at timestamptz not null default now()
);

insert into creator.model_profiles (id, provider, base_url, model_id, timeout_ms, max_retries, safety_preset)
values ('nim-default', 'nvidia_nim', 'http://localhost:9000', 'nim-video-v1', 15000, 2, 'kids_strict')
on conflict (id) do nothing;
