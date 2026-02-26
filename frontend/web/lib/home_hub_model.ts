import type { KidsHomeResponse, KidsMode, KidsProgressResponse, ProfileCard, RailItem } from "./experience_types";

export type ContinueWatchingItem = {
  childProfileID: string;
  profileName: string;
  ageBand: ProfileCard["age_band"];
  mode: KidsMode;
  href: string;
  episode: RailItem;
  watchedMinutesToday: number;
  completionPercent: number;
};

export type TopRankedItem = {
  episode: RailItem;
  score: number;
  scoreBreakdown: {
    ageFit: number;
    safetyBonus: number;
    reasonBonus: number;
  };
  sourceModes: KidsMode[];
};

export type HomeHubViewModel = {
  featured: RailItem | null;
  profiles: ProfileCard[];
  continueWatching: ContinueWatchingItem[];
  top10: TopRankedItem[];
  categoryRows: {
    forYou: RailItem[];
    knowledge: RailItem[];
    creative: RailItem[];
    adventure: RailItem[];
  };
};

export type HomeHubInput = {
  profile: ProfileCard;
  mode: KidsMode;
  home: KidsHomeResponse;
  progress: KidsProgressResponse;
};

function containsAny(value: string, keywords: string[]): boolean {
  const normalized = value.toLowerCase();
  return keywords.some((keyword) => normalized.includes(keyword));
}

function episodeScore(episode: RailItem): TopRankedItem["scoreBreakdown"] {
  const ageFit = Math.round(episode.age_fit_score * 100);
  const safetyBonus = episode.safety_applied ? 10 : 0;
  const reasonBonus = containsAny(episode.reason_code, ["learning", "progress"]) ? 8 : 0;

  return { ageFit, reasonBonus, safetyBonus };
}

function toTopRankedItem(episode: RailItem, sourceMode: KidsMode): TopRankedItem {
  const scoreBreakdown = episodeScore(episode);
  return {
    episode,
    score: scoreBreakdown.ageFit + scoreBreakdown.safetyBonus + scoreBreakdown.reasonBonus,
    scoreBreakdown,
    sourceModes: [sourceMode]
  };
}

function uniqueByEpisodeID(items: TopRankedItem[]): TopRankedItem[] {
  const deduped = new Map<string, TopRankedItem>();

  for (const item of items) {
    const existing = deduped.get(item.episode.episode_id);

    if (!existing) {
      deduped.set(item.episode.episode_id, item);
      continue;
    }

    if (item.score > existing.score) {
      deduped.set(item.episode.episode_id, {
        ...item,
        sourceModes: Array.from(new Set([...existing.sourceModes, ...item.sourceModes]))
      });
      continue;
    }

    if (item.score === existing.score) {
      deduped.set(item.episode.episode_id, {
        ...existing,
        sourceModes: Array.from(new Set([...existing.sourceModes, ...item.sourceModes]))
      });
    }
  }

  return Array.from(deduped.values());
}

function sortTopRanked(items: TopRankedItem[]): TopRankedItem[] {
  return [...items].sort((left, right) => {
    if (left.score !== right.score) {
      return right.score - left.score;
    }

    if (left.episode.title !== right.episode.title) {
      return left.episode.title.localeCompare(right.episode.title);
    }

    return left.episode.episode_id.localeCompare(right.episode.episode_id);
  });
}

function categorizeRailItems(episodes: RailItem[]): HomeHubViewModel["categoryRows"] {
  const knowledge = episodes.filter((episode) =>
    episode.learning_tags.some((tag) => containsAny(tag, ["science", "learn", "logic", "math", "nature", "knowledge", "history"]))
  );

  const creative = episodes.filter((episode) =>
    episode.learning_tags.some((tag) => containsAny(tag, ["creative", "art", "music", "build", "design", "maker"]))
  );

  const adventure = episodes.filter((episode) =>
    episode.learning_tags.some((tag) => containsAny(tag, ["adventure", "story", "explore", "quest", "team", "mission"]))
  );

  const forYou = episodes;

  return {
    forYou: forYou.slice(0, 12),
    knowledge: knowledge.slice(0, 12),
    creative: creative.slice(0, 12),
    adventure: adventure.slice(0, 12)
  };
}

export function buildHomeHubViewModel(inputs: HomeHubInput[]): HomeHubViewModel {
  const featured = inputs.flatMap((input) => input.home.rails).find(Boolean) ?? null;

  const continueWatching = inputs
    .map((input) => {
      const resume =
        input.home.rails.find((rail) => rail.episode_id === input.progress.last_episode_id) ?? input.home.rails[0] ?? null;

      if (!resume) {
        return null;
      }

      return {
        childProfileID: input.profile.profile_id,
        profileName: input.profile.name,
        ageBand: input.profile.age_band,
        mode: input.mode,
        href: `${input.profile.href}?child_profile_id=${encodeURIComponent(input.profile.profile_id)}`,
        episode: resume,
        watchedMinutesToday: input.progress.watched_minutes_today,
        completionPercent: input.progress.completion_percent
      } satisfies ContinueWatchingItem;
    })
    .filter((item): item is ContinueWatchingItem => item !== null);

  const ranked = uniqueByEpisodeID(
    inputs.flatMap((input) => input.home.rails.map((episode) => toTopRankedItem(episode, input.mode)))
  );

  const sorted = sortTopRanked(ranked);
  const top10 = sorted.slice(0, 10);
  const categoryRows = categorizeRailItems(sorted.map((item) => item.episode));

  return {
    featured,
    profiles: inputs.map((input) => input.profile),
    continueWatching,
    top10,
    categoryRows
  };
}
