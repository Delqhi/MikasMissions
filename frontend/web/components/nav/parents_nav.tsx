import type { Locale } from "../../lib/i18n";
import { withLocalePath } from "../../lib/i18n";
import type { LocalizedMessages } from "../../lib/messages/types";
import styles from "./parents_nav.module.css";

type ParentsNavProps = {
  active: ParentNavKey;
  locale: Locale;
  labels: LocalizedMessages["parents"]["nav"];
};

export type ParentNavKey = "dashboard" | "controls" | "compliance" | "onboarding";

type ParentNavItem = {
  key: ParentNavKey;
  label: string;
  href: string;
};

export function ParentsNav({ active, locale, labels }: ParentsNavProps) {
  const tabs: ParentNavItem[] = [
    { key: "dashboard", label: labels.dashboard, href: "/parents#dashboard" },
    { key: "controls", label: labels.controls, href: "/parents#controls" },
    { key: "compliance", label: labels.compliance, href: "/parents#compliance" },
    { key: "onboarding", label: labels.onboarding, href: "/parents/onboarding" }
  ];

  return (
    <nav aria-label="Parents navigation" className={styles.nav}>
      {tabs.map((tab) => (
        <a
          aria-current={tab.key === active ? "page" : undefined}
          className={tab.key === active ? styles.active : styles.idle}
          href={withLocalePath(locale, tab.href)}
          key={tab.key}
        >
          {tab.label}
        </a>
      ))}
    </nav>
  );
}
