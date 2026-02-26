import type { CSSProperties } from "react";
import type { RailItem } from "../../lib/experience_types";
import type { LocalizedMessages } from "../../lib/messages/types";
import { BigActionButton } from "./big_action_button";
import { LearningBadge } from "./learning_badge";
import styles from "./story_card.module.css";

type StoryCardProps = {
  item: RailItem;
  index: number;
  labels?: Pick<LocalizedMessages["kids"]["cards"], "openEpisode" | "ageFit" | "safety" | "safetyApplied" | "safetyOff">;
};

const defaultLabels = {
  openEpisode: "Open episode",
  ageFit: "Age fit",
  safety: "Safety",
  safetyApplied: "Applied",
  safetyOff: "Off"
} as const;

export function StoryCard({ item, index, labels }: StoryCardProps) {
  const delayStyle = {
    "--stagger-delay": `${index * 70}ms`
  } as CSSProperties;
  const text = labels ?? defaultLabels;
  const reasonLabel = typeof item.reason_code === "string" ? item.reason_code.replaceAll("_", " ") : "recommended";
  const tags = Array.isArray(item.learning_tags) ? item.learning_tags : [];
  const ageFitPercent = Math.round((Number(item.age_fit_score) || 0) * 100);

  return (
    <article className={styles.card} style={delayStyle}>
      <img alt="" className={styles.hero} loading="lazy" src={item.thumbnail_url} />
      <div className={styles.overlay}>
        <p className={styles.reason}>{reasonLabel}</p>
        <h3>{item.title}</h3>
        <p className={styles.summary}>{item.summary}</p>
        <div className={styles.badges}>
          {tags.map((tag) => (
            <LearningBadge key={`${item.episode_id}-${tag}`} label={tag} />
          ))}
        </div>
        <div className={styles.foot}>
          <BigActionButton hint={item.content_suitability} label={text.openEpisode} />
          <dl>
            <div>
              <dt>{text.ageFit}</dt>
              <dd>{ageFitPercent}%</dd>
            </div>
            <div>
              <dt>{text.safety}</dt>
              <dd>{item.safety_applied ? text.safetyApplied : text.safetyOff}</dd>
            </div>
          </dl>
        </div>
      </div>
    </article>
  );
}
