import type { ReactNode } from "react";
import { ParentsNav } from "../nav/parents_nav";
import styles from "./parents_shell.module.css";

type ParentsShellProps = {
  heading: string;
  description?: string;
  activeNav?: string;
  children: ReactNode;
};

export function ParentsShell({
  heading,
  description = "Consent, controls, reports, and emergency overrides in one place.",
  activeNav = "Dashboard",
  children
}: ParentsShellProps) {
  return (
    <div className={`${styles.page} mode-parent`}>
      <div aria-hidden="true" className={styles.backdrop} />

      <header className={styles.header}>
        <div>
          <p className={styles.kicker}>Parent Command</p>
          <h1>{heading}</h1>
          <span>{description}</span>
        </div>
        <div className={styles.actions}>
          <a className={styles.primaryLink} href="/parents/onboarding">
            Add child profile
          </a>
          <a className={styles.switchLink} href="/">
            Switch profile
          </a>
        </div>
      </header>

      <ParentsNav active={activeNav} />

      <main className={styles.main}>{children}</main>
    </div>
  );
}
