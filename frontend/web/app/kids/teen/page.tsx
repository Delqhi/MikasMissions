import { KidsShell } from "../../../components/layout/kids_shell";
import { SectionHeading } from "../../../components/ui/section_heading";
import { StoryCard } from "../../../components/ui/story_card";
import { EpisodeTile } from "../../../components/ui/episode_tile";
import { ParentGatePrompt } from "../../../components/ui/parent_gate_prompt";
import { fetchKidsHome } from "../../../lib/fetch_kids_home";
import { fetchKidsProgress } from "../../../lib/fetch_kids_progress";
import { accessTokenFromCookie } from "../../../lib/server_auth";
import styles from "../kids_page.module.css";

export const dynamic = "force-dynamic";

type KidsPageProps = {
  searchParams: Promise<Record<string, string | string[] | undefined>>;
};

export default async function KidsTeenPage({ searchParams }: KidsPageProps) {
  const params = await searchParams;
  const childProfileID = typeof params.child_profile_id === "string" ? params.child_profile_id : "";
  if (childProfileID === "") {
    return (
      <KidsShell
        activeNav="Explore"
        ageBand="12-16"
        mode="teen"
        navItems={[
          { label: "Explore", href: "/kids/teen#featured" },
          { label: "Watchlist", href: "/kids/teen#missions" },
          { label: "Series", href: "/kids/teen#picks" },
          { label: "Learning", href: "/kids/teen#status" },
          { label: "Reports", href: "/parents/onboarding" }
        ]}
        profileName="Setup required"
        sessionLimitMinutes={0}
        subtitle="Please create a child profile first"
        watchedMinutes={0}
      >
        <div className="stack">
          <section className={styles.focusCard}>
            <h3>Missing child profile</h3>
            <p>Open parent onboarding and create a profile to continue.</p>
            <a href="/parents/onboarding">Go to onboarding</a>
          </section>
        </div>
      </KidsShell>
    );
  }
  const token = await accessTokenFromCookie();
  const [home, progress] = await Promise.all([
    fetchKidsHome("teen", childProfileID, token),
    fetchKidsProgress("teen", childProfileID, token)
  ]);
  const featured = home.rails[0];

  return (
    <KidsShell
      activeNav="Explore"
      ageBand="12-16"
      mode="teen"
      navItems={[
        { label: "Explore", href: `/kids/teen?child_profile_id=${encodeURIComponent(childProfileID)}#featured` },
        { label: "Watchlist", href: `/kids/teen?child_profile_id=${encodeURIComponent(childProfileID)}#missions` },
        { label: "Series", href: `/kids/teen?child_profile_id=${encodeURIComponent(childProfileID)}#picks` },
        { label: "Learning", href: `/kids/teen?child_profile_id=${encodeURIComponent(childProfileID)}#status` },
        { label: "Reports", href: `/parents?child_profile_id=${encodeURIComponent(childProfileID)}` }
      ]}
      profileName="Mika Studio"
      sessionLimitMinutes={progress.session_limit_minutes}
      subtitle="Higher autonomy, safety filters still active"
      watchedMinutes={progress.watched_minutes_today}
    >
      <div className="stack">
        <nav className={styles.switchRow}>
          <a href={`/kids/early?child_profile_id=${encodeURIComponent(childProfileID)}`}>3-5</a>
          <a href={`/kids/core?child_profile_id=${encodeURIComponent(childProfileID)}`}>6-11</a>
          <a aria-current="page" href={`/kids/teen?child_profile_id=${encodeURIComponent(childProfileID)}`}>
            12-16
          </a>
          <a href={`/parents?child_profile_id=${encodeURIComponent(childProfileID)}`}>Parents</a>
        </nav>

        <section className={styles.heroPanel} id="featured">
          <img alt="" src={featured.thumbnail_url} />
          <div className={styles.heroShade} />
          <div className={styles.heroBody}>
            <span className={styles.heroKicker}>Teen featured</span>
            <h2>{featured.title}</h2>
            <p>{featured.summary}</p>
            <div className={styles.heroActions}>
              <a className={styles.heroPlay} href={`/kids/teen?child_profile_id=${encodeURIComponent(childProfileID)}`}>
                Continue
              </a>
              <a className={styles.heroInfo} href={`/parents?child_profile_id=${encodeURIComponent(childProfileID)}`}>
                Safety report
              </a>
            </div>
            <ul className={styles.heroMeta}>
              <li>{progress.mission_streak_days} active mission days</li>
              <li>{progress.watched_minutes_today} min today</li>
              <li>Age-aware recommendations</li>
            </ul>
          </div>
        </section>

        <section className="railSection" id="missions">
          <SectionHeading
            description="Dense but readable layout with shortcuts and explicit safety boundaries."
            title="Explore studio rails"
          />
          <div className={styles.storyRow}>
            {home.rails.map((item, index) => (
              <StoryCard index={index} item={item} key={item.episode_id} />
            ))}
          </div>
        </section>

        <section className="railSection" id="picks">
          <SectionHeading
            description="Transparent reason codes explain why each recommendation is shown."
            title="Quick picks"
          />
          <div className={styles.tileRow}>
            {home.rails.map((item) => (
              <EpisodeTile episode={item} key={item.episode_id} />
            ))}
          </div>
        </section>

        <section className={styles.calloutGrid} id="status">
          <article className={styles.focusCard}>
            <h3>Anti-rabbit-hole defaults</h3>
            <p>Session cap and finite exploration prevent endless algorithmic loops.</p>
          </article>
          <article className={styles.focusCard}>
            <h3>Watchlist with accountability</h3>
            <p>{progress.mission_streak_days} active mission days tracked without manipulative streak pressure.</p>
          </article>
          <ParentGatePrompt actionLabel="Open third-party discussion room" challengeType="device_confirm" />
        </section>
      </div>
    </KidsShell>
  );
}
