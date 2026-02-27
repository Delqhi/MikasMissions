import { buildHomeHubViewModel, type ContinueWatchingItem, type HomeHubInput } from "../lib/home_hub_model";
import type { KidsMode, ProfileCard, RailItem } from "../lib/experience_types";
import { fetchChildProfiles, type ChildProfileSummary } from "../lib/fetch_child_profiles";
import { fetchKidsHome } from "../lib/fetch_kids_home";
import { fetchKidsProgress } from "../lib/fetch_kids_progress";
import { getLocaleAndMessages, supportedLocales, withLocalePath, type Locale } from "../lib/i18n";
import { demoProfiles, kidsHomeFallback, kidsProgressFallback } from "../lib/mock_payloads";
import { accessTokenFromCookie, parentUserIDFromCookie } from "../lib/server_auth";
import styles from "./page.module.css";

export const dynamic = "force-dynamic";

function modeForAgeBand(ageBand: string): KidsMode {
  if (ageBand === "3-5") {
    return "early";
  }
  if (ageBand === "12-16") {
    return "teen";
  }
  return "core";
}

function subtitleForMode(mode: KidsMode): string {
  if (mode === "early") {
    return "Audio-guided missions";
  }
  if (mode === "teen") {
    return "Explore with active safety";
  }
  return "Curated learning rails";
}

function toProfileCard(profile: ChildProfileSummary): ProfileCard {
  const mode = modeForAgeBand(profile.age_band);
  return {
    profile_id: profile.child_profile_id,
    name: profile.display_name,
    age_band: profile.age_band,
    mode,
    subtitle: subtitleForMode(mode),
    href: `/kids/${mode}`
  };
}

function fallbackInput(profile: ProfileCard): HomeHubInput {
  return {
    profile,
    mode: profile.mode,
    home: {
      ...kidsHomeFallback[profile.mode],
      child_profile_id: profile.profile_id
    },
    progress: {
      ...kidsProgressFallback[profile.mode],
      child_profile_id: profile.profile_id
    }
  };
}

async function loadProfiles(token: string, parentUserID: string): Promise<ProfileCard[]> {
  if (token === "" || parentUserID === "") {
    return demoProfiles;
  }

  try {
    const profiles = await fetchChildProfiles(parentUserID, token);
    if (profiles.length === 0) {
      return demoProfiles;
    }
    return profiles.map(toProfileCard);
  } catch {
    return demoProfiles;
  }
}

async function loadHomeHubInputs(profiles: ProfileCard[], token: string): Promise<HomeHubInput[]> {
  return Promise.all(
    profiles.map(async (profile) => {
      try {
        const [home, progress] = await Promise.all([
          fetchKidsHome(profile.mode, profile.profile_id, token || undefined),
          fetchKidsProgress(profile.mode, profile.profile_id, token || undefined)
        ]);

        return {
          profile,
          mode: profile.mode,
          home,
          progress
        } satisfies HomeHubInput;
      } catch {
        return fallbackInput(profile);
      }
    })
  );
}

function modeHref(mode: KidsMode, inputs: HomeHubInput[], locale: Locale): string {
  const profile = inputs.find((item) => item.mode === mode)?.profile ?? inputs[0]?.profile ?? demoProfiles[0];
  return withLocalePath(locale, `/kids/${mode}?child_profile_id=${encodeURIComponent(profile.profile_id)}`);
}

function profileHref(profile: ProfileCard, locale: Locale): string {
  return withLocalePath(locale, `/kids/${profile.mode}?child_profile_id=${encodeURIComponent(profile.profile_id)}`);
}

function episodeHref(episode: RailItem, inputs: HomeHubInput[], locale: Locale): string {
  const mode = episode.content_suitability;
  const profile = inputs.find((item) => item.mode === mode)?.profile ?? inputs[0]?.profile ?? demoProfiles[0];
  return withLocalePath(locale, `/kids/${mode}?child_profile_id=${encodeURIComponent(profile.profile_id)}`);
}

function formatTagLine(episode: RailItem): string {
  const tags = episode.learning_tags
    .slice(0, 2)
    .map((tag) => tag.replaceAll("_", " "))
    .join(" / ");
  return tags || episode.reason_code.replaceAll("_", " ");
}

function formatDuration(durationMS: number): string {
  return `${Math.max(1, Math.round(durationMS / 60000))} min`;
}

type PosterRowProps = {
  title: string;
  subtitle: string;
  items: RailItem[];
  locale: Locale;
  inputs: HomeHubInput[];
  id?: string;
  ranked?: boolean;
};

function PosterRow({ title, subtitle, items, locale, inputs, id, ranked = false }: PosterRowProps) {
  if (items.length === 0) {
    return null;
  }

  return (
    <section className={styles.section} id={id}>
      <header className={styles.sectionHead}>
        <h2>{title}</h2>
        <p>{subtitle}</p>
      </header>

      <div className={styles.posterRail}>
        {items.map((episode, index) => (
          <a className={styles.posterCard} href={episodeHref(episode, inputs, locale)} key={`${title}-${episode.episode_id}-${index}`}>
            <img alt={episode.title} loading="lazy" src={episode.thumbnail_url} />
            <div className={styles.posterShade} />
            {ranked ? <span className={styles.rankBadge}>{index + 1}</span> : null}
            <div className={styles.posterMeta}>
              <h3>{episode.title}</h3>
              <p>{episode.summary}</p>
              <span>
                {formatDuration(episode.duration_ms)} / {formatTagLine(episode)}
              </span>
            </div>
          </a>
        ))}
      </div>
    </section>
  );
}

type ContinueWatchingRailProps = {
  title: string;
  subtitle: string;
  items: ContinueWatchingItem[];
  locale: Locale;
};

function ContinueWatchingRail({ title, subtitle, items, locale }: ContinueWatchingRailProps) {
  return (
    <section className={styles.section} id="continue">
      <header className={styles.sectionHead}>
        <h2>{title}</h2>
        <p>{subtitle}</p>
      </header>

      {items.length === 0 ? (
        <article className={styles.emptyState}>No adventures started yet.</article>
      ) : (
        <div className={styles.continueRail}>
          {items.map((item) => (
            <a className={styles.continueCard} href={withLocalePath(locale, item.href)} key={item.childProfileID}>
              <img alt={item.episode.title} loading="lazy" src={item.episode.thumbnail_url} />
              <div className={styles.continueShade} />
              <div className={styles.continueBody}>
                <span>{item.profileName}</span>
                <h3>{item.episode.title}</h3>
                <p>{formatDuration(item.episode.duration_ms)}</p>
                <div className={styles.progressTrack}>
                  <span style={{ width: `${Math.max(8, item.completionPercent)}%` }} />
                </div>
              </div>
            </a>
          ))}
        </div>
      )}
    </section>
  );
}

export default async function HomePage() {
  const [{ locale, messages }, token, parentUserID] = await Promise.all([
    getLocaleAndMessages(),
    accessTokenFromCookie(),
    parentUserIDFromCookie()
  ]);

  const profiles = await loadProfiles(token, parentUserID);
  const hubInputs = await loadHomeHubInputs(profiles, token);
  const hub = buildHomeHubViewModel(hubInputs);

  const featured = hub.featured ?? hubInputs[0]?.home.rails[0] ?? null;
  const featuredMode = featured?.content_suitability ?? "core";
  const featuredTitle = featured?.title ?? messages.common.brand;
  const featuredSummary = featured?.summary ?? messages.home.heroFallbackSummary;
  const featuredImage = featured?.thumbnail_url ?? hubInputs[0]?.home.rails[0]?.thumbnail_url ?? "";
  const featuredHref = featured ? episodeHref(featured, hubInputs, locale) : modeHref("core", hubInputs, locale);
  const featuredKicker =
    featuredMode === "early"
      ? messages.kids.page.featuredKicker.early
      : featuredMode === "teen"
        ? messages.kids.page.featuredKicker.teen
        : messages.kids.page.featuredKicker.core;

  const fallbackEpisodes = hubInputs.flatMap((input) => input.home.rails);
  const topEpisodes = (hub.top10.length > 0 ? hub.top10.map((item) => item.episode) : fallbackEpisodes).slice(0, 12);
  const forYouEpisodes = (hub.categoryRows.forYou.length > 0 ? hub.categoryRows.forYou : topEpisodes).slice(0, 12);
  const knowledgeEpisodes = (hub.categoryRows.knowledge.length > 0 ? hub.categoryRows.knowledge : topEpisodes).slice(0, 12);
  const creativeEpisodes = (hub.categoryRows.creative.length > 0 ? hub.categoryRows.creative : topEpisodes).slice(0, 12);
  const adventureEpisodes = (hub.categoryRows.adventure.length > 0 ? hub.categoryRows.adventure : topEpisodes).slice(0, 12);

  const profileHighlights = hub.profiles.map((profile) => {
    const highlight = hubInputs.find((input) => input.profile.profile_id === profile.profile_id)?.home.rails[0] ?? featured ?? null;
    return { profile, highlight };
  });

  const sparkCards = [
    {
      icon: "GO",
      title: messages.home.sections.forYou,
      subtitle: messages.home.sections.forYouSub,
      href: "#for-you"
    },
    {
      icon: "LAB",
      title: messages.home.sections.knowledge,
      subtitle: messages.home.sections.knowledgeSub,
      href: "#top-10"
    },
    {
      icon: "FUN",
      title: messages.home.sections.adventure,
      subtitle: messages.home.sections.adventureSub,
      href: "#profiles"
    }
  ] as const;

  return (
    <main className={styles.page}>
      <a className={styles.skipLink} href="#main-content">
        Skip to main content
      </a>

      <header className={styles.topBar}>
        <a className={styles.brand} href={withLocalePath(locale, "/")}>
          {messages.common.brand}
        </a>

        <nav aria-label="Main navigation" className={styles.topNav}>
          <a href="#profiles">{messages.home.sections.chooseProfile}</a>
          <a href="#continue">{messages.home.sections.continueWatching}</a>
          <a href="#for-you">{messages.home.sections.forYou}</a>
          <a href="#top-10">{messages.home.sections.top10}</a>
          <a href="#parents">{messages.common.navParents}</a>
        </nav>

        <div className={styles.topActions}>
          <a href={withLocalePath(locale, "/parents")}>{messages.common.openParentDashboard}</a>
        </div>

        <nav aria-label="Language navigation" className={styles.localeNav}>
          {supportedLocales.map((supportedLocale) => (
            <a
              aria-current={supportedLocale === locale ? "page" : undefined}
              href={withLocalePath(supportedLocale, "/")}
              key={supportedLocale}
            >
              {supportedLocale.toUpperCase()}
            </a>
          ))}
        </nav>
      </header>

      <section aria-label="Featured mission" className={styles.hero}>
        {featuredImage ? <img alt={featuredTitle} className={styles.heroImage} src={featuredImage} /> : null}
        <div className={styles.heroShade} />
        <div className={styles.heroContent}>
          <span className={styles.heroKicker}>{featuredKicker}</span>
          <h1>{featuredTitle}</h1>
          <p>{featuredSummary}</p>
          <div className={styles.heroActions}>
            <a className={styles.playButton} href={featuredHref}>
              {messages.home.playNow}
            </a>
            <a className={styles.ghostButton} href="#profiles">
              {messages.common.openMode}
            </a>
          </div>
          <ul className={styles.heroStats}>
            <li>
              <strong>{hub.profiles.length}</strong>
              <span>{messages.home.sections.chooseProfile}</span>
            </li>
            <li>
              <strong>{hub.continueWatching.length}</strong>
              <span>{messages.home.sections.continueWatching}</span>
            </li>
            <li>
              <strong>{topEpisodes.length}</strong>
              <span>{messages.home.sections.top10}</span>
            </li>
          </ul>
        </div>
      </section>

      <section className={styles.sparkStrip}>
        {sparkCards.map((card) => (
          <a className={styles.sparkCard} href={card.href} key={card.title}>
            <span className={styles.sparkIcon}>{card.icon}</span>
            <div>
              <h3>{card.title}</h3>
              <p>{card.subtitle}</p>
            </div>
          </a>
        ))}
      </section>

      <section className={styles.section} id="profiles">
        <header className={styles.sectionHead}>
          <h2>{messages.home.sections.chooseProfile}</h2>
          <p>{messages.home.sections.continueWatchingSub}</p>
        </header>

        <div className={styles.profileRail}>
          {profileHighlights.map(({ profile, highlight }) => (
            <a className={styles.profileCard} href={profileHref(profile, locale)} key={profile.profile_id}>
              {highlight ? <img alt={highlight.title} loading="lazy" src={highlight.thumbnail_url} /> : null}
              <div className={styles.profileShade} />
              <div className={styles.profileBody}>
                <span>{messages.common.profileAgePrefix} {profile.age_band}</span>
                <h3>{profile.name}</h3>
                <p>{profile.subtitle}</p>
              </div>
            </a>
          ))}
        </div>
      </section>

      <div id="main-content">
        <ContinueWatchingRail
          items={hub.continueWatching}
          locale={locale}
          subtitle={messages.home.sections.continueWatchingSub}
          title={messages.home.sections.continueWatching}
        />

        <PosterRow
          id="for-you"
          inputs={hubInputs}
          items={forYouEpisodes}
          locale={locale}
          subtitle={messages.home.sections.forYouSub}
          title={messages.home.sections.forYou}
        />

        <PosterRow
          inputs={hubInputs}
          items={knowledgeEpisodes}
          locale={locale}
          subtitle={messages.home.sections.knowledgeSub}
          title={messages.home.sections.knowledge}
        />

        <PosterRow
          inputs={hubInputs}
          items={creativeEpisodes}
          locale={locale}
          subtitle={messages.home.sections.creativeSub}
          title={messages.home.sections.creative}
        />

        <PosterRow
          inputs={hubInputs}
          items={adventureEpisodes}
          locale={locale}
          subtitle={messages.home.sections.adventureSub}
          title={messages.home.sections.adventure}
        />

        <PosterRow
          id="top-10"
          inputs={hubInputs}
          items={topEpisodes}
          locale={locale}
          ranked
          subtitle={messages.home.sections.top10Sub}
          title={messages.home.sections.top10}
        />
      </div>

      <section className={styles.parentPanel} id="parents">
        <div>
          <h2>{messages.home.sections.trust}</h2>
          <p>{messages.home.sections.trustSub}</p>
        </div>
        <div className={styles.parentActions}>
          <a className={styles.parentPrimary} href={withLocalePath(locale, "/parents")}>
            {messages.common.openParentDashboard}
          </a>
          <a className={styles.parentSecondary} href={withLocalePath(locale, "/parents/onboarding")}>
            {messages.home.startFamilySetup}
          </a>
        </div>
      </section>

      <footer className={styles.footer}>
        <p>{messages.common.brand}</p>
        <div>
          <a href={withLocalePath(locale, "/parents")}>{messages.common.navParents}</a>
          <a href={withLocalePath(locale, "/admin/studio")}>{messages.common.navStudio}</a>
        </div>
      </footer>
    </main>
  );
}
