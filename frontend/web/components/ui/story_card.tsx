import type { CSSProperties } from "react";
import type { RailItem } from "../../lib/experience_types";
import { BigActionButton } from "./big_action_button";
import { LearningBadge } from "./learning_badge";
import styles from "./story_card.module.css";

type StoryCardProps = {
  item: RailItem;
  index: number;
};

export function StoryCard({ item, index }: StoryCardProps) {
  const delayStyle = {
    "--stagger-delay": `${index * 70}ms`
  } as CSSProperties;

  return (
    <article className={styles.card} style={delayStyle}>
      <img alt="" className={styles.hero} loading="lazy" src={item.thumbnail_url} />
      <div className={styles.overlay}>
        <p className={styles.reason}>{item.reason_code.replaceAll("_", " ")}</p>
        <h3>{item.title}</h3>
        <p className={styles.summary}>{item.summary}</p>
        <div className={styles.badges}>
          {item.learning_tags.map((tag) => (
            <LearningBadge key={`${item.episode_id}-${tag}`} label={tag} />
          ))}
        </div>
        <div className={styles.foot}>
          <BigActionButton hint={item.content_suitability} label="Open episode" />
          <dl>
            <div>
              <dt>Age fit</dt>
              <dd>{Math.round(item.age_fit_score * 100)}%</dd>
            </div>
            <div>
              <dt>Safety</dt>
              <dd>{item.safety_applied ? "Applied" : "Off"}</dd>
            </div>
          </dl>
        </div>
      </div>
    </article>
  );
}
