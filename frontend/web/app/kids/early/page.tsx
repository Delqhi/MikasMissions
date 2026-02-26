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

export default async function KidsEarlyPage({ searchParams }: KidsPageProps) {
  const params = await searchParams;
  const childProfileID = typeof params.child_profile_id === "string" ? params.child_profile_id : "";
  if (childProfileID === "") {
    return (
      <KidsShell
        activeNav="Stories"
        ageBand="3-5"
        mode="early"
        navItems={[
          { label: "Stories", href: "/kids/early#featured" },
          { label: "Missions", href: "/kids/early#missions" },
          { label: "Favorites", href: "/kids/early#picks" },
          { label: "Bedtime", href: "/parents/onboarding" }
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
    fetchKidsHome("early", childProfileID, token),
    fetchKidsProgress("early", childProfileID, token)
  ]);
  const featured = home.rails[0];
  const featuredItems = home.rails.slice(0, 2);

  return (
    <KidsShell
      activeNav="Stories"
      ageBand="3-5"
      mode="early"
      navItems={[
        { label: "Stories", href: `/kids/early?child_profile_id=${encodeURIComponent(childProfileID)}#featured` },
        { label: "Missions", href: `/kids/early?child_profile_id=${encodeURIComponent(childProfileID)}#missions` },
        { label: "Favorites", href: `/kids/early?child_profile_id=${encodeURIComponent(childProfileID)}#picks` },
        { label: "Bedtime", href: `/parents?child_profile_id=${encodeURIComponent(childProfileID)}` }
      ]}
      profileName="Mika Mini"
      sessionLimitMinutes={progress.session_limit_minutes}
      subtitle="Audio-guided and image-first"
      watchedMinutes={progress.watched_minutes_today}
    >
      <div className="stack">
        <nav className={styles.switchRow}>
          <a aria-current="page" href={`/kids/early?child_profile_id=${encodeURIComponent(childProfileID)}`}>
            3-5
          </a>
          <a href={`/kids/core?child_profile_id=${encodeURIComponent(childProfileID)}`}>6-11</a>
          <a href={`/kids/teen?child_profile_id=${encodeURIComponent(childProfileID)}`}>12-16</a>
          <a href={`/parents?child_profile_id=${encodeURIComponent(childProfileID)}`}>Parents</a>
        </nav>

        <section className={styles.heroPanel} id="featured">
          <img alt="" src={featured.thumbnail_url} />
          <div className={styles.heroShade} />
          <div className={styles.heroBody}>
            <span className={styles.heroKicker}>Early featured</span>
            <h2>{featured.title}</h2>
            <p>{featured.summary}</p>
            <div className={styles.heroActions}>
              <a className={styles.heroPlay} href={`/kids/early?child_profile_id=${encodeURIComponent(childProfileID)}`}>
                Start story
              </a>
              <a className={styles.heroInfo} href={`/parents?child_profile_id=${encodeURIComponent(childProfileID)}`}>
                Guardian view
              </a>
            </div>
            <ul className={styles.heroMeta}>
              <li>{progress.session_limit_minutes} min limit</li>
              <li>{progress.watched_minutes_today} min today</li>
              <li>Guided-only navigation</li>
            </ul>
          </div>
        </section>

        <section className="railSection" id="missions">
          <SectionHeading
            description="Large visuals, guided audio prompts, and no free-search surface for early learners."
            title="Today\'s guided adventures"
          />
          <div className={styles.storyRow}>
            {featuredItems.map((item, index) => (
              <StoryCard index={index} item={item} key={item.episode_id} />
            ))}
          </div>
        </section>

        <section className="railSection" id="picks">
          <SectionHeading
            description="Every route has short click paths and clear return points to avoid confusion."
            title="Safe picks"
          />
          <div className={styles.tileRow}>
            {home.rails.map((item) => (
              <EpisodeTile episode={item} key={item.episode_id} />
            ))}
          </div>
        </section>

        <section className={styles.calloutGrid}>
          <article className={styles.focusCard}>
            <h3>Audio-first guidance</h3>
            <p>Focused prompts help non-readers complete actions without trial-and-error taps.</p>
          </article>
          <article className={styles.focusCard}>
            <h3>No open exploration feed</h3>
            <p>Finite curated rails prevent infinite browsing loops for ages 3-5.</p>
          </article>
          <ParentGatePrompt actionLabel="Open external activity link" challengeType="device_confirm" />
        </section>
      </div>
    </KidsShell>
  );
}
