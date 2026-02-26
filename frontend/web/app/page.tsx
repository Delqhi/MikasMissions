import { demoProfiles, kidsHomeFallback } from "../lib/mock_payloads";
import styles from "./page.module.css";

const featured = kidsHomeFallback.core.rails[0];
const controlHighlights = [
  { label: "Age modes", value: "3 adaptive worlds" },
  { label: "Safety", value: "Strict by default" },
  { label: "Parent controls", value: "Realtime updates" }
];

const rows = [
  {
    title: "Top picks this week",
    subtitle: "Most started missions",
    items: [...kidsHomeFallback.core.rails, ...kidsHomeFallback.teen.rails]
  },
  {
    title: "Early learning lane",
    subtitle: "Guided stories for mini explorers",
    items: kidsHomeFallback.early.rails
  },
  {
    title: "Teen spotlight",
    subtitle: "Creative and critical thinking missions",
    items: kidsHomeFallback.teen.rails
  }
];

const trustRows = [
  {
    title: "Filter coverage",
    detail: "Every recommendation rail ships with safety metadata and transparent reason codes."
  },
  {
    title: "Session limits",
    detail: "Daily watch caps are enforced and visible across all age modes without manual checks."
  },
  {
    title: "Parent gate",
    detail: "External links, purchases, and account changes always require adult verification."
  }
];

export default function HomePage() {
  return (
    <main className={styles.page}>
      <div aria-hidden="true" className={styles.texture} />

      <header className={styles.topBar}>
        <a className={styles.brand} href="/">
          MIKASMISSIONS
        </a>
        <nav className={styles.navLinks} aria-label="Main navigation">
          <a href="/parents/onboarding">Family setup</a>
          <a href="/parents">Parents</a>
          <a href="/admin/studio">Studio</a>
        </nav>
      </header>

      <section className={styles.hero} aria-label="Featured mission">
        <img alt="" className={styles.heroImage} src={featured.thumbnail_url} />
        <div className={styles.heroOverlay} />

        <div className={styles.heroContent}>
          <span className={styles.heroTag}>Heute im Fokus</span>
          <h1>{featured.title}</h1>
          <p>{featured.summary}</p>

          <div className={styles.heroActions}>
            <a className={styles.primaryButton} href="/kids/core?child_profile_id=child-core-01">
              Play now
            </a>
            <a className={styles.secondaryButton} href="/parents/onboarding">
              Start family setup
            </a>
          </div>

          <ul className={styles.chipRow}>
            {kidsHomeFallback.core.primary_actions.slice(0, 5).map((action) => (
              <li key={action}>{action}</li>
            ))}
          </ul>
        </div>

        <aside className={styles.heroPanel}>
          <h2>Safety cockpit</h2>
          <ul>
            {controlHighlights.map((item) => (
              <li key={item.label}>
                <span>{item.label}</span>
                <strong>{item.value}</strong>
              </li>
            ))}
          </ul>
          <a href="/parents">Open parent dashboard</a>
        </aside>
      </section>

      <section className={styles.section}>
        <div className={styles.sectionHead}>
          <h2>Who is watching?</h2>
          <a href="/parents/onboarding">Create profile</a>
        </div>

        <div className={styles.profileGrid}>
          {demoProfiles.map((profile) => (
            <a
              key={profile.profile_id}
              className={styles.profileCard}
              href={`${profile.href}?child_profile_id=${profile.profile_id}`}
            >
              <span className={styles.profileMode}>{profile.mode}</span>
              <strong>{profile.name}</strong>
              <p>{profile.subtitle}</p>
              <span className={styles.profileMeta}>Age {profile.age_band}</span>
            </a>
          ))}
        </div>
      </section>

      {rows.map((row) => (
        <section className={styles.section} key={row.title}>
          <div className={styles.sectionHead}>
            <h2>{row.title}</h2>
            <span>{row.subtitle}</span>
          </div>

          <div className={styles.rail}>
            {row.items.map((item, index) => (
              <article key={`${row.title}-${item.episode_id}`} className={styles.episodeCard}>
                <img alt="" src={item.thumbnail_url} />
                <div className={styles.episodeBody}>
                  <span className={styles.rank}>#{index + 1}</span>
                  <h3>{item.title}</h3>
                  <p>{item.summary}</p>
                  <em>{item.learning_tags.join(" · ")}</em>
                </div>
              </article>
            ))}
          </div>
        </section>
      ))}

      <section className={styles.section}>
        <div className={styles.sectionHead}>
          <h2>Built for trust</h2>
          <span>Parent-first guardrails with kid-friendly exploration</span>
        </div>

        <div className={styles.trustGrid}>
          {trustRows.map((item) => (
            <article className={styles.trustCard} key={item.title}>
              <h3>{item.title}</h3>
              <p>{item.detail}</p>
            </article>
          ))}
        </div>
      </section>

      <footer className={styles.footer}>
        <p>MikasMissions · Safe streaming for families</p>
        <div>
          <a href="/parents">Parent controls</a>
          <a href="/admin/studio">Admin studio</a>
        </div>
      </footer>
    </main>
  );
}
