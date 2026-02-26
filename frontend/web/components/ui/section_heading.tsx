import styles from "./section_heading.module.css";

type SectionHeadingProps = {
  title: string;
  description: string;
};

export function SectionHeading({ title, description }: SectionHeadingProps) {
  return (
    <header className={styles.header}>
      <h2>{title}</h2>
      <p>{description}</p>
    </header>
  );
}
