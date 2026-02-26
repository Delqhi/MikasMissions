import styles from "./big_action_button.module.css";

type BigActionButtonProps = {
  label: string;
  hint?: string;
  href?: string;
  variant?: "primary" | "secondary";
};

export function BigActionButton(props: BigActionButtonProps) {
  const { label, hint, href, variant = "primary" } = props;

  if (href) {
    return (
      <a className={`${styles.button} ${styles[variant]}`} href={href}>
        <span className={styles.label}>{label}</span>
        {hint ? <span className={styles.hint}>{hint}</span> : null}
      </a>
    );
  }

  return (
    <button className={`${styles.button} ${styles[variant]}`} type="button">
      <span className={styles.label}>{label}</span>
      {hint ? <span className={styles.hint}>{hint}</span> : null}
    </button>
  );
}
