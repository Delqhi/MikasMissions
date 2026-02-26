import { cookies } from "next/headers";
import { redirect } from "next/navigation";

import { apiBaseURL, safeJSONFetch } from "../../../lib/fetch_helpers";
import { getLocaleAndMessages, getLocaleFromRequest, withLocalePath } from "../../../lib/i18n";
import styles from "./page.module.css";

type SignupResponse = {
  parent_user_id: string;
};

type LoginResponse = {
  access_token: string;
  expires_in: number;
  parent_user_id: string;
};

type ChildProfileResponse = {
  child_profile_id: string;
};

function modeForAgeBand(ageBand: string): "early" | "core" | "teen" {
  switch (ageBand) {
    case "3-5":
      return "early";
    case "12-16":
      return "teen";
    default:
      return "core";
  }
}

async function completeOnboarding(formData: FormData) {
  "use server";

  const locale = await getLocaleFromRequest();
  const email = String(formData.get("email") ?? "").trim().toLowerCase();
  const password = String(formData.get("password") ?? "");
  const displayName = String(formData.get("display_name") ?? "").trim();
  const ageBand = String(formData.get("age_band") ?? "6-11");
  const avatar = String(formData.get("avatar") ?? "robot");

  const baseURL = apiBaseURL();

  try {
    await safeJSONFetch<SignupResponse>(new URL("/v1/parents/signup", baseURL).toString(), {
      method: "POST",
      body: {
        accepted_terms: true,
        country: "DE",
        email,
        language: locale,
        marketing: false,
        password
      }
    });
  } catch {
    // Existing account is valid for onboarding: continue with login.
  }

  const login = await safeJSONFetch<LoginResponse>(new URL("/v1/parents/login", baseURL).toString(), {
    method: "POST",
    body: { email, password }
  });

  await safeJSONFetch(new URL("/v1/parents/consent/verify", baseURL).toString(), {
    method: "POST",
    body: {
      challenge: "ok",
      method: "card",
      parent_user_id: login.parent_user_id
    },
    token: login.access_token
  });

  const profile = await safeJSONFetch<ChildProfileResponse>(new URL("/v1/children/profiles", baseURL).toString(), {
    method: "POST",
    body: {
      age_band: ageBand,
      avatar,
      display_name: displayName,
      parent_user_id: login.parent_user_id
    },
    token: login.access_token
  });

  const cookieStore = await cookies();
  cookieStore.set("mm_access_token", login.access_token, {
    httpOnly: true,
    maxAge: login.expires_in,
    path: "/",
    sameSite: "lax",
    secure: process.env.NODE_ENV === "production"
  });
  cookieStore.set("mm_parent_user_id", login.parent_user_id, {
    httpOnly: true,
    maxAge: login.expires_in,
    path: "/",
    sameSite: "lax",
    secure: process.env.NODE_ENV === "production"
  });

  const mode = modeForAgeBand(ageBand);
  redirect(withLocalePath(locale, `/kids/${mode}?child_profile_id=${encodeURIComponent(profile.child_profile_id)}`));
}

export default async function ParentOnboardingPage() {
  const { locale, messages } = await getLocaleAndMessages();

  return (
    <main className={styles.wrapper}>
      <section className={styles.hero}>
        <p className={styles.kicker}>{messages.parents.onboarding.heroKicker}</p>
        <h1>{messages.parents.onboarding.heroTitle}</h1>
        <p>{messages.parents.onboarding.heroText}</p>
        <ul>
          <li>{messages.parents.onboarding.bulletOne}</li>
          <li>{messages.parents.onboarding.bulletTwo}</li>
          <li>{messages.parents.onboarding.bulletThree}</li>
        </ul>
      </section>

      <section className={styles.card}>
        <h2>{messages.parents.onboarding.launchTitle}</h2>
        <p>{messages.parents.onboarding.launchText}</p>

        <form action={completeOnboarding} className={styles.form}>
          <label>
            <span>{messages.parents.onboarding.emailLabel}</span>
            <input name="email" placeholder="parent@example.com" required type="email" />
          </label>

          <label>
            <span>{messages.parents.onboarding.passwordLabel}</span>
            <input minLength={10} name="password" required type="password" />
          </label>

          <label>
            <span>{messages.parents.onboarding.childNameLabel}</span>
            <input name="display_name" placeholder="Mika" required type="text" />
          </label>

          <div className={styles.row}>
            <label>
              <span>{messages.parents.onboarding.ageBandLabel}</span>
              <select defaultValue="6-11" name="age_band">
                <option value="3-5">3-5</option>
                <option value="6-11">6-11</option>
                <option value="12-16">12-16</option>
              </select>
            </label>

            <label>
              <span>{messages.parents.onboarding.avatarLabel}</span>
              <input defaultValue="robot" name="avatar" type="text" />
            </label>
          </div>

          <button type="submit">{messages.parents.onboarding.submit}</button>
        </form>

        <a className={styles.backLink} href={withLocalePath(locale, "/")}>
          {messages.parents.onboarding.backToHome}
        </a>
      </section>
    </main>
  );
}
