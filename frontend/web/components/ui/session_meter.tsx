import { clampPercent } from "../../lib/formatters";
import styles from "./session_meter.module.css";

type SessionMeterProps = {
  watchedMinutes: number;
  sessionLimitMinutes: number;
};

export function SessionMeter({ watchedMinutes, sessionLimitMinutes }: SessionMeterProps) {
  const usagePercent = clampPercent((watchedMinutes / Math.max(1, sessionLimitMinutes)) * 100);
  const remaining = Math.max(0, sessionLimitMinutes - watchedMinutes);

  return (
    <section className={styles.panel} aria-label="Session status">
      <div className={styles.text}>
        <p>Session health</p>
        <strong>
          {watchedMinutes} / {sessionLimitMinutes} min
        </strong>
        <span>{remaining} min left before the configured cap</span>
      </div>
      <div className={styles.track}>
        <div className={styles.fill} style={{ width: `${usagePercent}%` }} />
      </div>
    </section>
  );
}
