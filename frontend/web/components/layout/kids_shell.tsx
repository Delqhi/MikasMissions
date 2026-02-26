import type { ReactNode } from "react";
import type { Locale } from "../../lib/i18n";
import type { KidsMode } from "../../lib/experience_types";
import type { LocalizedMessages } from "../../lib/messages/types";
import { modeClassName } from "../../lib/mode_theme";
import { KidsNav, type KidsNavItem } from "../nav/kids_nav";
import { ProfileOrb } from "../ui/profile_orb";
import { SessionMeter } from "../ui/session_meter";
import styles from "./kids_shell.module.css";

type KidsShellProps = {
  locale: Locale;
  messages: LocalizedMessages["kids"];
  mode: KidsMode;
  profileName: string;
  ageBand: string;
  subtitle: string;
  navItems: KidsNavItem[];
  activeNavKey: string;
  watchedMinutes: number;
  sessionLimitMinutes: number;
  children: ReactNode;
};

export function KidsShell(props: KidsShellProps) {
  const {
    locale,
    messages,
    mode,
    profileName,
    ageBand,
    subtitle,
    navItems,
    activeNavKey,
    watchedMinutes,
    sessionLimitMinutes,
    children
  } = props;
  const remainingMinutes = Math.max(0, sessionLimitMinutes - watchedMinutes);
  const modeBandLabel = mode === "early" ? messages.switcher.early : mode === "teen" ? messages.switcher.teen : messages.switcher.core;

  return (
    <div className={`${styles.page} ${modeClassName[mode]}`}>
      <a className={styles.skipLink} href="#kids-main">
        Skip to content
      </a>

      <header className={styles.top}>
        <div className={styles.headline}>
          <p className={styles.kicker}>{messages.shell.kicker}</p>
          <h1>{messages.shell.modeTitle[mode]}</h1>
          <p className={styles.note}>{messages.shell.note}</p>
          <ul className={styles.statusRow}>
            <li>
              <span>{messages.cards.sessionCap}</span>
              <strong>
                {remainingMinutes} {messages.cards.minLeftSuffix}
              </strong>
            </li>
            <li>
              <span>{profileName}</span>
              <strong>{modeBandLabel}</strong>
            </li>
          </ul>
        </div>
        <ProfileOrb ageBand={ageBand} name={profileName} subtitle={subtitle} />
      </header>

      <section className={styles.navTray}>
        <KidsNav activeKey={activeNavKey} items={navItems} locale={locale} />
      </section>
      <SessionMeter labels={messages.cards} sessionLimitMinutes={sessionLimitMinutes} watchedMinutes={watchedMinutes} />

      <main className={styles.main} id="kids-main">
        {children}
      </main>
    </div>
  );
}
