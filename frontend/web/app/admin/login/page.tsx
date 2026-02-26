import { cookies } from "next/headers";
import { redirect } from "next/navigation";

import { AdminShell } from "../../../components/layout/admin_shell";
import { apiBaseURL, safeJSONFetch } from "../../../lib/fetch_helpers";
import { getLocaleAndMessages, getLocaleFromRequest, withLocalePath } from "../../../lib/i18n";
import styles from "./page.module.css";

type AdminLoginResponse = {
  access_token: string;
  expires_in: number;
  admin_user_id: string;
};

async function loginAdmin(formData: FormData) {
  "use server";

  const locale = await getLocaleFromRequest();
  const email = String(formData.get("email") ?? "").trim().toLowerCase();
  const password = String(formData.get("password") ?? "");
  const baseURL = apiBaseURL();

  const response = await safeJSONFetch<AdminLoginResponse>(new URL("/v1/admin/login", baseURL).toString(), {
    method: "POST",
    body: { email, password }
  });

  const cookieStore = await cookies();
  cookieStore.set("mm_admin_access_token", response.access_token, {
    httpOnly: true,
    maxAge: response.expires_in,
    path: "/",
    sameSite: "lax",
    secure: process.env.NODE_ENV === "production"
  });
  cookieStore.set("mm_admin_user_id", response.admin_user_id, {
    httpOnly: true,
    maxAge: response.expires_in,
    path: "/",
    sameSite: "lax",
    secure: process.env.NODE_ENV === "production"
  });

  redirect(withLocalePath(locale, "/admin/studio"));
}

export default async function AdminLoginPage() {
  const { locale, messages } = await getLocaleAndMessages();

  return (
    <AdminShell
      activeNav="login"
      labels={messages.admin.shell}
      locale={locale}
      subtitle={messages.admin.login.subtitle}
      title={messages.admin.login.title}
    >
      <section className={styles.card}>
        <h2>{messages.admin.login.cardTitle}</h2>
        <p>{messages.admin.login.cardText}</p>

        <form action={loginAdmin} className={styles.form}>
          <label>
            <span>{messages.admin.login.emailLabel}</span>
            <input name="email" required type="email" />
          </label>

          <label>
            <span>{messages.admin.login.passwordLabel}</span>
            <input minLength={10} name="password" required type="password" />
          </label>

          <button type="submit">{messages.admin.login.submit}</button>
        </form>
      </section>
    </AdminShell>
  );
}
