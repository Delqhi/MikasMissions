import type { KidsHomeResponse, KidsMode } from "./experience_types";
import { safeJSONFetch, apiBaseURL, useFallbackData } from "./fetch_helpers";
import { kidsHomeFallback } from "./mock_payloads";

export async function fetchKidsHome(mode: KidsMode, childProfileID: string, token?: string): Promise<KidsHomeResponse> {
  const url = new URL("/v1/kids/home", apiBaseURL());
  url.searchParams.set("child_profile_id", childProfileID);
  url.searchParams.set("mode", mode);

  try {
    return await safeJSONFetch<KidsHomeResponse>(url.toString(), { token });
  } catch (error) {
    if (useFallbackData()) {
      return kidsHomeFallback[mode];
    }
    throw error;
  }
}
