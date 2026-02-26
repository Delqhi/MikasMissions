import { cookies } from "next/headers";
import { redirect } from "next/navigation";

import { AdminShell } from "../../../components/layout/admin_shell";
import { apiBaseURL, safeJSONFetch } from "../../../lib/fetch_helpers";
import { getLocaleAndMessages, getLocaleFromRequest, withLocalePath } from "../../../lib/i18n";
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

  const [token, locale] = await Promise.all([tokenFromCookie(), getLocaleFromRequest()]);
  if (token === "") {
    redirect(withLocalePath(locale, "/admin/login"));
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

  redirect(withLocalePath(locale, "/admin/studio"));
}

async function updateModelProfile(formData: FormData) {
  "use server";

  const [token, locale] = await Promise.all([tokenFromCookie(), getLocaleFromRequest()]);
  if (token === "") {
    redirect(withLocalePath(locale, "/admin/login"));
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

  redirect(withLocalePath(locale, "/admin/studio"));
}

async function startRun(formData: FormData) {
  "use server";

  const [token, locale] = await Promise.all([tokenFromCookie(), getLocaleFromRequest()]);
  if (token === "") {
    redirect(withLocalePath(locale, "/admin/login"));
  }

  const workflowID = String(formData.get("workflow_id") ?? "");
  if (workflowID === "") {
    redirect(withLocalePath(locale, "/admin/studio"));
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

  redirect(withLocalePath(locale, `/admin/runs?run_id=${encodeURIComponent(run.run_id)}`));
}

export const dynamic = "force-dynamic";

export default async function AdminStudioPage() {
  const [{ locale, messages }, token] = await Promise.all([getLocaleAndMessages(), tokenFromCookie()]);

  if (token === "") {
    return (
      <AdminShell
        activeNav="studio"
        labels={messages.admin.shell}
        locale={locale}
        subtitle={messages.admin.studio.noSessionText}
        title={messages.admin.studio.title}
      >
        <section className={styles.notice}>
          <h2>{messages.admin.studio.noSessionTitle}</h2>
          <p>{messages.admin.studio.noSessionText}</p>
          <a href={withLocalePath(locale, "/admin/login")}>{messages.admin.studio.openLogin}</a>
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
      activeNav="studio"
      labels={messages.admin.shell}
      locale={locale}
      subtitle={messages.admin.studio.subtitle}
      title={messages.admin.studio.title}
    >
      <section className={styles.hero}>
        <article>
          <h2>
            {workflows.workflows.length} {messages.admin.studio.workflowTemplates}
          </h2>
          <p>{messages.admin.studio.templatesText}</p>
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
          <h2>{messages.admin.studio.modelProfile}</h2>
          <form action={updateModelProfile} className={styles.form}>
            <label>
              <span>{messages.admin.studio.labelProfileId}</span>
              <input defaultValue={modelProfile.model_profile_id} name="model_profile_id" readOnly />
            </label>

            <label>
              <span>{messages.admin.studio.labelBaseUrl}</span>
              <input defaultValue={modelProfile.base_url} name="base_url" placeholder="NIM base url" required />
            </label>

            <label>
              <span>{messages.admin.studio.labelModelId}</span>
              <input defaultValue={modelProfile.model_id} name="model_id" placeholder="model id" required />
            </label>

            <label>
              <span>{messages.admin.studio.labelTimeout}</span>
              <input defaultValue={modelProfile.timeout_ms} min={500} name="timeout_ms" type="number" />
            </label>

            <label>
              <span>{messages.admin.studio.labelMaxRetries}</span>
              <input defaultValue={modelProfile.max_retries} min={0} name="max_retries" type="number" />
            </label>

            <label>
              <span>{messages.admin.studio.labelSafetyPreset}</span>
              <input defaultValue={modelProfile.safety_preset} name="safety_preset" required />
            </label>

            <button type="submit">{messages.admin.studio.updateModelProfile}</button>
          </form>
        </article>

        <article className={styles.panel}>
          <h2>{messages.admin.studio.createWorkflow}</h2>
          <form action={createWorkflow} className={styles.form}>
            <label>
              <span>{messages.admin.studio.labelName}</span>
              <input name="name" placeholder={messages.admin.studio.labelName} required />
            </label>

            <label>
              <span>{messages.admin.studio.labelDescription}</span>
              <input name="description" placeholder={messages.admin.studio.labelDescription} required />
            </label>

            <label>
              <span>{messages.admin.studio.labelContentSuitability}</span>
              <select defaultValue="core" name="content_suitability">
                <option value="early">early</option>
                <option value="core">core</option>
                <option value="teen">teen</option>
              </select>
            </label>

            <label>
              <span>{messages.admin.studio.labelAgeBand}</span>
              <select defaultValue="6-11" name="age_band">
                <option value="3-5">3-5</option>
                <option value="6-11">6-11</option>
                <option value="12-16">12-16</option>
              </select>
            </label>

            <label>
              <span>{messages.admin.studio.labelPipelineSteps}</span>
              <input defaultValue="prompt,generate,qc" name="steps" required />
            </label>

            <label>
              <span>{messages.admin.studio.labelModelProfileId}</span>
              <input defaultValue="nim-default" name="model_profile_id" required />
            </label>

            <label>
              <span>{messages.admin.studio.labelSafetyProfile}</span>
              <input defaultValue="strict" name="safety_profile" required />
            </label>

            <button type="submit">{messages.admin.studio.createWorkflowButton}</button>
          </form>
        </article>
      </section>

      <section className={styles.panel}>
        <h2>{messages.admin.studio.workflowTemplates}</h2>
        <ul className={styles.list}>
          {workflows.workflows.length === 0 ? (
            <li className={styles.empty}>{messages.admin.studio.noTemplates}</li>
          ) : (
            workflows.workflows.map((workflow) => (
              <li key={workflow.workflow_id}>
                <div className={styles.workflowMeta}>
                  <strong>{workflow.name}</strong>
                  <p>{workflow.description}</p>
                  <span>
                    {workflow.age_band} · {workflow.content_suitability} · v{workflow.version} · {workflow.model_profile_id}
                  </span>
                  <em>
                    {messages.admin.studio.stepsPrefix}: {workflow.steps.join(" → ")}
                  </em>
                </div>
                <form action={startRun} className={styles.inlineForm}>
                  <input name="workflow_id" type="hidden" value={workflow.workflow_id} />
                  <label>
                    <span>{messages.admin.studio.labelTitleHint}</span>
                    <input defaultValue="Mika Mission" name="title_hint" placeholder="Title hint" required />
                  </label>
                  <label>
                    <span>{messages.admin.studio.labelTopic}</span>
                    <input defaultValue="teamwork" name="topic" placeholder="Topic" required />
                  </label>
                  <button type="submit">{messages.admin.studio.startRun}</button>
                </form>
              </li>
            ))
          )}
        </ul>
      </section>
    </AdminShell>
  );
}
