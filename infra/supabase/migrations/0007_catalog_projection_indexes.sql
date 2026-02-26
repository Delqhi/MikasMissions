create index if not exists idx_catalog_episodes_age_band_published
on catalog.episodes (age_band, published_at desc);

create index if not exists idx_catalog_episodes_published
on catalog.episodes (published_at desc);
