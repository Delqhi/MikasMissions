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

export default async function KidsCorePage({ searchParams }: KidsPageProps) {
  const params = await searchParams;
  const childProfileID = typeof params.child_profile_id === "string" ? params.child_profile_id : "";
  if (childProfileID === "") {
    return (
      <KidsShell
        activeNav="Home"
        ageBand="6-11"
        mode="core"
        navItems={[
          { label: "Home", href: "/kids/core#featured" },
          { label: "Missions", href: "/kids/core#missions" },
          { label: "Discover", href: "/kids/core#picks" },
          { label: "Progress", href: "/kids/core#status" },
          { label: "Library", href: "/parents/onboarding" }
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
    fetchKidsHome("core", childProfileID, token),
    fetchKidsProgress("core", childProfileID, token)
  ]);
  const featured = home.rails[0];
  const featuredItems = home.rails.slice(0, 2);

  return (
    <KidsShell
      activeNav="Home"
      ageBand="6-11"
      mode="core"
      navItems={[
        { label: "Home", href: `/kids/core?child_profile_id=${encodeURIComponent(childProfileID)}#featured` },
        { label: "Missions", href: `/kids/core?child_profile_id=${encodeURIComponent(childProfileID)}#missions` },
        { label: "Discover", href: `/kids/core?child_profile_id=${encodeURIComponent(childProfileID)}#picks` },
        { label: "Progress", href: `/kids/core?child_profile_id=${encodeURIComponent(childProfileID)}#status` },
        { label: "Library", href: `/parents?child_profile_id=${encodeURIComponent(childProfileID)}` }
      ]}
      profileName="Mika Explorer"
      sessionLimitMinutes={progress.session_limit_minutes}
      subtitle="Curated search and learning rails"
      watchedMinutes={progress.watched_minutes_today}
    >
      <div className="stack">
        <nav className={styles.switchRow}>
          <a href={`/kids/early?child_profile_id=${encodeURIComponent(childProfileID)}`}>3-5</a>
          <a aria-current="page" href={`/kids/core?child_profile_id=${encodeURIComponent(childProfileID)}`}>
            6-11
          </a>
          <a href={`/kids/teen?child_profile_id=${encodeURIComponent(childProfileID)}`}>12-16</a>
          <a href={`/parents?child_profile_id=${encodeURIComponent(childProfileID)}`}>Parents</a>
        </nav>

        <section className={styles.heroPanel} id="featured">
          <img alt="" src={featured.thumbnail_url} />
          <div className={styles.heroShade} />
          <div className={styles.heroBody}>
            <span className={styles.heroKicker}>Core featured</span>
            <h2>{featured.title}</h2>
            <p>{featured.summary}</p>
            <div className={styles.heroActions}>
              <a className={styles.heroPlay} href={`/kids/core?child_profile_id=${encodeURIComponent(childProfileID)}`}>
                Play episode
              </a>
              <a className={styles.heroInfo} href={`/parents?child_profile_id=${encodeURIComponent(childProfileID)}`}>
                Parent controls
              </a>
            </div>
            <ul className={styles.heroMeta}>
              <li>{progress.completion_percent}% completed</li>
              <li>{progress.watched_minutes_today} min watched today</li>
              <li>Safety filters active</li>
            </ul>
          </div>
        </section>

        <section className="railSection" id="missions">
          <SectionHeading
            description="Mission progression is visible, motivating, and free from streak pressure mechanics."
            title="Continue your learning path"
          />
          <div className={styles.storyRow}>
            {featuredItems.map((item, index) => (
              <StoryCard index={index} item={item} key={item.episode_id} />
            ))}
          </div>
        </section>

        <section className="railSection" id="picks">
          <SectionHeading
            description="Exploration remains bounded through safe rails, quality labels, and age-fit scores."
            title="Explore safely"
          />
          <div className={styles.tileRow}>
            {home.rails.map((item) => (
              <EpisodeTile episode={item} key={item.episode_id} />
            ))}
          </div>
        </section>

        <section className={styles.calloutGrid} id="status">
          <article className={styles.focusCard}>
            <h3>Curated search only</h3>
            <p>Children navigate by pre-defined themes instead of unrestricted text search.</p>
          </article>
          <article className={styles.focusCard}>
            <h3>Progress clarity</h3>
            <p>{progress.completion_percent}% mission completion today with transparent milestones.</p>
          </article>
          <ParentGatePrompt actionLabel="Share to external platform" challengeType="pin" />
        </section>
      </div>
    </KidsShell>
  );
}
