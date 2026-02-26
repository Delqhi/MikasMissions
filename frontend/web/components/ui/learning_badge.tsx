import styles from "./learning_badge.module.css";

type LearningBadgeProps = {
  label: string;
};

export function LearningBadge({ label }: LearningBadgeProps) {
  return <span className={styles.badge}>{label.replaceAll("_", " ")}</span>;
}
