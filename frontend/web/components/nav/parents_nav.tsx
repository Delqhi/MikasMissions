import styles from "./parents_nav.module.css";

type ParentsNavProps = {
  active: string;
};

type ParentNavItem = {
  label: string;
  href: string;
};

const tabs: ParentNavItem[] = [
  { label: "Dashboard", href: "/parents#dashboard" },
  { label: "Controls", href: "/parents#controls" },
  { label: "Compliance", href: "/parents#compliance" },
  { label: "Onboarding", href: "/parents/onboarding" }
];

export function ParentsNav({ active }: ParentsNavProps) {
  return (
    <nav aria-label="Parents navigation" className={styles.nav}>
      {tabs.map((tab) => (
        <a
          aria-current={tab.label === active ? "page" : undefined}
          className={tab.label === active ? styles.active : styles.idle}
          href={tab.href}
          key={tab.label}
        >
          {tab.label}
        </a>
      ))}
    </nav>
  );
}
