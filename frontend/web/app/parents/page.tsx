import { ParentsShell } from "../../components/layout/parents_shell";
import { ParentControlPanel } from "../../components/ui/parent_control_panel";
import { SectionHeading } from "../../components/ui/section_heading";
import { fetchChildProfiles } from "../../lib/fetch_child_profiles";
import { fetchKidsProgress } from "../../lib/fetch_kids_progress";
import { fetchParentControls } from "../../lib/fetch_parent_controls";
import { accessTokenFromCookie, parentUserIDFromCookie } from "../../lib/server_auth";
import styles from "./page.module.css";

export const dynamic = "force-dynamic";

type ParentsPageProps = {
  searchParams: Promise<Record<string, string | string[] | undefined>>;
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

export default async function ParentsPage({ searchParams }: ParentsPageProps) {
  const token = await accessTokenFromCookie();
  const parentUserID = await parentUserIDFromCookie();
  if (token === "" || parentUserID === "") {
    return (
      <ParentsShell
        activeNav="Onboarding"
        description="A valid parent session is required before controls can be edited."
        heading="Parent onboarding required"
      >
        <section className={styles.notice}>
          <h3>Session missing</h3>
          <p>Start onboarding to create an authenticated parent session and unlock all controls.</p>
          <a href="/parents/onboarding">Open onboarding</a>
        </section>
      </ParentsShell>
    );
  }

  const profiles = await fetchChildProfiles(parentUserID, token);
  if (profiles.length === 0) {
    return (
      <ParentsShell
        activeNav="Onboarding"
        description="Create at least one child profile to activate kids mode and policy controls."
        heading="No child profiles found"
      >
        <section className={styles.notice}>
          <h3>Create first child profile</h3>
          <p>Run onboarding once to create the first profile and initialize guardrails.</p>
          <a href="/parents/onboarding">Create profile</a>
        </section>
      </ParentsShell>
    );
  }

  const params = await searchParams;
  const requestedProfileID = typeof params.child_profile_id === "string" ? params.child_profile_id : "";
  const selectedProfile =
    profiles.find((profile) => profile.child_profile_id === requestedProfileID) ?? profiles[0];

  const [controls, selectedProgress, progressByProfile] = await Promise.all([
    fetchParentControls(selectedProfile.child_profile_id, parentUserID, token),
    fetchKidsProgress(modeForAgeBand(selectedProfile.age_band), selectedProfile.child_profile_id, token),
    Promise.all(
      profiles.map(async (profile) => {
        const progress = await fetchKidsProgress(modeForAgeBand(profile.age_band), profile.child_profile_id, token);
        return { profile, progress };
      })
    )
  ]);

  const watchedTodayTotal = progressByProfile.reduce((sum, item) => sum + item.progress.watched_minutes_today, 0);
  const selectedKidsModeHref =
    `/kids/${modeForAgeBand(selectedProfile.age_band)}` +
    `?child_profile_id=${encodeURIComponent(selectedProfile.child_profile_id)}`;

  return (
    <ParentsShell
      description="Set strict defaults once, review daily behavior quickly, and adjust per child in seconds."
      heading="Family safety and learning operations"
    >
      <div className={styles.layout}>
        <section className={styles.hero} id="dashboard">
          <div className={styles.heroBody}>
            <span className={styles.kicker}>Live profile</span>
            <h2>
              {selectedProfile.display_name} · {selectedProfile.age_band}
            </h2>
            <p>
              Active safety mode is <strong>{controls.controls.safety_mode}</strong>. Session usage is continuously
              tracked and policy changes are audit logged.
            </p>
            <div className={styles.heroActions}>
              <a className={styles.primaryAction} href={selectedKidsModeHref}>
                Open kids mode
              </a>
              <a className={styles.secondaryAction} href="/parents/onboarding">
                Add another child
              </a>
            </div>
          </div>

          <ul className={styles.heroStats}>
            <li>
              <span>Total watched today</span>
              <strong>{watchedTodayTotal} min</strong>
            </li>
            <li>
              <span>Current session usage</span>
              <strong>
                {selectedProgress.session_minutes_used}/{selectedProgress.session_limit_minutes} min
              </strong>
            </li>
            <li>
              <span>Bedtime window</span>
              <strong>{controls.controls.bedtime_window}</strong>
            </li>
          </ul>
        </section>

        <section className={styles.stack}>
          <SectionHeading
            description="Switch profile context, jump into kids mode, and manage controls without losing orientation."
            title="Profiles and quick actions"
          />

          <section className={styles.reportGrid}>
            {profiles.map((profile) => {
              const isSelected = profile.child_profile_id === selectedProfile.child_profile_id;
              return (
                <article className={styles.reportCard} key={profile.child_profile_id}>
                  <span className={styles.cardTag}>{isSelected ? "Selected" : profile.age_band}</span>
                  <h3>{profile.display_name}</h3>
                  <p>
                    Age band {profile.age_band}. Policies can be tuned independently while preserving global safety
                    defaults.
                  </p>
                  <div className={styles.cardActions}>
                    <a
                      className={styles.cardPrimary}
                      href={`/parents?child_profile_id=${encodeURIComponent(profile.child_profile_id)}`}
                    >
                      Manage profile
                    </a>
                    <a
                      className={styles.cardSecondary}
                      href={`/kids/${modeForAgeBand(profile.age_band)}?child_profile_id=${encodeURIComponent(profile.child_profile_id)}`}
                    >
                      Open kids mode
                    </a>
                  </div>
                </article>
              );
            })}
          </section>
        </section>

        <ParentControlPanel childProfileID={controls.child_profile_id} initialControls={controls.controls} />

        <section className={styles.stack} id="compliance">
          <SectionHeading
            description="A compact compliance view for daily checks and policy confidence."
            title="Compliance snapshot"
          />
          <section className={styles.metricsGrid}>
            <article className={styles.metricCard}>
              <h3>Session cap compliance</h3>
              <p>
                {selectedProfile.display_name}: {selectedProgress.session_minutes_used}/
                {selectedProgress.session_limit_minutes} min
              </p>
              <span>Policy retained</span>
            </article>

            <article className={styles.metricCard}>
              <h3>Safety filter coverage</h3>
              <p>100% of recommendation rails include reason codes and safety metadata.</p>
              <span>Trust signal</span>
            </article>

            <article className={styles.metricCard}>
              <h3>Audit readiness</h3>
              <p>Every control mutation writes to the parent audit trail and can be reviewed by run.</p>
              <span>Operational log</span>
            </article>
          </section>
        </section>

        <section className={styles.gate}>
          <h3>Adaptive parent gate</h3>
          <p>
            External links, purchases, and account changes require verification. Current policy: <strong>device
            confirm + PIN fallback</strong>.
          </p>
        </section>
      </div>
    </ParentsShell>
  );
}
