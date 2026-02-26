import { cookies } from "next/headers";

export async function accessTokenFromCookie(): Promise<string> {
  const cookieStore = await cookies();
  return cookieStore.get("mm_access_token")?.value ?? "";
}

export async function parentUserIDFromCookie(): Promise<string> {
  const cookieStore = await cookies();
  return cookieStore.get("mm_parent_user_id")?.value ?? "";
}
