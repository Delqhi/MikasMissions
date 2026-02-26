alter table profiles.child_profiles drop constraint if exists child_profiles_age_band_check;

alter table profiles.child_profiles
  add constraint child_profiles_age_band_check
  check (age_band in ('3-5', '6-11', '12-16'));
