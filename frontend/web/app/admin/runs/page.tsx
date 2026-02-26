import { cookies } from "next/headers";
import { redirect } from "next/navigation";

import { AdminShell } from "../../../components/layout/admin_shell";
import { apiBaseURL, safeJSONFetch } from "../../../lib/fetch_helpers";
import { getLocaleAndMessages, getLocaleFromRequest, withLocalePath } from "../../../lib/i18n";
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
  const [token, locale] = await Promise.all([tokenFromCookie(), getLocaleFromRequest()]);
  if (token === "") {
    redirect(withLocalePath(locale, "/admin/login"));
  }
  const runID = String(formData.get("run_id") ?? "");
  await safeJSONFetch(new URL(`/v1/admin/runs/${encodeURIComponent(runID)}/retry`, apiBaseURL()).toString(), {
    method: "POST",
    token
  });
  redirect(withLocalePath(locale, `/admin/runs?run_id=${encodeURIComponent(runID)}`));
}

async function cancelRun(formData: FormData) {
  "use server";
  const [token, locale] = await Promise.all([tokenFromCookie(), getLocaleFromRequest()]);
  if (token === "") {
    redirect(withLocalePath(locale, "/admin/login"));
  }
  const runID = String(formData.get("run_id") ?? "");
  await safeJSONFetch(new URL(`/v1/admin/runs/${encodeURIComponent(runID)}/cancel`, apiBaseURL()).toString(), {
    method: "POST",
    token
  });
  redirect(withLocalePath(locale, `/admin/runs?run_id=${encodeURIComponent(runID)}`));
}

export const dynamic = "force-dynamic";

export default async function AdminRunsPage({ searchParams }: RunsPageProps) {
  const [{ locale, messages }, token, params] = await Promise.all([getLocaleAndMessages(), tokenFromCookie(), searchParams]);

  if (token === "") {
    return (
      <AdminShell
        activeNav="runs"
        labels={messages.admin.shell}
        locale={locale}
        subtitle={messages.admin.runs.subtitle}
        title={messages.admin.runs.title}
      >
        <section className={styles.notice}>
          <h2>{messages.admin.runs.noSessionTitle}</h2>
          <a href={withLocalePath(locale, "/admin/login")}>{messages.admin.runs.openLogin}</a>
        </section>
      </AdminShell>
    );
  }

  const runID = typeof params.run_id === "string" ? params.run_id : "";
  if (runID === "") {
    return (
      <AdminShell
        activeNav="runs"
        labels={messages.admin.shell}
        locale={locale}
        subtitle={messages.admin.runs.subtitle}
        title={messages.admin.runs.title}
      >
        <section className={styles.notice}>
          <h2>{messages.admin.runs.noRunTitle}</h2>
          <a href={withLocalePath(locale, "/admin/studio")}>{messages.admin.runs.backToStudio}</a>
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
      activeNav="runs"
      labels={messages.admin.shell}
      locale={locale}
      subtitle={messages.admin.runs.subtitle}
      title={`${messages.admin.runs.title} · ${run.run_id}`}
    >
      <section className={styles.summary}>
        <article className={styles.metaCard}>
          <span>{messages.admin.runs.status}</span>
          <strong className={styles.statusPill} data-status={run.status.toLowerCase()}>
            {run.status}
          </strong>
        </article>

        <article className={styles.metaCard}>
          <span>{messages.admin.runs.workflow}</span>
          <strong>{run.workflow_id}</strong>
        </article>

        <article className={styles.metaCard}>
          <span>{messages.admin.runs.priority}</span>
          <strong>{run.priority}</strong>
        </article>
      </section>

      <section className={styles.panel}>
        <h2>{messages.admin.runs.runControls}</h2>
        <p>
          {messages.admin.runs.lastError}: {run.last_error || messages.admin.runs.none}
        </p>
        <div className={styles.actions}>
          <form action={retryRun}>
            <input name="run_id" type="hidden" value={run.run_id} />
            <button type="submit">{messages.admin.runs.retryRun}</button>
          </form>
          <form action={cancelRun}>
            <input name="run_id" type="hidden" value={run.run_id} />
            <button type="submit">{messages.admin.runs.cancelRun}</button>
          </form>
          <a href={withLocalePath(locale, "/admin/studio")}>{messages.admin.runs.backToStudio}</a>
        </div>
      </section>

      <section className={styles.panel}>
        <h2>{messages.admin.runs.runLogs}</h2>
        <ul className={styles.logList}>
          {logs.logs.length === 0 ? (
            <li className={styles.empty}>{messages.admin.runs.noLogs}</li>
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
