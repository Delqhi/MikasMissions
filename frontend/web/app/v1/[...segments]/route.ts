import { NextResponse } from "next/server";
import type { AgeBand, KidsMode, ParentalControls } from "../../../lib/experience_types";
import { kidsHomeFallback, kidsProgressFallback, parentControlsFallback } from "../../../lib/mock_payloads";

type RouteContext = {
  params: Promise<{ segments: string[] }>;
};

type ChildProfile = {
  age_band: AgeBand;
  avatar: string;
  child_profile_id: string;
  display_name: string;
  parent_user_id: string;
};

type ParentUser = {
  email: string;
  parent_user_id: string;
  password: string;
};

type AdminUser = {
  admin_user_id: string;
  email: string;
  password: string;
};

type AdminWorkflow = {
  workflow_id: string;
  name: string;
  description: string;
  content_suitability: "early" | "core" | "teen";
  age_band: AgeBand;
  steps: string[];
  model_profile_id: string;
  safety_profile: string;
  version: number;
};

type ModelProfile = {
  model_profile_id: string;
  provider: "nvidia_nim";
  base_url: string;
  model_id: string;
  timeout_ms: number;
  max_retries: number;
  safety_preset: string;
};

type RunRecord = {
  run_id: string;
  workflow_id: string;
  status: string;
  priority: string;
  auto_publish: boolean;
  input_payload: Record<string, unknown>;
  last_error: string;
};

type RunLog = {
  run_id: string;
  step: string;
  status: string;
  message: string;
  event_time: string;
};

const parentUsersByEmail = new Map<string, ParentUser>();
const parentProfilesByUser = new Map<string, ChildProfile[]>();
const parentControlsByChild = new Map<string, ParentalControls>();
const parentTokenToUser = new Map<string, string>();

const adminUsersByEmail = new Map<string, AdminUser>();
const adminTokenToUser = new Map<string, string>();
const workflowsByID = new Map<string, AdminWorkflow>();
const modelProfilesByID = new Map<string, ModelProfile>();
const runsByID = new Map<string, RunRecord>();
const runLogsByRunID = new Map<string, RunLog[]>();

function createID(prefix: string): string {
  return `${prefix}-${Math.random().toString(36).slice(2, 10)}`;
}

function nowISO(): string {
  return new Date().toISOString();
}

function json(data: unknown, status = 200): NextResponse {
  return NextResponse.json(data, { status });
}

function authTokenFromRequest(request: Request): string {
  const value = request.headers.get("authorization") ?? "";
  if (!value.toLowerCase().startsWith("bearer ")) {
    return "";
  }
  return value.slice(7).trim();
}

function modeForAgeBand(ageBand: AgeBand): KidsMode {
  switch (ageBand) {
    case "3-5":
      return "early";
    case "12-16":
      return "teen";
    default:
      return "core";
  }
}

function findProfileByID(childProfileID: string): ChildProfile | null {
  for (const profiles of parentProfilesByUser.values()) {
    const found = profiles.find((profile) => profile.child_profile_id === childProfileID);
    if (found) {
      return found;
    }
  }
  return null;
}

function ensureSeedData(): void {
  if (parentUsersByEmail.size === 0) {
    const user: ParentUser = {
      email: "parent@example.com",
      parent_user_id: "parent-demo-01",
      password: "ChangeMe-Parent-2026"
    };
    parentUsersByEmail.set(user.email, user);
    parentProfilesByUser.set(user.parent_user_id, [
      {
        age_band: "6-11",
        avatar: "robot",
        child_profile_id: "child-core-01",
        display_name: "Mika Explorer",
        parent_user_id: user.parent_user_id
      }
    ]);
    parentControlsByChild.set("child-core-01", parentControlsFallback.controls);
  }

  if (adminUsersByEmail.size === 0) {
    const admin: AdminUser = {
      admin_user_id: "admin-demo-01",
      email: "admin@example.com",
      password: "ChangeMe-Admin-2026"
    };
    adminUsersByEmail.set(admin.email, admin);
  }

  if (modelProfilesByID.size === 0) {
    modelProfilesByID.set("nim-default", {
      model_profile_id: "nim-default",
      provider: "nvidia_nim",
      base_url: "http://127.0.0.1:9000",
      model_id: "nim-video-v1",
      timeout_ms: 12000,
      max_retries: 2,
      safety_preset: "kids_strict"
    });
  }

  if (workflowsByID.size === 0) {
    workflowsByID.set("wf-demo-01", {
      workflow_id: "wf-demo-01",
      name: "Mika Mission Builder",
      description: "Demo workflow for safe kids mission generation.",
      content_suitability: "core",
      age_band: "6-11",
      steps: ["prompt", "generate", "qc", "publish"],
      model_profile_id: "nim-default",
      safety_profile: "strict",
      version: 1
    });
  }
}

async function bodyJSON<T extends Record<string, unknown>>(request: Request): Promise<T> {
  try {
    return (await request.json()) as T;
  } catch {
    return {} as T;
  }
}

function requireParent(request: Request, fallbackParentUserID = ""): string {
  const token = authTokenFromRequest(request);
  const fromToken = token ? parentTokenToUser.get(token) ?? "" : "";
  return fromToken || fallbackParentUserID;
}

function requireAdmin(request: Request): string {
  const token = authTokenFromRequest(request);
  if (!token) {
    return "";
  }
  return adminTokenToUser.get(token) ?? "";
}

function handleChildrenProfilesGET(request: Request): NextResponse {
  const url = new URL(request.url);
  const parentUserID = requireParent(request, url.searchParams.get("parent_user_id") ?? "");
  if (!parentUserID) {
    return json({ code: "missing_parent_user_id", message: "parent_user_id is required" }, 400);
  }
  return json({ profiles: parentProfilesByUser.get(parentUserID) ?? [] });
}

async function handleChildrenProfilesPOST(request: Request): Promise<NextResponse> {
  const body = await bodyJSON<{
    age_band?: AgeBand;
    avatar?: string;
    display_name?: string;
    parent_user_id?: string;
  }>(request);
  const parentUserID = requireParent(request, body.parent_user_id ?? "");
  if (!parentUserID) {
    return json({ code: "missing_parent_session", message: "parent session is required" }, 401);
  }

  const ageBand: AgeBand = body.age_band === "3-5" || body.age_band === "12-16" ? body.age_band : "6-11";
  const profile: ChildProfile = {
    age_band: ageBand,
    avatar: body.avatar ?? "robot",
    child_profile_id: createID("child"),
    display_name: body.display_name?.trim() || "Mika",
    parent_user_id: parentUserID
  };
  const current = parentProfilesByUser.get(parentUserID) ?? [];
  parentProfilesByUser.set(parentUserID, [...current, profile]);
  parentControlsByChild.set(profile.child_profile_id, {
    ...parentControlsFallback.controls,
    session_limit_minutes: ageBand === "3-5" ? 35 : ageBand === "12-16" ? 75 : 50
  });
  return json({ child_profile_id: profile.child_profile_id }, 201);
}

function handleKidsHomeGET(request: Request): NextResponse {
  const url = new URL(request.url);
  const modeInput = url.searchParams.get("mode");
  const mode: KidsMode = modeInput === "early" || modeInput === "teen" ? modeInput : "core";
  const childProfileID = url.searchParams.get("child_profile_id") ?? kidsHomeFallback[mode].child_profile_id;
  return json({ ...kidsHomeFallback[mode], child_profile_id: childProfileID });
}

function handleKidsProgressGET(childProfileID: string): NextResponse {
  const profile = findProfileByID(childProfileID);
  const mode = modeForAgeBand(profile?.age_band ?? "6-11");
  return json({ ...kidsProgressFallback[mode], child_profile_id: childProfileID });
}

function handleParentControlsGET(request: Request, childProfileID: string): NextResponse {
  const url = new URL(request.url);
  const parentUserID = requireParent(request, url.searchParams.get("parent_user_id") ?? "");
  if (!parentUserID) {
    return json({ code: "missing_parent_session", message: "parent session is required" }, 401);
  }
  const controls = parentControlsByChild.get(childProfileID) ?? parentControlsFallback.controls;
  parentControlsByChild.set(childProfileID, controls);
  return json({ child_profile_id: childProfileID, controls });
}

async function handleParentControlsPUT(request: Request, childProfileID: string): Promise<NextResponse> {
  const body = await bodyJSON<ParentalControls>(request);
  const controls: ParentalControls = {
    autoplay: Boolean(body.autoplay),
    chat_enabled: Boolean(body.chat_enabled),
    external_links: Boolean(body.external_links),
    session_limit_minutes:
      typeof body.session_limit_minutes === "number" && Number.isFinite(body.session_limit_minutes)
        ? body.session_limit_minutes
        : parentControlsFallback.controls.session_limit_minutes,
    bedtime_window: typeof body.bedtime_window === "string" ? body.bedtime_window : parentControlsFallback.controls.bedtime_window,
    safety_mode: body.safety_mode === "balanced" ? "balanced" : "strict"
  };
  parentControlsByChild.set(childProfileID, controls);
  return json({ child_profile_id: childProfileID, controls });
}

async function handleParentSignupPOST(request: Request): Promise<NextResponse> {
  const body = await bodyJSON<{ email?: string; password?: string }>(request);
  const email = (body.email ?? "").trim().toLowerCase();
  const password = body.password ?? "";
  if (!email || !password) {
    return json({ code: "invalid_payload", message: "email and password are required" }, 400);
  }
  const existing = parentUsersByEmail.get(email);
  if (existing) {
    return json({ parent_user_id: existing.parent_user_id });
  }
  const user: ParentUser = {
    email,
    parent_user_id: createID("parent"),
    password
  };
  parentUsersByEmail.set(email, user);
  parentProfilesByUser.set(user.parent_user_id, []);
  return json({ parent_user_id: user.parent_user_id }, 201);
}

async function handleParentLoginPOST(request: Request): Promise<NextResponse> {
  const body = await bodyJSON<{ email?: string; password?: string }>(request);
  const email = (body.email ?? "").trim().toLowerCase();
  const password = body.password ?? "";
  if (!email || !password) {
    return json({ code: "invalid_payload", message: "email and password are required" }, 400);
  }
  let user = parentUsersByEmail.get(email);
  if (!user) {
    user = {
      email,
      parent_user_id: createID("parent"),
      password
    };
    parentUsersByEmail.set(email, user);
    parentProfilesByUser.set(user.parent_user_id, []);
  }
  if (user.password !== password) {
    return json({ code: "invalid_credentials", message: "email or password is invalid" }, 401);
  }
  const token = createID("ptok");
  parentTokenToUser.set(token, user.parent_user_id);
  return json({ access_token: token, expires_in: 60 * 60 * 8, parent_user_id: user.parent_user_id });
}

async function handleParentConsentVerifyPOST(request: Request): Promise<NextResponse> {
  const body = await bodyJSON<{ parent_user_id?: string }>(request);
  const parentUserID = requireParent(request, body.parent_user_id ?? "");
  if (!parentUserID) {
    return json({ code: "missing_parent_session", message: "parent session is required" }, 401);
  }
  return json({ parent_user_id: parentUserID, verified: true });
}

async function handleAdminLoginPOST(request: Request): Promise<NextResponse> {
  const body = await bodyJSON<{ email?: string; password?: string }>(request);
  const email = (body.email ?? "").trim().toLowerCase();
  const password = body.password ?? "";
  if (!email || !password) {
    return json({ code: "invalid_payload", message: "email and password are required" }, 400);
  }
  let admin = adminUsersByEmail.get(email);
  if (!admin) {
    admin = { admin_user_id: createID("admin"), email, password };
    adminUsersByEmail.set(email, admin);
  }
  if (admin.password !== password) {
    return json({ code: "invalid_credentials", message: "email or password is invalid" }, 401);
  }
  const token = createID("atok");
  adminTokenToUser.set(token, admin.admin_user_id);
  return json({ access_token: token, expires_in: 60 * 60 * 8, admin_user_id: admin.admin_user_id });
}

function handleAdminWorkflowsGET(request: Request): NextResponse {
  if (!requireAdmin(request)) {
    return json({ code: "missing_admin_session", message: "admin session is required" }, 401);
  }
  return json({ workflows: [...workflowsByID.values()] });
}

async function handleAdminWorkflowsPOST(request: Request): Promise<NextResponse> {
  if (!requireAdmin(request)) {
    return json({ code: "missing_admin_session", message: "admin session is required" }, 401);
  }
  const body = await bodyJSON<{
    name?: string;
    description?: string;
    content_suitability?: "early" | "core" | "teen";
    age_band?: AgeBand;
    steps?: string[];
    model_profile_id?: string;
    safety_profile?: string;
  }>(request);
  const workflow: AdminWorkflow = {
    workflow_id: createID("wf"),
    name: body.name?.trim() || "Untitled Workflow",
    description: body.description?.trim() || "No description",
    content_suitability: body.content_suitability === "early" || body.content_suitability === "teen" ? body.content_suitability : "core",
    age_band: body.age_band === "3-5" || body.age_band === "12-16" ? body.age_band : "6-11",
    steps: Array.isArray(body.steps) && body.steps.length > 0 ? body.steps : ["prompt", "generate", "qc"],
    model_profile_id: body.model_profile_id?.trim() || "nim-default",
    safety_profile: body.safety_profile?.trim() || "strict",
    version: 1
  };
  workflowsByID.set(workflow.workflow_id, workflow);
  return json(workflow, 201);
}

function handleAdminModelProfileGET(request: Request, profileID: string): NextResponse {
  if (!requireAdmin(request)) {
    return json({ code: "missing_admin_session", message: "admin session is required" }, 401);
  }
  return json(modelProfilesByID.get(profileID) ?? modelProfilesByID.get("nim-default"));
}

async function handleAdminModelProfilePUT(request: Request, profileID: string): Promise<NextResponse> {
  if (!requireAdmin(request)) {
    return json({ code: "missing_admin_session", message: "admin session is required" }, 401);
  }
  const body = await bodyJSON<ModelProfile>(request);
  const profile: ModelProfile = {
    model_profile_id: profileID,
    provider: "nvidia_nim",
    base_url: typeof body.base_url === "string" ? body.base_url : "http://127.0.0.1:9000",
    model_id: typeof body.model_id === "string" ? body.model_id : "nim-video-v1",
    timeout_ms: typeof body.timeout_ms === "number" ? body.timeout_ms : 12000,
    max_retries: typeof body.max_retries === "number" ? body.max_retries : 2,
    safety_preset: typeof body.safety_preset === "string" ? body.safety_preset : "kids_strict"
  };
  modelProfilesByID.set(profileID, profile);
  return json(profile);
}

async function handleAdminStartRunPOST(request: Request, workflowID: string): Promise<NextResponse> {
  if (!requireAdmin(request)) {
    return json({ code: "missing_admin_session", message: "admin session is required" }, 401);
  }
  const body = await bodyJSON<{ input_payload?: Record<string, unknown>; priority?: string; auto_publish?: boolean }>(request);
  const runID = createID("run");
  const run: RunRecord = {
    run_id: runID,
    workflow_id: workflowID,
    status: "QUEUED",
    priority: body.priority?.trim() || "normal",
    auto_publish: Boolean(body.auto_publish),
    input_payload: body.input_payload ?? {},
    last_error: ""
  };
  runsByID.set(runID, run);
  runLogsByRunID.set(runID, [
    {
      run_id: runID,
      step: "orchestrator",
      status: "QUEUED",
      message: "Run accepted and queued.",
      event_time: nowISO()
    }
  ]);
  return json({ run_id: runID, status: run.status }, 201);
}

function handleAdminRunGET(request: Request, runID: string): NextResponse {
  if (!requireAdmin(request)) {
    return json({ code: "missing_admin_session", message: "admin session is required" }, 401);
  }
  const run = runsByID.get(runID);
  if (!run) {
    return json({ code: "not_found", message: "run not found" }, 404);
  }
  return json(run);
}

function handleAdminRunLogsGET(request: Request, runID: string): NextResponse {
  if (!requireAdmin(request)) {
    return json({ code: "missing_admin_session", message: "admin session is required" }, 401);
  }
  return json({ run_id: runID, logs: runLogsByRunID.get(runID) ?? [] });
}

function mutateRunStatus(request: Request, runID: string, status: string, message: string): NextResponse {
  if (!requireAdmin(request)) {
    return json({ code: "missing_admin_session", message: "admin session is required" }, 401);
  }
  const run = runsByID.get(runID);
  if (!run) {
    return json({ code: "not_found", message: "run not found" }, 404);
  }
  run.status = status;
  run.last_error = status === "CANCELLED" ? "Cancelled by admin request." : "";
  runsByID.set(runID, run);
  const logs = runLogsByRunID.get(runID) ?? [];
  logs.push({
    run_id: runID,
    step: "orchestrator",
    status,
    message,
    event_time: nowISO()
  });
  runLogsByRunID.set(runID, logs);
  return json({ run_id: runID, status: run.status });
}

export async function GET(request: Request, context: RouteContext): Promise<NextResponse> {
  ensureSeedData();
  const { segments } = await context.params;

  if (segments[0] === "children" && segments[1] === "profiles" && segments.length === 2) {
    return handleChildrenProfilesGET(request);
  }
  if (segments[0] === "kids" && segments[1] === "home" && segments.length === 2) {
    return handleKidsHomeGET(request);
  }
  if (segments[0] === "kids" && segments[1] === "progress" && segments[2] && segments.length === 3) {
    return handleKidsProgressGET(segments[2]);
  }
  if (segments[0] === "parents" && segments[1] === "controls" && segments[2] && segments.length === 3) {
    return handleParentControlsGET(request, segments[2]);
  }
  if (segments[0] === "admin" && segments[1] === "workflows" && segments.length === 2) {
    return handleAdminWorkflowsGET(request);
  }
  if (segments[0] === "admin" && segments[1] === "model-profiles" && segments[2] && segments.length === 3) {
    return handleAdminModelProfileGET(request, segments[2]);
  }
  if (segments[0] === "admin" && segments[1] === "runs" && segments[2] && segments.length === 3) {
    return handleAdminRunGET(request, segments[2]);
  }
  if (segments[0] === "admin" && segments[1] === "runs" && segments[2] && segments[3] === "logs" && segments.length === 4) {
    return handleAdminRunLogsGET(request, segments[2]);
  }

  return json({ code: "not_found", message: "route not found" }, 404);
}

export async function POST(request: Request, context: RouteContext): Promise<NextResponse> {
  ensureSeedData();
  const { segments } = await context.params;

  if (segments[0] === "parents" && segments[1] === "signup" && segments.length === 2) {
    return handleParentSignupPOST(request);
  }
  if (segments[0] === "parents" && segments[1] === "login" && segments.length === 2) {
    return handleParentLoginPOST(request);
  }
  if (segments[0] === "parents" && segments[1] === "consent" && segments[2] === "verify" && segments.length === 3) {
    return handleParentConsentVerifyPOST(request);
  }
  if (segments[0] === "children" && segments[1] === "profiles" && segments.length === 2) {
    return handleChildrenProfilesPOST(request);
  }
  if (segments[0] === "admin" && segments[1] === "login" && segments.length === 2) {
    return handleAdminLoginPOST(request);
  }
  if (segments[0] === "admin" && segments[1] === "workflows" && segments.length === 2) {
    return handleAdminWorkflowsPOST(request);
  }
  if (segments[0] === "admin" && segments[1] === "workflows" && segments[2] && segments[3] === "runs" && segments.length === 4) {
    return handleAdminStartRunPOST(request, segments[2]);
  }
  if (segments[0] === "admin" && segments[1] === "runs" && segments[2] && segments[3] === "retry" && segments.length === 4) {
    return mutateRunStatus(request, segments[2], "RETRYING", "Run was re-queued by admin.");
  }
  if (segments[0] === "admin" && segments[1] === "runs" && segments[2] && segments[3] === "cancel" && segments.length === 4) {
    return mutateRunStatus(request, segments[2], "CANCELLED", "Run cancelled by admin.");
  }

  return json({ code: "not_found", message: "route not found" }, 404);
}

export async function PUT(request: Request, context: RouteContext): Promise<NextResponse> {
  ensureSeedData();
  const { segments } = await context.params;

  if (segments[0] === "parents" && segments[1] === "controls" && segments[2] && segments.length === 3) {
    return handleParentControlsPUT(request, segments[2]);
  }
  if (segments[0] === "admin" && segments[1] === "model-profiles" && segments[2] && segments.length === 3) {
    return handleAdminModelProfilePUT(request, segments[2]);
  }

  return json({ code: "not_found", message: "route not found" }, 404);
}
