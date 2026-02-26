import type { ReactNode } from "react";
import type { KidsMode } from "../../lib/experience_types";
import { modeClassName, modeLabel } from "../../lib/mode_theme";
import { KidsNav, type KidsNavItem } from "../nav/kids_nav";
import { ProfileOrb } from "../ui/profile_orb";
import { SessionMeter } from "../ui/session_meter";
import styles from "./kids_shell.module.css";

type KidsShellProps = {
  mode: KidsMode;
  profileName: string;
  ageBand: string;
  subtitle: string;
  navItems: KidsNavItem[];
  activeNav: string;
  watchedMinutes: number;
  sessionLimitMinutes: number;
  children: ReactNode;
};

export function KidsShell(props: KidsShellProps) {
  const {
    mode,
    profileName,
    ageBand,
    subtitle,
    navItems,
    activeNav,
    watchedMinutes,
    sessionLimitMinutes,
    children
  } = props;

  return (
    <div className={`${styles.page} ${modeClassName[mode]}`}>
      <header className={styles.top}>
        <div>
          <p className={styles.kicker}>MikasMissions</p>
          <h1>{modeLabel[mode]}</h1>
          <p className={styles.note}>Strict safety defaults and age-fit rails are always active.</p>
        </div>
        <ProfileOrb ageBand={ageBand} name={profileName} subtitle={subtitle} />
      </header>

      <KidsNav activeItem={activeNav} items={navItems} />
      <SessionMeter sessionLimitMinutes={sessionLimitMinutes} watchedMinutes={watchedMinutes} />

      <main>{children}</main>
    </div>
  );
}
