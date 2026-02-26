import styles from "./profile_orb.module.css";

type ProfileOrbProps = {
  name: string;
  ageBand: string;
  subtitle: string;
};

function initialsFromName(name: string): string {
  const tokens = name.trim().split(/\s+/);
  const letters = tokens.slice(0, 2).map((token) => token[0]?.toUpperCase() ?? "");
  return letters.join("") || "MM";
}

export function ProfileOrb({ name, ageBand, subtitle }: ProfileOrbProps) {
  return (
    <article className={styles.card} aria-label={`${name} profile`}>
      <div className={styles.orb} aria-hidden="true">
        {initialsFromName(name)}
      </div>
      <div className={styles.meta}>
        <h3>{name}</h3>
        <p>{subtitle}</p>
        <span>{ageBand}</span>
      </div>
    </article>
  );
}
