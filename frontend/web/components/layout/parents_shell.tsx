import type { ReactNode } from "react";
import type { Locale } from "../../lib/i18n";
import { withLocalePath } from "../../lib/i18n";
import type { LocalizedMessages } from "../../lib/messages/types";
import { ParentsNav, type ParentNavKey } from "../nav/parents_nav";
import styles from "./parents_shell.module.css";

type ParentsShellProps = {
  locale: Locale;
  messages: LocalizedMessages["parents"];
  heading: string;
  description?: string;
  activeNav?: ParentNavKey;
  children: ReactNode;
};

export function ParentsShell({
  locale,
  messages,
  heading,
  description = messages.shell.defaultDescription,
  activeNav = "dashboard",
  children
}: ParentsShellProps) {
  return (
    <div className={`${styles.page} mode-parent`}>
      <div aria-hidden="true" className={styles.backdrop} />

      <header className={styles.header}>
        <div>
          <p className={styles.kicker}>{messages.shell.kicker}</p>
          <h1>{heading}</h1>
          <span>{description}</span>
        </div>
        <div className={styles.actions}>
          <a className={styles.primaryLink} href={withLocalePath(locale, "/parents/onboarding")}>
            {messages.shell.addChildProfile}
          </a>
          <a className={styles.switchLink} href={withLocalePath(locale, "/")}>
            {messages.shell.switchProfile}
          </a>
        </div>
      </header>

      <ParentsNav active={activeNav} labels={messages.nav} locale={locale} />

      <main className={styles.main}>{children}</main>
    </div>
  );
}
