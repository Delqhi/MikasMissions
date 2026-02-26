import styles from "./kids_nav.module.css";

export type KidsNavItem = {
  label: string;
  href: string;
};

type KidsNavProps = {
  items: KidsNavItem[];
  activeItem: string;
};

export function KidsNav({ items, activeItem }: KidsNavProps) {
  return (
    <nav aria-label="Kids navigation" className={styles.nav}>
      {items.map((item) => (
        <a
          aria-current={item.label === activeItem ? "page" : undefined}
          className={item.label === activeItem ? styles.active : styles.idle}
          href={item.href}
          key={`${item.label}-${item.href}`}
        >
          {item.label}
        </a>
      ))}
    </nav>
  );
}
