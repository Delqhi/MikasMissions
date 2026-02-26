import { formatDurationMS } from "../../lib/formatters";
import type { RailItem } from "../../lib/experience_types";
import type { LocalizedMessages } from "../../lib/messages/types";
import { LearningBadge } from "./learning_badge";
import styles from "./episode_tile.module.css";

type EpisodeTileProps = {
  episode: RailItem;
  labels?: Pick<LocalizedMessages["kids"]["cards"], "episodePrefix">;
};

const defaultLabels = {
  episodePrefix: "Episode"
} as const;

export function EpisodeTile({ episode, labels }: EpisodeTileProps) {
  const text = labels ?? defaultLabels;
  const tags = Array.isArray(episode.learning_tags) ? episode.learning_tags : [];

  return (
    <article className={styles.tile} aria-label={`${text.episodePrefix} ${episode.title}`}>
      <img alt="" className={styles.thumb} loading="lazy" src={episode.thumbnail_url} />
      <div className={styles.body}>
        <div className={styles.topline}>
          <h3>{episode.title}</h3>
          <span>{formatDurationMS(episode.duration_ms)}</span>
        </div>
        <p>{episode.summary}</p>
        <div className={styles.badges}>
          {tags.map((tag) => (
            <LearningBadge key={`${episode.episode_id}-${tag}`} label={tag} />
          ))}
        </div>
      </div>
    </article>
  );
}
