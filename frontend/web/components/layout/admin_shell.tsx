import type { ReactNode } from "react";
import type { Locale } from "../../lib/i18n";
import { withLocalePath } from "../../lib/i18n";
import type { LocalizedMessages } from "../../lib/messages/types";
import styles from "./admin_shell.module.css";

type AdminNavItem = {
  key: "studio" | "runs" | "login";
  label: string;
  href: string;
};

type AdminShellProps = {
  locale: Locale;
  labels: LocalizedMessages["admin"]["shell"];
  title: string;
  subtitle: string;
  activeNav: AdminNavItem["key"];
  children: ReactNode;
};

export function AdminShell({ locale, labels, title, subtitle, activeNav, children }: AdminShellProps) {
  const navItems: AdminNavItem[] = [
    { key: "studio", label: labels.nav.studio, href: "/admin/studio" },
    { key: "runs", label: labels.nav.runs, href: "/admin/runs" },
    { key: "login", label: labels.nav.login, href: "/admin/login" }
  ];

  return (
    <div className={styles.page}>
      <div aria-hidden="true" className={styles.backdrop} />

      <header className={styles.header}>
        <div>
          <p className={styles.kicker}>{labels.kicker}</p>
          <h1>{title}</h1>
          <span>{subtitle}</span>
        </div>

        <div className={styles.actions}>
          <a className={styles.homeLink} href={withLocalePath(locale, "/")}>
            {labels.home}
          </a>
        </div>
      </header>

      <nav aria-label="Admin navigation" className={styles.nav}>
        {navItems.map((item) => (
          <a
            aria-current={item.key === activeNav ? "page" : undefined}
            className={item.key === activeNav ? styles.active : styles.idle}
            href={withLocalePath(locale, item.href)}
            key={item.key}
          >
            {item.label}
          </a>
        ))}
      </nav>

      <main className={styles.main}>{children}</main>
    </div>
  );
}
