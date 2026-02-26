export type AgeBand = "3-5" | "6-11" | "12-16";

export type SafetyMode = "strict" | "balanced";

export type ContentSuitability = "early" | "core" | "teen";

export type KidsMode = "early" | "core" | "teen";

export type RailItem = {
  episode_id: string;
  title: string;
  summary: string;
  thumbnail_url: string;
  duration_ms: number;
  age_band: AgeBand;
  content_suitability: ContentSuitability;
  learning_tags: string[];
  reason_code: string;
  safety_applied: boolean;
  age_fit_score: number;
};

export type KidsHomeResponse = {
  child_profile_id: string;
  mode: KidsMode;
  safety_mode: SafetyMode;
  primary_actions: string[];
  rails: RailItem[];
};

export type KidsProgressResponse = {
  child_profile_id: string;
  watched_minutes_today: number;
  watched_minutes_7d: number;
  completion_percent: number;
  mission_streak_days: number;
  session_limit_minutes: number;
  session_minutes_used: number;
  session_capped: boolean;
  last_episode_id: string;
};

export type ParentalControls = {
  autoplay: boolean;
  chat_enabled: boolean;
  external_links: boolean;
  session_limit_minutes: number;
  bedtime_window: string;
  safety_mode: SafetyMode;
};

export type ParentControlsResponse = {
  child_profile_id: string;
  controls: ParentalControls;
};

export type ProfileCard = {
  profile_id: string;
  name: string;
  age_band: AgeBand;
  mode: KidsMode;
  subtitle: string;
  href: string;
};
