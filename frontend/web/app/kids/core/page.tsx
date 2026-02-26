import { KidsShell } from "../../../components/layout/kids_shell";
import { SectionHeading } from "../../../components/ui/section_heading";
import { StoryCard } from "../../../components/ui/story_card";
import { EpisodeTile } from "../../../components/ui/episode_tile";
import { ParentGatePrompt } from "../../../components/ui/parent_gate_prompt";
import { fetchKidsHome } from "../../../lib/fetch_kids_home";
import { fetchKidsProgress } from "../../../lib/fetch_kids_progress";
import { getLocaleAndMessages, withLocalePath } from "../../../lib/i18n";
import { accessTokenFromCookie } from "../../../lib/server_auth";
import styles from "../kids_page.module.css";

export const dynamic = "force-dynamic";

type KidsPageProps = {
  searchParams: Promise<Record<string, string | string[] | undefined>>;
};

export default async function KidsCorePage({ searchParams }: KidsPageProps) {
  const [{ locale, messages }, params] = await Promise.all([getLocaleAndMessages(), searchParams]);

  const childProfileID = typeof params.child_profile_id === "string" ? params.child_profile_id : "";

  if (childProfileID === "") {
    return (
      <KidsShell
        activeNavKey="home"
        ageBand="6-11"
        locale={locale}
        messages={messages.kids}
        mode="core"
        navItems={[
          { key: "home", label: messages.kids.nav.home, href: "/kids/core#featured" },
          { key: "missions", label: messages.kids.nav.missions, href: "/kids/core#continue" },
          { key: "discover", label: messages.kids.nav.discover, href: "/kids/core#recommended" },
          { key: "progress", label: messages.kids.nav.progress, href: "/kids/core#safety" },
          { key: "library", label: messages.kids.nav.library, href: "/parents/onboarding" }
        ]}
        profileName={messages.kids.page.missingProfileTitle}
        sessionLimitMinutes={0}
        subtitle={messages.kids.page.missingProfileText}
        watchedMinutes={0}
      >
        <div className="stack">
          <section className={styles.focusCard}>
            <h3>{messages.kids.page.missingProfileTitle}</h3>
            <p>{messages.kids.page.missingProfileText}</p>
            <a href={withLocalePath(locale, "/parents/onboarding")}>{messages.kids.page.goToOnboarding}</a>
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
  const continueItems = home.rails.slice(0, 3);

  return (
    <KidsShell
      activeNavKey="home"
      ageBand="6-11"
      locale={locale}
      messages={messages.kids}
      mode="core"
      navItems={[
        { key: "home", label: messages.kids.nav.home, href: `/kids/core?child_profile_id=${encodeURIComponent(childProfileID)}#featured` },
        {
          key: "missions",
          label: messages.kids.nav.missions,
          href: `/kids/core?child_profile_id=${encodeURIComponent(childProfileID)}#continue`
        },
        {
          key: "discover",
          label: messages.kids.nav.discover,
          href: `/kids/core?child_profile_id=${encodeURIComponent(childProfileID)}#recommended`
        },
        {
          key: "progress",
          label: messages.kids.nav.progress,
          href: `/kids/core?child_profile_id=${encodeURIComponent(childProfileID)}#safety`
        },
        {
          key: "library",
          label: messages.kids.nav.library,
          href: `/parents?child_profile_id=${encodeURIComponent(childProfileID)}`
        }
      ]}
      profileName="Mika Explorer"
      sessionLimitMinutes={progress.session_limit_minutes}
      subtitle={messages.kids.page.curatedSearchTitle}
      watchedMinutes={progress.watched_minutes_today}
    >
      <div className="stack">
        <nav className={styles.switchRow}>
          <a href={withLocalePath(locale, `/kids/early?child_profile_id=${encodeURIComponent(childProfileID)}`)}>
            {messages.kids.switcher.early}
          </a>
          <a aria-current="page" href={withLocalePath(locale, `/kids/core?child_profile_id=${encodeURIComponent(childProfileID)}`)}>
            {messages.kids.switcher.core}
          </a>
          <a href={withLocalePath(locale, `/kids/teen?child_profile_id=${encodeURIComponent(childProfileID)}`)}>
            {messages.kids.switcher.teen}
          </a>
          <a href={withLocalePath(locale, `/parents?child_profile_id=${encodeURIComponent(childProfileID)}`)}>
            {messages.kids.switcher.parents}
          </a>
        </nav>

        <section className={styles.heroPanel} id="featured">
          <img alt="" src={featured.thumbnail_url} />
          <div className={styles.heroShade} />
          <div className={styles.heroBody}>
            <span className={styles.heroKicker}>{messages.kids.page.featuredKicker.core}</span>
            <h2>{featured.title}</h2>
            <p>{featured.summary}</p>
            <div className={styles.heroActions}>
              <a className={styles.heroPlay} href={withLocalePath(locale, `/kids/core?child_profile_id=${encodeURIComponent(childProfileID)}`)}>
                {messages.kids.page.playFeatured}
              </a>
              <a className={styles.heroInfo} href={withLocalePath(locale, `/parents?child_profile_id=${encodeURIComponent(childProfileID)}`)}>
                {messages.kids.page.parentControls}
              </a>
            </div>
            <ul className={styles.heroMeta}>
              <li>{progress.completion_percent}% completed</li>
              <li>{progress.watched_minutes_today} min today</li>
              <li>{messages.kids.page.curatedSearchTitle}</li>
            </ul>
          </div>
        </section>

        <section className="railSection" id="continue">
          <SectionHeading description={messages.kids.page.continueDescription} title={messages.kids.page.continueTitle} />
          <div className={styles.storyRow}>
            {continueItems.map((item, index) => (
              <StoryCard index={index} item={item} key={item.episode_id} labels={messages.kids.cards} />
            ))}
          </div>
        </section>

        <section className="railSection" id="recommended">
          <SectionHeading description={messages.kids.page.exploreDescription} title={messages.kids.page.exploreTitle} />
          <div className={styles.tileRow}>
            {home.rails.map((item) => (
              <EpisodeTile episode={item} key={item.episode_id} labels={messages.kids.cards} />
            ))}
          </div>
        </section>

        <section className={styles.calloutGrid} id="safety">
          <article className={styles.focusCard}>
            <h3>{messages.kids.page.curatedSearchTitle}</h3>
            <p>{messages.kids.page.curatedSearchDescription}</p>
          </article>
          <article className={styles.focusCard}>
            <h3>{messages.kids.page.progressClarityTitle}</h3>
            <p>{messages.kids.page.progressClarityDescription}</p>
          </article>
          <ParentGatePrompt actionLabel="Share to external platform" challengeType="pin" labels={messages.kids.gate} />
        </section>
      </div>
    </KidsShell>
  );
}
