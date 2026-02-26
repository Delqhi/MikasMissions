import { ParentsShell } from "../../components/layout/parents_shell";
import { ParentControlPanel } from "../../components/ui/parent_control_panel";
import { SectionHeading } from "../../components/ui/section_heading";
import { fetchChildProfiles } from "../../lib/fetch_child_profiles";
import { fetchKidsProgress } from "../../lib/fetch_kids_progress";
import { fetchParentControls } from "../../lib/fetch_parent_controls";
import { getLocaleAndMessages, withLocalePath } from "../../lib/i18n";
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
  const [{ locale, messages }, token, parentUserID] = await Promise.all([
    getLocaleAndMessages(),
    accessTokenFromCookie(),
    parentUserIDFromCookie()
  ]);

  if (token === "" || parentUserID === "") {
    return (
      <ParentsShell
        activeNav="onboarding"
        description={messages.parents.page.sessionMissingText}
        heading={messages.parents.page.onboardingRequired}
        locale={locale}
        messages={messages.parents}
      >
        <section className={styles.notice}>
          <h3>{messages.parents.page.sessionMissingTitle}</h3>
          <p>{messages.parents.page.sessionMissingText}</p>
          <a href={withLocalePath(locale, "/parents/onboarding")}>{messages.parents.page.openOnboarding}</a>
        </section>
      </ParentsShell>
    );
  }

  const profiles = await fetchChildProfiles(parentUserID, token);
  if (profiles.length === 0) {
    return (
      <ParentsShell
        activeNav="onboarding"
        description={messages.parents.page.createProfileText}
        heading={messages.parents.page.noProfilesFound}
        locale={locale}
        messages={messages.parents}
      >
        <section className={styles.notice}>
          <h3>{messages.parents.page.createProfileTitle}</h3>
          <p>{messages.parents.page.createProfileText}</p>
          <a href={withLocalePath(locale, "/parents/onboarding")}>{messages.parents.page.openOnboarding}</a>
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
      description={messages.parents.page.familySafetyDescription}
      heading={messages.parents.page.familySafetyHeading}
      locale={locale}
      messages={messages.parents}
    >
      <div className={styles.layout}>
        <section className={styles.hero} id="dashboard">
          <div className={styles.heroBody}>
            <span className={styles.kicker}>{messages.parents.page.liveProfile}</span>
            <h2>
              {selectedProfile.display_name} · {selectedProfile.age_band}
            </h2>
            <p>
              Active safety mode is <strong>{controls.controls.safety_mode}</strong>. Session usage is continuously
              tracked and policy changes are audit logged.
            </p>
            <div className={styles.heroActions}>
              <a className={styles.primaryAction} href={withLocalePath(locale, selectedKidsModeHref)}>
                {messages.parents.page.openKidsMode}
              </a>
              <a className={styles.secondaryAction} href={withLocalePath(locale, "/parents/onboarding")}>
                {messages.parents.page.addAnotherChild}
              </a>
            </div>
          </div>

          <ul className={styles.heroStats}>
            <li>
              <span>{messages.parents.page.totalWatchedToday}</span>
              <strong>{watchedTodayTotal} min</strong>
            </li>
            <li>
              <span>{messages.parents.page.currentSessionUsage}</span>
              <strong>
                {selectedProgress.session_minutes_used}/{selectedProgress.session_limit_minutes} min
              </strong>
            </li>
            <li>
              <span>{messages.parents.page.bedtimeWindow}</span>
              <strong>{controls.controls.bedtime_window}</strong>
            </li>
          </ul>
        </section>

        <section className={styles.stack}>
          <SectionHeading
            description={messages.parents.page.profilesQuickActionsDescription}
            title={messages.parents.page.profilesQuickActionsTitle}
          />

          <section className={styles.reportGrid}>
            {profiles.map((profile) => {
              const isSelected = profile.child_profile_id === selectedProfile.child_profile_id;
              return (
                <article className={styles.reportCard} key={profile.child_profile_id}>
                  <span className={styles.cardTag}>{isSelected ? messages.parents.page.selected : profile.age_band}</span>
                  <h3>{profile.display_name}</h3>
                  <p>
                    Age band {profile.age_band}. Policies can be tuned independently while preserving global safety
                    defaults.
                  </p>
                  <div className={styles.cardActions}>
                    <a
                      className={styles.cardPrimary}
                      href={withLocalePath(
                        locale,
                        `/parents?child_profile_id=${encodeURIComponent(profile.child_profile_id)}`
                      )}
                    >
                      {messages.parents.page.manageProfile}
                    </a>
                    <a
                      className={styles.cardSecondary}
                      href={withLocalePath(
                        locale,
                        `/kids/${modeForAgeBand(profile.age_band)}?child_profile_id=${encodeURIComponent(profile.child_profile_id)}`
                      )}
                    >
                      {messages.parents.page.openKidsMode}
                    </a>
                  </div>
                </article>
              );
            })}
          </section>
        </section>

        <ParentControlPanel
          childProfileID={controls.child_profile_id}
          initialControls={controls.controls}
          messages={messages.parents.controlPanel}
        />

        <section className={styles.stack} id="compliance">
          <SectionHeading
            description={messages.parents.page.complianceSnapshotDescription}
            title={messages.parents.page.complianceSnapshotTitle}
          />
          <section className={styles.metricsGrid}>
            <article className={styles.metricCard}>
              <h3>{messages.parents.page.sessionCapCompliance}</h3>
              <p>
                {selectedProfile.display_name}: {selectedProgress.session_minutes_used}/
                {selectedProgress.session_limit_minutes} min
              </p>
              <span>{messages.parents.page.policyRetained}</span>
            </article>

            <article className={styles.metricCard}>
              <h3>{messages.parents.page.safetyFilterCoverage}</h3>
              <p>100% of recommendation rails include reason codes and safety metadata.</p>
              <span>{messages.parents.page.trustSignal}</span>
            </article>

            <article className={styles.metricCard}>
              <h3>{messages.parents.page.auditReadiness}</h3>
              <p>Every control mutation writes to the parent audit trail and can be reviewed by run.</p>
              <span>{messages.parents.page.operationalLog}</span>
            </article>
          </section>
        </section>

        <section className={styles.gate}>
          <h3>{messages.parents.page.adaptiveGateTitle}</h3>
          <p>
            {messages.parents.page.adaptiveGateText} <strong>device confirm + PIN fallback</strong>.
          </p>
        </section>
      </div>
    </ParentsShell>
  );
}
