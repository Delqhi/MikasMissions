import type { KidsMode, KidsProgressResponse } from "./experience_types";
import { safeJSONFetch, apiBaseURL, useFallbackData } from "./fetch_helpers";
import { kidsProgressFallback } from "./mock_payloads";

export async function fetchKidsProgress(mode: KidsMode, childProfileID: string, token?: string): Promise<KidsProgressResponse> {
  const url = new URL(`/v1/kids/progress/${childProfileID}`, apiBaseURL());

  try {
    return await safeJSONFetch<KidsProgressResponse>(url.toString(), { token });
  } catch (error) {
    if (useFallbackData() || process.env.NODE_ENV === "production") {
      return kidsProgressFallback[mode];
    }
    throw error;
  }
}
