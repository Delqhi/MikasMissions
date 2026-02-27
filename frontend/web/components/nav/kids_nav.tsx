import type { Locale } from "../../lib/i18n";
import { withLocalePath } from "../../lib/i18n";
import styles from "./kids_nav.module.css";

export type KidsNavItem = {
  key: string;
  label: string;
  href: string;
};

type KidsNavProps = {
  items: KidsNavItem[];
  activeKey: string;
  locale: Locale;
};

function isParentEntry(item: KidsNavItem): boolean {
  return item.key.toLowerCase().includes("parent") || item.href.toLowerCase().includes("/parents");
}

export function KidsNav({ items, activeKey, locale }: KidsNavProps) {
  return (
    <nav aria-label="Kids navigation" className={styles.nav}>
      {items.map((item) => {
        const active = item.key === activeKey;
        const parentEntry = isParentEntry(item);
        const className = [active ? styles.active : styles.idle, parentEntry ? styles.parentEntry : ""]
          .filter(Boolean)
          .join(" ");

        return (
          <a
            aria-current={active ? "page" : undefined}
            className={className}
            data-parent-entry={parentEntry ? "true" : undefined}
            href={withLocalePath(locale, item.href)}
            key={`${item.key}-${item.href}`}
          >
            {item.label}
          </a>
        );
      })}
    </nav>
  );
}
