import { cookies } from "next/headers";
import { redirect } from "next/navigation";

import { AdminShell } from "../../../components/layout/admin_shell";
import { apiBaseURL, safeJSONFetch } from "../../../lib/fetch_helpers";
import styles from "./page.module.css";

type AdminLoginResponse = {
  access_token: string;
  expires_in: number;
  admin_user_id: string;
};

async function loginAdmin(formData: FormData) {
  "use server";

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

  redirect("/admin/studio");
}

export default function AdminLoginPage() {
  return (
    <AdminShell
      activeNav="Login"
      subtitle="Role-based entry for workflow orchestration and model profile management."
      title="Admin studio login"
    >
      <section className={styles.card}>
        <h2>Authenticate as admin</h2>
        <p>Use your admin credentials to access workflow templates, model controls, and run operations.</p>

        <form action={loginAdmin} className={styles.form}>
          <label>
            <span>Admin email</span>
            <input name="email" required type="email" />
          </label>

          <label>
            <span>Passwort</span>
            <input minLength={10} name="password" required type="password" />
          </label>

          <button type="submit">Login</button>
        </form>
      </section>
    </AdminShell>
  );
}
