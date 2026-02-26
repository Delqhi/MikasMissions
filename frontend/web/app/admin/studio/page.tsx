import { cookies } from "next/headers";
import { redirect } from "next/navigation";

import { AdminShell } from "../../../components/layout/admin_shell";
import { apiBaseURL, safeJSONFetch } from "../../../lib/fetch_helpers";
import styles from "./page.module.css";

type AdminWorkflow = {
  workflow_id: string;
  name: string;
  description: string;
  content_suitability: "early" | "core" | "teen";
  age_band: "3-5" | "6-11" | "12-16";
  steps: string[];
  model_profile_id: string;
  safety_profile: string;
  version: number;
};

type WorkflowsResponse = {
  workflows: AdminWorkflow[];
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

type RunResponse = {
  run_id: string;
  status: string;
};

async function tokenFromCookie(): Promise<string> {
  const cookieStore = await cookies();
  return cookieStore.get("mm_admin_access_token")?.value ?? "";
}

async function createWorkflow(formData: FormData) {
  "use server";

  const token = await tokenFromCookie();
  if (token === "") {
    redirect("/admin/login");
  }

  const body = {
    name: String(formData.get("name") ?? "").trim(),
    description: String(formData.get("description") ?? "").trim(),
    content_suitability: String(formData.get("content_suitability") ?? "core"),
    age_band: String(formData.get("age_band") ?? "6-11"),
    steps: String(formData.get("steps") ?? "prompt,generate,qc")
      .split(",")
      .map((item) => item.trim())
      .filter(Boolean),
    model_profile_id: String(formData.get("model_profile_id") ?? "nim-default"),
    safety_profile: String(formData.get("safety_profile") ?? "strict")
  };

  await safeJSONFetch(new URL("/v1/admin/workflows", apiBaseURL()).toString(), {
    method: "POST",
    body,
    token
  });

  redirect("/admin/studio");
}

async function updateModelProfile(formData: FormData) {
  "use server";

  const token = await tokenFromCookie();
  if (token === "") {
    redirect("/admin/login");
  }

  const profileID = String(formData.get("model_profile_id") ?? "nim-default");
  const body = {
    model_profile_id: profileID,
    provider: "nvidia_nim",
    base_url: String(formData.get("base_url") ?? "http://127.0.0.1:9000"),
    model_id: String(formData.get("model_id") ?? "nim-video-v1"),
    timeout_ms: Number(formData.get("timeout_ms") ?? 12000),
    max_retries: Number(formData.get("max_retries") ?? 2),
    safety_preset: String(formData.get("safety_preset") ?? "kids_strict")
  };

  await safeJSONFetch(new URL(`/v1/admin/model-profiles/${encodeURIComponent(profileID)}`, apiBaseURL()).toString(), {
    method: "PUT",
    body,
    token
  });

  redirect("/admin/studio");
}

async function startRun(formData: FormData) {
  "use server";

  const token = await tokenFromCookie();
  if (token === "") {
    redirect("/admin/login");
  }

  const workflowID = String(formData.get("workflow_id") ?? "");
  if (workflowID === "") {
    redirect("/admin/studio");
  }

  const run = await safeJSONFetch<RunResponse>(
    new URL(`/v1/admin/workflows/${encodeURIComponent(workflowID)}/runs`, apiBaseURL()).toString(),
    {
      method: "POST",
      body: {
        input_payload: {
          title_hint: String(formData.get("title_hint") ?? "Mission"),
          topic: String(formData.get("topic") ?? "science")
        },
        priority: String(formData.get("priority") ?? "normal"),
        auto_publish: false
      },
      token
    }
  );

  redirect(`/admin/runs?run_id=${encodeURIComponent(run.run_id)}`);
}

export const dynamic = "force-dynamic";

export default async function AdminStudioPage() {
  const token = await tokenFromCookie();
  if (token === "") {
    return (
      <AdminShell
        activeNav="Studio"
        subtitle="Authenticate first to configure workflows and start generation runs."
        title="Admin studio"
      >
        <section className={styles.notice}>
          <h2>Admin session missing</h2>
          <p>Login as admin to access workflow operations.</p>
          <a href="/admin/login">Open admin login</a>
        </section>
      </AdminShell>
    );
  }

  const [workflows, modelProfile] = await Promise.all([
    safeJSONFetch<WorkflowsResponse>(new URL("/v1/admin/workflows", apiBaseURL()).toString(), { token }),
    safeJSONFetch<ModelProfile>(new URL("/v1/admin/model-profiles/nim-default", apiBaseURL()).toString(), { token })
  ]);

  return (
    <AdminShell
      activeNav="Studio"
      subtitle="Configure model profiles, maintain workflow templates, and launch controlled generation runs."
      title="Admin studio"
    >
      <section className={styles.hero}>
        <article>
          <h2>{workflows.workflows.length} workflow templates</h2>
          <p>Templates are versioned and tied to age-band and safety suitability.</p>
        </article>
        <article>
          <h2>{modelProfile.model_profile_id}</h2>
          <p>
            Provider {modelProfile.provider} · timeout {modelProfile.timeout_ms}ms · retries {modelProfile.max_retries}
          </p>
        </article>
      </section>

      <section className={styles.grid}>
        <article className={styles.panel}>
          <h2>Model profile</h2>
          <form action={updateModelProfile} className={styles.form}>
            <label>
              <span>Profile ID</span>
              <input defaultValue={modelProfile.model_profile_id} name="model_profile_id" readOnly />
            </label>

            <label>
              <span>Base URL</span>
              <input defaultValue={modelProfile.base_url} name="base_url" placeholder="NIM base url" required />
            </label>

            <label>
              <span>Model ID</span>
              <input defaultValue={modelProfile.model_id} name="model_id" placeholder="model id" required />
            </label>

            <label>
              <span>Timeout (ms)</span>
              <input defaultValue={modelProfile.timeout_ms} min={500} name="timeout_ms" type="number" />
            </label>

            <label>
              <span>Max retries</span>
              <input defaultValue={modelProfile.max_retries} min={0} name="max_retries" type="number" />
            </label>

            <label>
              <span>Safety preset</span>
              <input defaultValue={modelProfile.safety_preset} name="safety_preset" required />
            </label>

            <button type="submit">Update model profile</button>
          </form>
        </article>

        <article className={styles.panel}>
          <h2>Create workflow</h2>
          <form action={createWorkflow} className={styles.form}>
            <label>
              <span>Name</span>
              <input name="name" placeholder="Name" required />
            </label>

            <label>
              <span>Beschreibung</span>
              <input name="description" placeholder="Beschreibung" required />
            </label>

            <label>
              <span>Content suitability</span>
              <select defaultValue="core" name="content_suitability">
                <option value="early">early</option>
                <option value="core">core</option>
                <option value="teen">teen</option>
              </select>
            </label>

            <label>
              <span>Age band</span>
              <select defaultValue="6-11" name="age_band">
                <option value="3-5">3-5</option>
                <option value="6-11">6-11</option>
                <option value="12-16">12-16</option>
              </select>
            </label>

            <label>
              <span>Pipeline steps</span>
              <input defaultValue="prompt,generate,qc" name="steps" required />
            </label>

            <label>
              <span>Model profile ID</span>
              <input defaultValue="nim-default" name="model_profile_id" required />
            </label>

            <label>
              <span>Safety profile</span>
              <input defaultValue="strict" name="safety_profile" required />
            </label>

            <button type="submit">Create workflow</button>
          </form>
        </article>
      </section>

      <section className={styles.panel}>
        <h2>Workflow templates</h2>
        <ul className={styles.list}>
          {workflows.workflows.length === 0 ? (
            <li className={styles.empty}>No templates yet. Create your first workflow above.</li>
          ) : (
            workflows.workflows.map((workflow) => (
              <li key={workflow.workflow_id}>
                <div className={styles.workflowMeta}>
                  <strong>{workflow.name}</strong>
                  <p>{workflow.description}</p>
                  <span>
                    {workflow.age_band} · {workflow.content_suitability} · v{workflow.version} · {workflow.model_profile_id}
                  </span>
                  <em>Steps: {workflow.steps.join(" → ")}</em>
                </div>
                <form action={startRun} className={styles.inlineForm}>
                  <input name="workflow_id" type="hidden" value={workflow.workflow_id} />
                  <label>
                    <span>Title hint</span>
                    <input defaultValue="Mika Mission" name="title_hint" placeholder="Title hint" required />
                  </label>
                  <label>
                    <span>Topic</span>
                    <input defaultValue="teamwork" name="topic" placeholder="Topic" required />
                  </label>
                  <button type="submit">Start run</button>
                </form>
              </li>
            ))
          )}
        </ul>
      </section>
    </AdminShell>
  );
}
