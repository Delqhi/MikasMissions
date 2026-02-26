import type {
  KidsHomeResponse,
  KidsMode,
  KidsProgressResponse,
  ParentControlsResponse,
  ProfileCard,
  RailItem
} from "./experience_types";

const baseRail = {
  thumbnail_url: "https://images.unsplash.com/photo-1502082553048-f009c37129b9?auto=format&fit=crop&w=900&q=80",
  duration_ms: 540000,
  safety_applied: true,
  age_fit_score: 0.94
};

const earlyRails: RailItem[] = [
  {
    episode_id: "ep-rainbow-quest",
    title: "Rainbow Quest",
    summary: "Color matching adventure with guided voice prompts.",
    age_band: "3-5",
    learning_tags: ["colors", "coordination"],
    reason_code: "guided_story",
    content_suitability: "early",
    ...baseRail
  },
  {
    episode_id: "ep-sound-garden",
    title: "Sound Garden",
    summary: "Find friendly sounds and learn animal voices.",
    age_band: "3-5",
    learning_tags: ["listening", "animals"],
    reason_code: "safe_curation",
    content_suitability: "early",
    ...baseRail
  }
];

const coreRails: RailItem[] = [
  {
    episode_id: "ep-space-builders",
    title: "Space Builders",
    summary: "Mission-based series on planets and engineering basics.",
    age_band: "6-11",
    learning_tags: ["science", "problem_solving"],
    reason_code: "learning_path",
    content_suitability: "core",
    ...baseRail
  },
  {
    episode_id: "ep-logic-lab",
    title: "Logic Lab",
    summary: "Short puzzle episodes with progress checkpoints.",
    age_band: "6-11",
    learning_tags: ["logic", "math"],
    reason_code: "progress_resume",
    content_suitability: "core",
    ...baseRail
  },
  {
    episode_id: "ep-earth-club",
    title: "Earth Club",
    summary: "Hands-on sustainability challenges for kids.",
    age_band: "6-11",
    learning_tags: ["nature", "teamwork"],
    reason_code: "safe_explore",
    content_suitability: "core",
    ...baseRail
  }
];

const teenRails: RailItem[] = [
  {
    episode_id: "ep-future-studio",
    title: "Future Studio",
    summary: "Creator stories from design, coding, and science labs.",
    age_band: "12-16",
    learning_tags: ["creativity", "career_sparks"],
    reason_code: "interest_match",
    content_suitability: "teen",
    ...baseRail
  },
  {
    episode_id: "ep-debate-room",
    title: "Debate Room",
    summary: "Structured argument episodes with moderation safety.",
    age_band: "12-16",
    learning_tags: ["critical_thinking", "communication"],
    reason_code: "age_appropriate_explore",
    content_suitability: "teen",
    ...baseRail
  }
];

export const kidsHomeFallback: Record<KidsMode, KidsHomeResponse> = {
  early: {
    child_profile_id: "child-early-01",
    mode: "early",
    safety_mode: "strict",
    primary_actions: ["Start", "Mission", "Favorites", "Pause"],
    rails: earlyRails
  },
  core: {
    child_profile_id: "child-core-01",
    mode: "core",
    safety_mode: "strict",
    primary_actions: ["Resume", "Mission", "Explore", "Progress", "Library"],
    rails: coreRails
  },
  teen: {
    child_profile_id: "child-teen-01",
    mode: "teen",
    safety_mode: "strict",
    primary_actions: ["Watch", "Watchlist", "Explore", "Learn", "Report"],
    rails: teenRails
  }
};

export const kidsProgressFallback: Record<KidsMode, KidsProgressResponse> = {
  early: {
    child_profile_id: "child-early-01",
    watched_minutes_today: 18,
    watched_minutes_7d: 90,
    completion_percent: 42,
    mission_streak_days: 1,
    session_limit_minutes: 35,
    session_minutes_used: 18,
    session_capped: false,
    last_episode_id: "ep-rainbow-quest"
  },
  core: {
    child_profile_id: "child-core-01",
    watched_minutes_today: 32,
    watched_minutes_7d: 140,
    completion_percent: 61,
    mission_streak_days: 3,
    session_limit_minutes: 50,
    session_minutes_used: 32,
    session_capped: false,
    last_episode_id: "ep-space-builders"
  },
  teen: {
    child_profile_id: "child-teen-01",
    watched_minutes_today: 44,
    watched_minutes_7d: 210,
    completion_percent: 57,
    mission_streak_days: 2,
    session_limit_minutes: 75,
    session_minutes_used: 44,
    session_capped: false,
    last_episode_id: "ep-future-studio"
  }
};

export const parentControlsFallback: ParentControlsResponse = {
  child_profile_id: "child-core-01",
  controls: {
    autoplay: false,
    chat_enabled: false,
    external_links: false,
    session_limit_minutes: 50,
    bedtime_window: "20:00-07:00",
    safety_mode: "strict"
  }
};

export const demoProfiles: ProfileCard[] = [
  {
    profile_id: "child-early-01",
    name: "Mika Mini",
    age_band: "3-5",
    mode: "early",
    subtitle: "Audio-guided missions",
    href: "/kids/early"
  },
  {
    profile_id: "child-core-01",
    name: "Mika Explorer",
    age_band: "6-11",
    mode: "core",
    subtitle: "Curated learning rails",
    href: "/kids/core"
  },
  {
    profile_id: "child-teen-01",
    name: "Mika Studio",
    age_band: "12-16",
    mode: "teen",
    subtitle: "Explore with active safety",
    href: "/kids/teen"
  }
];
