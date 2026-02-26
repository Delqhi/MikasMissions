import { cookies } from "next/headers";
import { redirect } from "next/navigation";

import { AdminShell } from "../../../components/layout/admin_shell";
import { apiBaseURL, safeJSONFetch } from "../../../lib/fetch_helpers";
import styles from "./page.module.css";

type RunResponse = {
  run_id: string;
  workflow_id: string;
  status: string;
  priority: string;
  auto_publish: boolean;
  input_payload: Record<string, unknown>;
  last_error: string;
};

type RunLogsResponse = {
  run_id: string;
  logs: Array<{
    run_id: string;
    step: string;
    status: string;
    message: string;
    event_time: string;
  }>;
};

type RunsPageProps = {
  searchParams: Promise<Record<string, string | string[] | undefined>>;
};

async function tokenFromCookie(): Promise<string> {
  const cookieStore = await cookies();
  return cookieStore.get("mm_admin_access_token")?.value ?? "";
}

async function retryRun(formData: FormData) {
  "use server";
  const token = await tokenFromCookie();
  if (token === "") {
    redirect("/admin/login");
  }
  const runID = String(formData.get("run_id") ?? "");
  await safeJSONFetch(new URL(`/v1/admin/runs/${encodeURIComponent(runID)}/retry`, apiBaseURL()).toString(), {
    method: "POST",
    token
  });
  redirect(`/admin/runs?run_id=${encodeURIComponent(runID)}`);
}

async function cancelRun(formData: FormData) {
  "use server";
  const token = await tokenFromCookie();
  if (token === "") {
    redirect("/admin/login");
  }
  const runID = String(formData.get("run_id") ?? "");
  await safeJSONFetch(new URL(`/v1/admin/runs/${encodeURIComponent(runID)}/cancel`, apiBaseURL()).toString(), {
    method: "POST",
    token
  });
  redirect(`/admin/runs?run_id=${encodeURIComponent(runID)}`);
}

export const dynamic = "force-dynamic";

export default async function AdminRunsPage({ searchParams }: RunsPageProps) {
  const token = await tokenFromCookie();
  if (token === "") {
    return (
      <AdminShell
        activeNav="Runs"
        subtitle="Authenticate to inspect execution logs and run controls."
        title="Run operations"
      >
        <section className={styles.notice}>
          <h2>Admin session missing</h2>
          <a href="/admin/login">Open admin login</a>
        </section>
      </AdminShell>
    );
  }

  const params = await searchParams;
  const runID = typeof params.run_id === "string" ? params.run_id : "";
  if (runID === "") {
    return (
      <AdminShell
        activeNav="Runs"
        subtitle="Select a run id from Studio to inspect status and step logs."
        title="Run operations"
      >
        <section className={styles.notice}>
          <h2>No run selected</h2>
          <a href="/admin/studio">Back to studio</a>
        </section>
      </AdminShell>
    );
  }

  const [run, logs] = await Promise.all([
    safeJSONFetch<RunResponse>(new URL(`/v1/admin/runs/${encodeURIComponent(runID)}`, apiBaseURL()).toString(), { token }),
    safeJSONFetch<RunLogsResponse>(new URL(`/v1/admin/runs/${encodeURIComponent(runID)}/logs`, apiBaseURL()).toString(), {
      token
    })
  ]);

  return (
    <AdminShell
      activeNav="Runs"
      subtitle="Inspect each step, retry failures, or cancel active jobs with full status context."
      title={`Run ${run.run_id}`}
    >
      <section className={styles.summary}>
        <article className={styles.metaCard}>
          <span>Status</span>
          <strong className={styles.statusPill} data-status={run.status.toLowerCase()}>
            {run.status}
          </strong>
        </article>

        <article className={styles.metaCard}>
          <span>Workflow</span>
          <strong>{run.workflow_id}</strong>
        </article>

        <article className={styles.metaCard}>
          <span>Priority</span>
          <strong>{run.priority}</strong>
        </article>
      </section>

      <section className={styles.panel}>
        <h2>Run controls</h2>
        <p>Last error: {run.last_error || "none"}</p>
        <div className={styles.actions}>
          <form action={retryRun}>
            <input name="run_id" type="hidden" value={run.run_id} />
            <button type="submit">Retry run</button>
          </form>
          <form action={cancelRun}>
            <input name="run_id" type="hidden" value={run.run_id} />
            <button type="submit">Cancel run</button>
          </form>
          <a href="/admin/studio">Back to studio</a>
        </div>
      </section>

      <section className={styles.panel}>
        <h2>Run logs</h2>
        <ul className={styles.logList}>
          {logs.logs.length === 0 ? (
            <li className={styles.empty}>No logs available for this run yet.</li>
          ) : (
            logs.logs.map((entry, index) => (
              <li key={`${entry.event_time}-${index}`}>
                <div className={styles.logMeta}>
                  <span>{entry.event_time}</span>
                  <strong>{entry.step}</strong>
                  <em data-status={entry.status.toLowerCase()}>{entry.status}</em>
                </div>
                <p>{entry.message}</p>
              </li>
            ))
          )}
        </ul>
      </section>
    </AdminShell>
  );
}
