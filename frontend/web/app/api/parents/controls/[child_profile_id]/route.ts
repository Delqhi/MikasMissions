import { cookies } from "next/headers";
import { NextResponse } from "next/server";

import { apiBaseURL } from "../../../../../lib/fetch_helpers";

type Params = {
  child_profile_id: string;
};

export async function PUT(request: Request, context: { params: Promise<Params> }) {
  const { child_profile_id: childProfileID } = await context.params;
  const cookieStore = await cookies();
  const accessToken = cookieStore.get("mm_access_token")?.value;
  const parentUserID = cookieStore.get("mm_parent_user_id")?.value;

  if (!accessToken || !parentUserID) {
    return NextResponse.json(
      {
        code: "missing_session",
        message: "parent session is required"
      },
      { status: 401 }
    );
  }

  const target = new URL(`/v1/parents/controls/${childProfileID}`, apiBaseURL());
  target.searchParams.set("parent_user_id", parentUserID);

  const response = await fetch(target.toString(), {
    method: "PUT",
    headers: {
      Authorization: `Bearer ${accessToken}`,
      "Content-Type": "application/json"
    },
    body: await request.text()
  });

  const body = await response.text();
  return new NextResponse(body, {
    status: response.status,
    headers: {
      "Content-Type": response.headers.get("content-type") ?? "application/json"
    }
  });
}
