import { apiBaseURL, safeJSONFetch } from "./fetch_helpers";

export type ChildProfileSummary = {
  age_band: "3-5" | "6-11" | "12-16";
  avatar: string;
  child_profile_id: string;
  display_name: string;
  parent_user_id: string;
};

type ListChildProfilesResponse = {
  profiles: ChildProfileSummary[];
};

export async function fetchChildProfiles(parentUserID: string, token?: string): Promise<ChildProfileSummary[]> {
  const url = new URL("/v1/children/profiles", apiBaseURL());
  if (parentUserID) {
    url.searchParams.set("parent_user_id", parentUserID);
  }
  const response = await safeJSONFetch<ListChildProfilesResponse>(url.toString(), { token });
  return response.profiles;
}
