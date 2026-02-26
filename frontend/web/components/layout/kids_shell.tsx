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

  return (
    <div className={`${styles.page} ${modeClassName[mode]}`}>
      <header className={styles.top}>
        <div>
          <p className={styles.kicker}>{messages.shell.kicker}</p>
          <h1>{messages.shell.modeTitle[mode]}</h1>
          <p className={styles.note}>{messages.shell.note}</p>
        </div>
        <ProfileOrb ageBand={ageBand} name={profileName} subtitle={subtitle} />
      </header>

      <KidsNav activeKey={activeNavKey} items={navItems} locale={locale} />
      <SessionMeter labels={messages.cards} sessionLimitMinutes={sessionLimitMinutes} watchedMinutes={watchedMinutes} />

      <main>{children}</main>
    </div>
  );
}
