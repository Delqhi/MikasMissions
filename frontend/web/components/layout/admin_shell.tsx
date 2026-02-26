import type { ReactNode } from "react";
import styles from "./admin_shell.module.css";

type AdminNavItem = {
  label: string;
  href: string;
};

const navItems: AdminNavItem[] = [
  { label: "Studio", href: "/admin/studio" },
  { label: "Runs", href: "/admin/runs" },
  { label: "Login", href: "/admin/login" }
];

type AdminShellProps = {
  title: string;
  subtitle: string;
  activeNav: "Studio" | "Runs" | "Login";
  children: ReactNode;
};

export function AdminShell({ title, subtitle, activeNav, children }: AdminShellProps) {
  return (
    <div className={styles.page}>
      <div aria-hidden="true" className={styles.backdrop} />

      <header className={styles.header}>
        <div>
          <p className={styles.kicker}>Admin Console</p>
          <h1>{title}</h1>
          <span>{subtitle}</span>
        </div>

        <div className={styles.actions}>
          <a className={styles.homeLink} href="/">
            Home
          </a>
        </div>
      </header>

      <nav aria-label="Admin navigation" className={styles.nav}>
        {navItems.map((item) => (
          <a
            aria-current={item.label === activeNav ? "page" : undefined}
            className={item.label === activeNav ? styles.active : styles.idle}
            href={item.href}
            key={item.label}
          >
            {item.label}
          </a>
        ))}
      </nav>

      <main className={styles.main}>{children}</main>
    </div>
  );
}
