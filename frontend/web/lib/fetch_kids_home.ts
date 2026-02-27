import type { KidsHomeResponse, KidsMode } from "./experience_types";
import { safeJSONFetch, apiBaseURL, useFallbackData } from "./fetch_helpers";
import { kidsHomeFallback } from "./mock_payloads";

function isRecord(input: unknown): input is Record<string, unknown> {
  return typeof input === "object" && input !== null;
}

function isAgeBand(input: unknown): input is "3-5" | "6-11" | "12-16" {
  return input === "3-5" || input === "6-11" || input === "12-16";
}

function isSuitability(input: unknown): input is "early" | "core" | "teen" {
  return input === "early" || input === "core" || input === "teen";
}

function isSafetyMode(input: unknown): input is "strict" | "balanced" {
  return input === "strict" || input === "balanced";
}

function normalizeRailItems(input: unknown, mode: KidsMode): KidsHomeResponse["rails"] {
  const fallbackRails = kidsHomeFallback[mode].rails;
  const rails = Array.isArray(input) ? input : [];

  if (rails.length === 0) {
    return fallbackRails;
  }

  return rails.map((entry, index) => {
    const fallback = fallbackRails[index % fallbackRails.length];
    const item = isRecord(entry) ? entry : {};
    const learningTags = Array.isArray(item.learning_tags)
      ? item.learning_tags.filter((tag): tag is string => typeof tag === "string")
      : fallback.learning_tags;
    const durationMS = Number(item.duration_ms);
    const ageFitScore = Number(item.age_fit_score);

    return {
      episode_id:
        typeof item.episode_id === "string" && item.episode_id.length > 0
          ? item.episode_id
          : `${fallback.episode_id}-${index + 1}`,
      title: typeof item.title === "string" && item.title.length > 0 ? item.title : fallback.title,
      summary: typeof item.summary === "string" ? item.summary : fallback.summary,
      thumbnail_url:
        typeof item.thumbnail_url === "string" && item.thumbnail_url.length > 0 ? item.thumbnail_url : fallback.thumbnail_url,
      duration_ms: Number.isFinite(durationMS) ? Math.max(0, durationMS) : fallback.duration_ms,
      age_band: isAgeBand(item.age_band) ? item.age_band : fallback.age_band,
      content_suitability: isSuitability(item.content_suitability) ? item.content_suitability : fallback.content_suitability,
      learning_tags: learningTags.length > 0 ? learningTags : fallback.learning_tags,
      reason_code: typeof item.reason_code === "string" && item.reason_code.length > 0 ? item.reason_code : fallback.reason_code,
      safety_applied: typeof item.safety_applied === "boolean" ? item.safety_applied : fallback.safety_applied,
      age_fit_score: Number.isFinite(ageFitScore) ? ageFitScore : fallback.age_fit_score
    };
  });
}

function normalizePrimaryActions(input: unknown, mode: KidsMode): string[] {
  if (!Array.isArray(input)) {
    return kidsHomeFallback[mode].primary_actions;
  }

  const actions = input.filter((action): action is string => typeof action === "string");
  return actions.length > 0 ? actions : kidsHomeFallback[mode].primary_actions;
}

function normalizeKidsHomeResponse(payload: unknown, mode: KidsMode, childProfileID: string): KidsHomeResponse {
  const fallback = kidsHomeFallback[mode];
  const input = isRecord(payload) ? payload : {};
  const responseChildID =
    typeof input.child_profile_id === "string" && input.child_profile_id.length > 0 ? input.child_profile_id : childProfileID;

  return {
    child_profile_id: responseChildID,
    mode,
    safety_mode: isSafetyMode(input.safety_mode) ? input.safety_mode : fallback.safety_mode,
    primary_actions: normalizePrimaryActions(input.primary_actions, mode),
    rails: normalizeRailItems(input.rails, mode)
  };
}

export async function fetchKidsHome(mode: KidsMode, childProfileID: string, token?: string): Promise<KidsHomeResponse> {
  const url = new URL("/v1/kids/home", apiBaseURL());
  url.searchParams.set("child_profile_id", childProfileID);
  url.searchParams.set("mode", mode);

  try {
    const payload = await safeJSONFetch<unknown>(url.toString(), { token });
    return normalizeKidsHomeResponse(payload, mode, childProfileID);
  } catch (error) {
    if (useFallbackData() || process.env.NODE_ENV === "production") {
      return kidsHomeFallback[mode];
    }
    throw error;
  }
}
