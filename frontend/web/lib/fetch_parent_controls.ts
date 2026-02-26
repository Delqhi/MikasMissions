import type { ParentControlsResponse } from "./experience_types";
import { safeJSONFetch, apiBaseURL, useFallbackData } from "./fetch_helpers";
import { parentControlsFallback } from "./mock_payloads";

export async function fetchParentControls(childProfileID: string, parentUserID: string, token?: string): Promise<ParentControlsResponse> {
  const url = new URL(`/v1/parents/controls/${childProfileID}`, apiBaseURL());
  url.searchParams.set("parent_user_id", parentUserID);

  try {
    return await safeJSONFetch<ParentControlsResponse>(url.toString(), { token });
  } catch (error) {
    if (useFallbackData()) {
      return parentControlsFallback;
    }
    throw error;
  }
}
