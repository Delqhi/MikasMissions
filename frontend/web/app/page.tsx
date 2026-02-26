import {
  buildHomeHubViewModel,
  type HomeHubInput,
  type TopRankedItem
} from "../lib/home_hub_model";
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

function episodeHref(episode: RailItem, inputs: HomeHubInput[], locale: Locale): string {
  const mode = episode.content_suitability;
  const profile = inputs.find((item) => item.mode === mode)?.profile ?? inputs[0]?.profile ?? demoProfiles[0];
  return withLocalePath(locale, `/kids/${mode}?child_profile_id=${encodeURIComponent(profile.profile_id)}`);
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

  return (
    <a className={styles.episodeCard} href={href}>
      <img alt="" loading="lazy" src={item.episode.thumbnail_url} />
      <div className={styles.episodeBody}>
        {typeof rank === "number" ? <span className={styles.rank}>#{rank + 1}</span> : null}
        <h3>{item.episode.title}</h3>
        <p>{item.episode.summary}</p>
        <div className={styles.cardMeta}>
          <strong>
            {scoreLabel}: {item.score}
          </strong>
          <span>{item.episode.learning_tags.join(" · ")}</span>
        </div>
      </div>
    </a>
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

  return (
    <main className={styles.page}>
      <div aria-hidden="true" className={styles.noise} />

      <header className={styles.topBar}>
        <a className={styles.brand} href={withLocalePath(locale, "/")}>
          {messages.common.brand}
        </a>

        <nav aria-label="Main navigation" className={styles.navLinks}>
          <a href={withLocalePath(locale, "/parents/onboarding")}>{messages.common.navFamilySetup}</a>
          <a href={withLocalePath(locale, "/parents")}>{messages.common.navParents}</a>
          <a href={withLocalePath(locale, "/admin/studio")}>{messages.common.navStudio}</a>
        </nav>

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
        <div className={styles.heroVisual}>
          {featured ? <img alt="" className={styles.heroImage} src={featured.thumbnail_url} /> : null}
          <div className={styles.heroShade} />

          <div className={styles.heroContent}>
            <span className={styles.heroTag}>{messages.home.heroTag}</span>
            <h1>{featured?.title ?? messages.common.brand}</h1>
            <p>{featured?.summary ?? fallbackSummary}</p>

            <div className={styles.heroActions}>
              <a className={styles.primaryButton} href={modeHref("core", hubInputs, locale)}>
                {messages.home.playNow}
              </a>
              <a className={styles.secondaryButton} href={withLocalePath(locale, "/parents/onboarding")}>
                {messages.home.startFamilySetup}
              </a>
            </div>

            {featured ? (
              <ul className={styles.actionChips}>
                {featured.learning_tags.slice(0, 4).map((tag) => (
                  <li key={tag}>{tag.replaceAll("_", " ")}</li>
                ))}
              </ul>
            ) : null}
          </div>
        </div>

        <aside className={styles.heroPanel}>
          <h2>{messages.home.sections.trust}</h2>
          <ul>
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
          <a href={withLocalePath(locale, "/parents")}>{messages.common.openParentDashboard}</a>
        </aside>
      </section>

      <section className={styles.modeDeck}>
        <article className={`${styles.modeCard} ${styles.mode_early}`}>
          <span className={styles.modeChip}>
            {messages.common.profileAgePrefix} 3-5
          </span>
          <h3>{messages.home.modeDeck.earlyTitle}</h3>
          <p>{messages.home.modeDeck.earlyDetail}</p>
          <a href={modeHref("early", hubInputs, locale)}>{messages.common.openMode}</a>
        </article>

        <article className={`${styles.modeCard} ${styles.mode_core}`}>
          <span className={styles.modeChip}>
            {messages.common.profileAgePrefix} 6-11
          </span>
          <h3>{messages.home.modeDeck.coreTitle}</h3>
          <p>{messages.home.modeDeck.coreDetail}</p>
          <a href={modeHref("core", hubInputs, locale)}>{messages.common.openMode}</a>
        </article>

        <article className={`${styles.modeCard} ${styles.mode_teen}`}>
          <span className={styles.modeChip}>
            {messages.common.profileAgePrefix} 12-16
          </span>
          <h3>{messages.home.modeDeck.teenTitle}</h3>
          <p>{messages.home.modeDeck.teenDetail}</p>
          <a href={modeHref("teen", hubInputs, locale)}>{messages.common.openMode}</a>
        </article>
      </section>

      <section className={styles.section}>
        <div className={styles.sectionHead}>
          <h2>{messages.home.sections.chooseProfile}</h2>
          <a href={withLocalePath(locale, "/parents/onboarding")}>{messages.common.createProfile}</a>
        </div>

        <div className={styles.profileGrid}>
          {hub.profiles.map((profile) => (
            <a
              className={`${styles.profileCard} ${styles[`profile_${profile.mode}`]}`}
              href={withLocalePath(
                locale,
                `/kids/${profile.mode}?child_profile_id=${encodeURIComponent(profile.profile_id)}`
              )}
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

      <section className={styles.section}>
        <div className={styles.sectionHead}>
          <h2>{messages.home.sections.continueWatching}</h2>
          <span>{messages.home.sections.continueWatchingSub}</span>
        </div>

        <div className={styles.rail}>
          {hub.continueWatching.length === 0 ? (
            <article className={styles.emptyCard}>{messages.home.labels.noContinueWatching}</article>
          ) : (
            hub.continueWatching.map((item) => (
              <a className={styles.continueCard} href={withLocalePath(locale, item.href)} key={item.childProfileID}>
                <img alt="" loading="lazy" src={item.episode.thumbnail_url} />
                <div className={styles.continueBody}>
                  <span>{item.profileName}</span>
                  <h3>{item.episode.title}</h3>
                  <p>
                    {item.watchedMinutesToday} min today · {item.completionPercent}%
                  </p>
                </div>
              </a>
            ))
          )}
        </div>
      </section>

      <section className={styles.section}>
        <div className={styles.sectionHead}>
          <h2>{messages.home.sections.top10}</h2>
          <span>{messages.home.sections.top10Sub}</span>
        </div>

        <div className={styles.rail}>
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

      <section className={styles.section}>
        <div className={styles.sectionHead}>
          <h2>{messages.home.sections.forYou}</h2>
          <span>{messages.home.sections.forYouSub}</span>
        </div>
        <div className={styles.rail}>
          {hub.categoryRows.forYou.map((item) => (
            <RankedCard
              inputs={hubInputs}
              item={{
                episode: item,
                score: Math.round(item.age_fit_score * 100),
                scoreBreakdown: { ageFit: Math.round(item.age_fit_score * 100), reasonBonus: 0, safetyBonus: 0 },
                sourceModes: [item.content_suitability]
              }}
              key={`${item.episode_id}-for-you`}
              locale={locale}
              scoreLabel={messages.home.labels.topRankLabel}
            />
          ))}
        </div>
      </section>

      <section className={styles.section}>
        <div className={styles.sectionHead}>
          <h2>{messages.home.sections.knowledge}</h2>
          <span>{messages.home.sections.knowledgeSub}</span>
        </div>
        <div className={styles.rail}>
          {hub.categoryRows.knowledge.map((item) => (
            <RankedCard
              inputs={hubInputs}
              item={{
                episode: item,
                score: Math.round(item.age_fit_score * 100),
                scoreBreakdown: { ageFit: Math.round(item.age_fit_score * 100), reasonBonus: 0, safetyBonus: 0 },
                sourceModes: [item.content_suitability]
              }}
              key={`${item.episode_id}-knowledge`}
              locale={locale}
              scoreLabel={messages.home.labels.topRankLabel}
            />
          ))}
        </div>
      </section>

      <section className={styles.section}>
        <div className={styles.sectionHead}>
          <h2>{messages.home.sections.creative}</h2>
          <span>{messages.home.sections.creativeSub}</span>
        </div>
        <div className={styles.rail}>
          {hub.categoryRows.creative.map((item) => (
            <RankedCard
              inputs={hubInputs}
              item={{
                episode: item,
                score: Math.round(item.age_fit_score * 100),
                scoreBreakdown: { ageFit: Math.round(item.age_fit_score * 100), reasonBonus: 0, safetyBonus: 0 },
                sourceModes: [item.content_suitability]
              }}
              key={`${item.episode_id}-creative`}
              locale={locale}
              scoreLabel={messages.home.labels.topRankLabel}
            />
          ))}
        </div>
      </section>

      <section className={styles.section}>
        <div className={styles.sectionHead}>
          <h2>{messages.home.sections.adventure}</h2>
          <span>{messages.home.sections.adventureSub}</span>
        </div>
        <div className={styles.rail}>
          {hub.categoryRows.adventure.map((item) => (
            <RankedCard
              inputs={hubInputs}
              item={{
                episode: item,
                score: Math.round(item.age_fit_score * 100),
                scoreBreakdown: { ageFit: Math.round(item.age_fit_score * 100), reasonBonus: 0, safetyBonus: 0 },
                sourceModes: [item.content_suitability]
              }}
              key={`${item.episode_id}-adventure`}
              locale={locale}
              scoreLabel={messages.home.labels.topRankLabel}
            />
          ))}
        </div>
      </section>

      <section className={styles.section}>
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

      <footer className={styles.footer}>
        <p>{messages.common.brand} · Family streaming, safe by default</p>
        <div>
          <a href={withLocalePath(locale, "/parents")}>{messages.common.navParents}</a>
          <a href={withLocalePath(locale, "/admin/studio")}>{messages.common.navStudio}</a>
        </div>
      </footer>
    </main>
  );
}
