import { formatDurationMS } from "../../lib/formatters";
import type { RailItem } from "../../lib/experience_types";
import { LearningBadge } from "./learning_badge";
import styles from "./episode_tile.module.css";

type EpisodeTileProps = {
  episode: RailItem;
};

export function EpisodeTile({ episode }: EpisodeTileProps) {
  return (
    <article className={styles.tile} aria-label={`Episode ${episode.title}`}>
      <img alt="" className={styles.thumb} loading="lazy" src={episode.thumbnail_url} />
      <div className={styles.body}>
        <div className={styles.topline}>
          <h3>{episode.title}</h3>
          <span>{formatDurationMS(episode.duration_ms)}</span>
        </div>
        <p>{episode.summary}</p>
        <div className={styles.badges}>
          {episode.learning_tags.map((tag) => (
            <LearningBadge key={`${episode.episode_id}-${tag}`} label={tag} />
          ))}
        </div>
      </div>
    </article>
  );
}
