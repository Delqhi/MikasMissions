import { clampPercent } from "../../lib/formatters";
import type { LocalizedMessages } from "../../lib/messages/types";
import styles from "./session_meter.module.css";

type SessionMeterProps = {
  watchedMinutes: number;
  sessionLimitMinutes: number;
  labels?: Pick<LocalizedMessages["kids"]["cards"], "sessionHealth" | "sessionCap" | "minLeftSuffix">;
};

const defaultLabels = {
  sessionHealth: "Session health",
  sessionCap: "Session cap",
  minLeftSuffix: "min left"
} as const;

export function SessionMeter({ watchedMinutes, sessionLimitMinutes, labels }: SessionMeterProps) {
  const usagePercent = clampPercent((watchedMinutes / Math.max(1, sessionLimitMinutes)) * 100);
  const remaining = Math.max(0, sessionLimitMinutes - watchedMinutes);
  const text = labels ?? defaultLabels;

  return (
    <section className={styles.panel} aria-label="Session status">
      <div className={styles.text}>
        <p>{text.sessionHealth}</p>
        <strong>
          {watchedMinutes} / {sessionLimitMinutes} min
        </strong>
        <span>
          {text.sessionCap}: {remaining} {text.minLeftSuffix}
        </span>
      </div>
      <div className={styles.track}>
        <div className={styles.fill} style={{ width: `${usagePercent}%` }} />
      </div>
    </section>
  );
}
