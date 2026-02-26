import { buildHomeHubViewModel, type HomeHubInput, type TopRankedItem } from "../lib/home_hub_model";
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

function modeLabel(mode: KidsMode): string {
  if (mode === "early") {
    return "Early";
  }
  if (mode === "teen") {
    return "Teen";
  }
  return "Core";
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

function episodeHref(episode: RailItem, inputs: HomeHubInput[], locale: Locale): string {
  const mode = episode.content_suitability;
  const profile = inputs.find((item) => item.mode === mode)?.profile ?? inputs[0]?.profile ?? demoProfiles[0];
  return withLocalePath(locale, `/kids/${mode}?child_profile_id=${encodeURIComponent(profile.profile_id)}`);
}

function rankedFromEpisode(episode: RailItem): TopRankedItem {
  const ageFit = Math.round(episode.age_fit_score * 100);
  const safetyBonus = episode.safety_applied ? 10 : 0;

  return {
    episode,
    score: ageFit + safetyBonus,
    scoreBreakdown: {
      ageFit,
      safetyBonus,
      reasonBonus: 0
    },
    sourceModes: [episode.content_suitability]
  };
}

function toRankedRows(episodes: RailItem[], fallback: TopRankedItem[], limit = 8): TopRankedItem[] {
  const derived = episodes.map(rankedFromEpisode);
  const base = derived.length > 0 ? derived : fallback;
  return base.slice(0, limit);
}

type RankedCardProps = {
  item: TopRankedItem;
  rank?: number;
  locale: Locale;
  inputs: HomeHubInput[];
  scoreLabel: string;
};

function RankedCard({ item, rank, locale, inputs, scoreLabel }: RankedCardProps) {
  const href = episodeHref(item.episode, inputs, locale);
  const tagLine = item.episode.learning_tags
    .slice(0, 3)
    .map((tag) => tag.replaceAll("_", " "))
    .join(" · ");
  const source = Array.from(new Set(item.sourceModes.map(modeLabel))).join(" / ");

  return (
    <a className={styles.rankedCard} href={href}>
      <img alt={item.episode.title} loading="lazy" src={item.episode.thumbnail_url} />
      <div className={styles.rankedBody}>
        <div className={styles.rankedTopline}>
          {typeof rank === "number" ? <span className={styles.rankPill}>#{rank + 1}</span> : null}
          <span className={styles.scorePill}>
            {scoreLabel}: {item.score}
          </span>
          <span className={styles.modePill}>{source}</span>
        </div>
        <h3>{item.episode.title}</h3>
        <p>{item.episode.summary}</p>
        <span className={styles.tagLine}>{tagLine}</span>
      </div>
    </a>
  );
}

type ChannelCardProps = {
  title: string;
  subtitle: string;
  items: TopRankedItem[];
  locale: Locale;
  inputs: HomeHubInput[];
  ageFitLabel: string;
};

function ChannelCard({ title, subtitle, items, locale, inputs, ageFitLabel }: ChannelCardProps) {
  return (
    <article className={styles.channelCard}>
      <header className={styles.channelHead}>
        <h3>{title}</h3>
        <p>{subtitle}</p>
      </header>
      <ul className={styles.channelList}>
        {items.map((item) => (
          <li key={`${title}-${item.episode.episode_id}`}>
            <a className={styles.channelItem} href={episodeHref(item.episode, inputs, locale)}>
              <strong>{item.episode.title}</strong>
              <span>
                {Math.round(item.episode.age_fit_score * 100)}% {ageFitLabel}
              </span>
            </a>
          </li>
        ))}
      </ul>
    </article>
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
  const fallbackSummary = messages.home.heroFallbackSummary;
  const topRows = hub.top10;
  const forYouRows =
    hub.categoryRows.forYou.length > 0
      ? hub.categoryRows.forYou.slice(0, 8).map(rankedFromEpisode)
      : hub.top10.slice(0, 8);
  const knowledgeRows = toRankedRows(hub.categoryRows.knowledge, topRows, 5);
  const creativeRows = toRankedRows(hub.categoryRows.creative, topRows, 5);
  const adventureRows = toRankedRows(hub.categoryRows.adventure, topRows, 5);

  const readinessIndex = Math.round(
    topRows.reduce((sum, item) => sum + item.score, 0) / Math.max(topRows.length, 1)
  );
  const safeCoverage = topRows.filter((item) => item.episode.safety_applied).length;
  const knowledgeCoverage = new Set(hub.top10.flatMap((item) => item.episode.learning_tags)).size;
  const activeSessions = hub.continueWatching.length;

  return (
    <main className={styles.page}>
      <a className={styles.skipLink} href="#main-content">
        Skip to main content
      </a>
      <div aria-hidden="true" className={styles.gridGlow} />

      <header className={styles.topBar}>
        <a className={styles.brand} href={withLocalePath(locale, "/")}>
          {messages.common.brand}
        </a>

        <nav aria-label="Main navigation" className={styles.topNav}>
          <a href="#worlds">{messages.home.sections.chooseProfile}</a>
          <a href="#missions">{messages.home.sections.continueWatching}</a>
          <a href="#leaderboard">{messages.home.sections.top10}</a>
          <a href="#trust">{messages.home.sections.trust}</a>
        </nav>

        <div className={styles.topActions}>
          <a href={withLocalePath(locale, "/parents")}>{messages.common.navParents}</a>
          <a href={withLocalePath(locale, "/parents/onboarding")}>{messages.common.navFamilySetup}</a>
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

      <section aria-label="Featured mission" className={styles.heroLayout}>
        <div className={styles.heroStage}>
          {featured ? <img alt={featured.title} className={styles.heroImage} src={featured.thumbnail_url} /> : null}
          <div className={styles.heroShade} />

          <div className={styles.heroContent}>
            <span className={styles.heroTag}>{messages.home.heroTag}</span>
            <h1>{featured?.title ?? messages.common.brand}</h1>
            <p>{featured?.summary ?? fallbackSummary}</p>

            <div className={styles.heroActions}>
              <a className={styles.launchButton} href={modeHref("core", hubInputs, locale)}>
                {messages.home.playNow}
              </a>
              <a className={styles.ghostButton} href="#worlds">
                {messages.common.openMode}
              </a>
            </div>

            <ul className={styles.heroStats}>
              <li>
                <strong>{hub.profiles.length}</strong>
                <span>{messages.home.sections.chooseProfile}</span>
              </li>
              <li>
                <strong>{activeSessions}</strong>
                <span>{messages.home.sections.continueWatching}</span>
              </li>
              <li>
                <strong>{readinessIndex}</strong>
                <span>{messages.home.labels.scoreLabel}</span>
              </li>
            </ul>
          </div>
        </div>

        <aside className={styles.parentPanel}>
          <h2>{messages.home.sections.trust}</h2>
          <p>{messages.home.sections.trustSub}</p>
          <ul className={styles.parentStats}>
            <li>
              <span>{messages.home.safetyHighlights.ageModes}</span>
              <strong>{messages.home.safetyHighlights.ageModesValue}</strong>
            </li>
            <li>
              <span>{messages.home.safetyHighlights.safetyFilter}</span>
              <strong>{messages.home.safetyHighlights.safetyFilterValue}</strong>
            </li>
            <li>
              <span>{messages.home.safetyHighlights.parentControl}</span>
              <strong>{messages.home.safetyHighlights.parentControlValue}</strong>
            </li>
            <li>
              <span>{messages.home.safetyHighlights.sessionCaps}</span>
              <strong>{messages.home.safetyHighlights.sessionCapsValue}</strong>
            </li>
          </ul>
          <a className={styles.parentPrimary} href={withLocalePath(locale, "/parents")}>
            {messages.common.openParentDashboard}
          </a>
          <a className={styles.parentSecondary} href={withLocalePath(locale, "/parents/onboarding")}>
            {messages.home.startFamilySetup}
          </a>
        </aside>
      </section>

      <section aria-label="Mission operations" className={styles.controlStrip}>
        <article className={styles.controlCard}>
          <span className={styles.controlLabel}>{messages.home.labels.scoreLabel}</span>
          <strong className={styles.controlValue}>{readinessIndex}</strong>
          <p>{messages.home.sections.top10}</p>
        </article>
        <article className={styles.controlCard}>
          <span className={styles.controlLabel}>{messages.home.safetyHighlights.safetyFilter}</span>
          <strong className={styles.controlValue}>
            {safeCoverage}/{Math.max(topRows.length, 1)}
          </strong>
          <p>{messages.home.trustCards.filterCoverageTitle}</p>
        </article>
        <article className={styles.controlCard}>
          <span className={styles.controlLabel}>{messages.home.sections.knowledge}</span>
          <strong className={styles.controlValue}>{knowledgeCoverage}</strong>
          <p>{messages.home.sections.knowledgeSub}</p>
        </article>
        <article className={styles.controlCard}>
          <span className={styles.controlLabel}>{messages.home.sections.continueWatching}</span>
          <strong className={styles.controlValue}>{activeSessions}</strong>
          <p>{messages.home.sections.continueWatchingSub}</p>
        </article>
      </section>

      <section className={styles.worlds} id="worlds">
        <div className={styles.sectionHead}>
          <h2>{messages.home.sections.chooseProfile}</h2>
          <span>{messages.home.sections.continueWatchingSub}</span>
        </div>

        <div className={styles.worldGrid}>
          <article className={`${styles.worldCard} ${styles.world_early}`}>
            <span className={styles.worldAge}>
              {messages.common.profileAgePrefix} 3-5
            </span>
            <h3>{messages.home.modeDeck.earlyTitle}</h3>
            <p>{messages.home.modeDeck.earlyDetail}</p>
            <a href={modeHref("early", hubInputs, locale)}>{messages.common.openMode}</a>
          </article>

          <article className={`${styles.worldCard} ${styles.world_core}`}>
            <span className={styles.worldAge}>
              {messages.common.profileAgePrefix} 6-11
            </span>
            <h3>{messages.home.modeDeck.coreTitle}</h3>
            <p>{messages.home.modeDeck.coreDetail}</p>
            <a href={modeHref("core", hubInputs, locale)}>{messages.common.openMode}</a>
          </article>

          <article className={`${styles.worldCard} ${styles.world_teen}`}>
            <span className={styles.worldAge}>
              {messages.common.profileAgePrefix} 12-16
            </span>
            <h3>{messages.home.modeDeck.teenTitle}</h3>
            <p>{messages.home.modeDeck.teenDetail}</p>
            <a href={modeHref("teen", hubInputs, locale)}>{messages.common.openMode}</a>
          </article>
        </div>

        <div className={styles.profileRail}>
          {hub.profiles.map((profile) => (
            <a
              className={`${styles.profileChip} ${styles[`profile_${profile.mode}`]}`}
              href={withLocalePath(locale, `/kids/${profile.mode}?child_profile_id=${encodeURIComponent(profile.profile_id)}`)}
              key={profile.profile_id}
            >
              <div className={styles.profileOrb}>{profile.name.slice(0, 1)}</div>
              <div>
                <strong>{profile.name}</strong>
                <p>{profile.subtitle}</p>
              </div>
              <span className={styles.profileMeta}>
                {messages.common.profileAgePrefix} {profile.age_band}
              </span>
            </a>
          ))}
        </div>
      </section>

      <div id="main-content">
        <section className={styles.section} id="missions">
          <div className={styles.sectionHead}>
            <h2>{messages.home.sections.continueWatching}</h2>
            <span>{messages.home.sections.continueWatchingSub}</span>
          </div>

          <div className={styles.continueRail}>
            {hub.continueWatching.length === 0 ? (
              <article className={styles.emptyCard}>{messages.home.labels.noContinueWatching}</article>
            ) : (
              hub.continueWatching.map((item) => (
                <a className={styles.continueCard} href={withLocalePath(locale, item.href)} key={item.childProfileID}>
                  <img alt={item.episode.title} loading="lazy" src={item.episode.thumbnail_url} />
                  <div className={styles.continueBody}>
                    <span>{item.profileName}</span>
                    <h3>{item.episode.title}</h3>
                    <p>{item.watchedMinutesToday} min today</p>
                    <div className={styles.progressTrack}>
                      <span style={{ width: `${Math.max(item.completionPercent, 8)}%` }} />
                    </div>
                  </div>
                </a>
              ))
            )}
          </div>
        </section>

        <section className={styles.section}>
          <div className={styles.sectionHead}>
            <h2>{messages.home.sections.forYou}</h2>
            <span>{messages.home.sections.forYouSub}</span>
          </div>
          <div className={styles.rankedRail}>
            {forYouRows.map((item) => (
              <RankedCard
                inputs={hubInputs}
                item={item}
                key={`${item.episode.episode_id}-for-you`}
                locale={locale}
                scoreLabel={messages.home.labels.scoreLabel}
              />
            ))}
          </div>
        </section>

        <section className={styles.section}>
          <div className={styles.sectionHead}>
            <h2>{messages.home.sections.forYou}</h2>
            <span>{messages.home.sections.top10Sub}</span>
          </div>
          <div className={styles.channelGrid}>
            <ChannelCard
              ageFitLabel={messages.kids.cards.ageFit}
              inputs={hubInputs}
              items={knowledgeRows}
              locale={locale}
              subtitle={messages.home.sections.knowledgeSub}
              title={messages.home.sections.knowledge}
            />
            <ChannelCard
              ageFitLabel={messages.kids.cards.ageFit}
              inputs={hubInputs}
              items={creativeRows}
              locale={locale}
              subtitle={messages.home.sections.creativeSub}
              title={messages.home.sections.creative}
            />
            <ChannelCard
              ageFitLabel={messages.kids.cards.ageFit}
              inputs={hubInputs}
              items={adventureRows}
              locale={locale}
              subtitle={messages.home.sections.adventureSub}
              title={messages.home.sections.adventure}
            />
          </div>
        </section>

        <section className={styles.section} id="leaderboard">
          <div className={styles.sectionHead}>
            <h2>{messages.home.sections.top10}</h2>
            <span>{messages.home.sections.top10Sub}</span>
          </div>
          <div className={styles.rankedRail}>
            {hub.top10.map((item, index) => (
              <RankedCard
                inputs={hubInputs}
                item={item}
                key={item.episode.episode_id}
                locale={locale}
                rank={index}
                scoreLabel={messages.home.labels.scoreLabel}
              />
            ))}
          </div>
        </section>

        <section className={styles.section} id="trust">
          <div className={styles.sectionHead}>
            <h2>{messages.home.sections.trust}</h2>
            <span>{messages.home.sections.trustSub}</span>
          </div>
          <div className={styles.trustGrid}>
            <article className={styles.trustCard}>
              <h3>{messages.home.trustCards.filterCoverageTitle}</h3>
              <p>{messages.home.trustCards.filterCoverageDetail}</p>
            </article>
            <article className={styles.trustCard}>
              <h3>{messages.home.trustCards.sessionLimitsTitle}</h3>
              <p>{messages.home.trustCards.sessionLimitsDetail}</p>
            </article>
            <article className={styles.trustCard}>
              <h3>{messages.home.trustCards.parentGateTitle}</h3>
              <p>{messages.home.trustCards.parentGateDetail}</p>
            </article>
          </div>
        </section>
      </div>

      <footer className={styles.footer}>
        <p>{messages.common.brand} · missions first, parent controls always available</p>
        <div>
          <a href={withLocalePath(locale, "/parents")}>{messages.common.navParents}</a>
          <a href={withLocalePath(locale, "/admin/studio")}>{messages.common.navStudio}</a>
        </div>
      </footer>
    </main>
  );
}
