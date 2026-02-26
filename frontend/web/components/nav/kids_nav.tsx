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

export function KidsNav({ items, activeKey, locale }: KidsNavProps) {
  return (
    <nav aria-label="Kids navigation" className={styles.nav}>
      {items.map((item) => (
        <a
          aria-current={item.key === activeKey ? "page" : undefined}
          className={item.key === activeKey ? styles.active : styles.idle}
          href={withLocalePath(locale, item.href)}
          key={`${item.key}-${item.href}`}
        >
          {item.label}
        </a>
      ))}
    </nav>
  );
}
