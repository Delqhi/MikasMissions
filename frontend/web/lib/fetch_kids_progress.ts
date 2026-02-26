import type { KidsMode, KidsProgressResponse } from "./experience_types";
import { safeJSONFetch, apiBaseURL, useFallbackData } from "./fetch_helpers";
import { kidsProgressFallback } from "./mock_payloads";

function isRecord(input: unknown): input is Record<string, unknown> {
  return typeof input === "object" && input !== null;
}

function toNumber(input: unknown, fallback: number): number {
  const parsed = Number(input);
  return Number.isFinite(parsed) ? parsed : fallback;
}

function normalizeKidsProgressResponse(payload: unknown, mode: KidsMode, childProfileID: string): KidsProgressResponse {
  const fallback = kidsProgressFallback[mode];
  const input = isRecord(payload) ? payload : {};
  const responseChildID =
    typeof input.child_profile_id === "string" && input.child_profile_id.length > 0 ? input.child_profile_id : childProfileID;

  return {
    child_profile_id: responseChildID,
    watched_minutes_today: Math.max(0, toNumber(input.watched_minutes_today, fallback.watched_minutes_today)),
    watched_minutes_7d: Math.max(0, toNumber(input.watched_minutes_7d, fallback.watched_minutes_7d)),
    completion_percent: Math.max(0, Math.min(100, toNumber(input.completion_percent, fallback.completion_percent))),
    mission_streak_days: Math.max(0, Math.round(toNumber(input.mission_streak_days, fallback.mission_streak_days))),
    session_limit_minutes: Math.max(0, toNumber(input.session_limit_minutes, fallback.session_limit_minutes)),
    session_minutes_used: Math.max(0, toNumber(input.session_minutes_used, fallback.session_minutes_used)),
    session_capped: typeof input.session_capped === "boolean" ? input.session_capped : fallback.session_capped,
    last_episode_id: typeof input.last_episode_id === "string" ? input.last_episode_id : fallback.last_episode_id
  };
}

export async function fetchKidsProgress(mode: KidsMode, childProfileID: string, token?: string): Promise<KidsProgressResponse> {
  const url = new URL(`/v1/kids/progress/${childProfileID}`, apiBaseURL());

  try {
    const payload = await safeJSONFetch<unknown>(url.toString(), { token });
    return normalizeKidsProgressResponse(payload, mode, childProfileID);
  } catch (error) {
    if (useFallbackData() || process.env.NODE_ENV === "production") {
      return kidsProgressFallback[mode];
    }
    throw error;
  }
}
