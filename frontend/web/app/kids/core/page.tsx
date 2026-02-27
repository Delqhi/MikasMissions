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
          { key: "library", label: messages.kids.nav.library, href: "/kids/core#recommended" }
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
            <div className={styles.focusActions}>
              <a href={withLocalePath(locale, "/parents/onboarding")}>{messages.kids.page.goToOnboarding}</a>
            </div>
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

  const rails = Array.isArray(home.rails) ? home.rails : [];
  const featured = rails[0] ?? null;
  const continueItems = rails.slice(0, 3);
  const quickItems = rails.slice(0, 4);
  const minutesLeft = Math.max(0, progress.session_limit_minutes - progress.watched_minutes_today);

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
          href: `/kids/core?child_profile_id=${encodeURIComponent(childProfileID)}#recommended`
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
          <a className={styles.parentSwitch} href={withLocalePath(locale, `/parents?child_profile_id=${encodeURIComponent(childProfileID)}`)}>
            {messages.kids.switcher.parents}
          </a>
        </nav>

        <section className={styles.heroPanel} id="featured">
          {featured ? <img alt={featured.title} src={featured.thumbnail_url} /> : null}
          <div className={styles.heroShade} />
          <div className={styles.heroBody}>
            <span className={styles.heroKicker}>{messages.kids.page.featuredKicker.core}</span>
            <h2>{featured?.title ?? messages.kids.page.exploreTitle}</h2>
            <p>{featured?.summary ?? messages.kids.page.exploreDescription}</p>
            <div className={styles.heroActions}>
              <a className={styles.heroPlay} href={withLocalePath(locale, `/kids/core?child_profile_id=${encodeURIComponent(childProfileID)}#continue`)}>
                {messages.kids.page.playFeatured}
              </a>
              <a className={styles.heroBrowse} href={withLocalePath(locale, `/kids/core?child_profile_id=${encodeURIComponent(childProfileID)}#recommended`)}>
                {messages.kids.nav.discover}
              </a>
            </div>
            <ul className={styles.heroMeta}>
              <li>
                <strong>{progress.completion_percent}%</strong>
                <span>{messages.kids.page.progressClarityTitle}</span>
              </li>
              <li>
                <strong>{progress.mission_streak_days}</strong>
                <span>{messages.kids.nav.missions}</span>
              </li>
              <li>
                <strong>{minutesLeft} min</strong>
                <span>{messages.kids.cards.sessionCap}</span>
              </li>
            </ul>
          </div>
        </section>

        <section className={styles.moodGrid}>
          <article className={styles.moodCard}>
            <span className={styles.moodPill}>{messages.kids.nav.progress}</span>
            <h3>
              {progress.completion_percent}% {messages.kids.page.progressClarityTitle}
            </h3>
            <p>{messages.kids.page.progressClarityDescription}</p>
          </article>
          <article className={styles.moodCard}>
            <span className={styles.moodPill}>{messages.kids.nav.discover}</span>
            <h3>{messages.kids.page.exploreTitle}</h3>
            <p>{messages.kids.page.exploreDescription}</p>
          </article>
          <article className={styles.moodCard}>
            <span className={styles.moodPill}>{messages.kids.cards.sessionCap}</span>
            <h3>{minutesLeft} min</h3>
            <p>{messages.kids.page.safetyDescription}</p>
          </article>
        </section>

        <section className={styles.quickStrip}>
          {quickItems.map((item) => (
            <a
              className={styles.quickCard}
              href={withLocalePath(locale, `/kids/core?child_profile_id=${encodeURIComponent(childProfileID)}#recommended`)}
              key={`${item.episode_id}-quick`}
            >
              <strong>{item.title}</strong>
              <span>{Math.round((Number(item.age_fit_score) || 0) * 100)}% {messages.kids.cards.ageFit}</span>
              <em>{messages.kids.nav.missions}</em>
            </a>
          ))}
        </section>

        <section className={styles.missionBanner}>
          <h3>{messages.kids.page.continueTitle}</h3>
          <p>{messages.kids.page.continueDescription}</p>
          <a href={withLocalePath(locale, `/kids/core?child_profile_id=${encodeURIComponent(childProfileID)}#continue`)}>
            {messages.kids.page.playFeatured}
          </a>
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
            {rails.map((item) => (
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
          <article className={styles.focusCard}>
            <h3>{messages.kids.page.parentControls}</h3>
            <p>{messages.kids.page.safetyDescription}</p>
            <div className={styles.focusActions}>
              <a href={withLocalePath(locale, `/parents?child_profile_id=${encodeURIComponent(childProfileID)}`)}>
                {messages.kids.switcher.parents}
              </a>
            </div>
          </article>
          <ParentGatePrompt actionLabel="Share to external platform" challengeType="pin" labels={messages.kids.gate} />
        </section>
      </div>
    </KidsShell>
  );
}
