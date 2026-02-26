"use client";

import { useState } from "react";
import type { ParentalControls } from "../../lib/experience_types";
import type { LocalizedMessages } from "../../lib/messages/types";
import styles from "./parent_control_panel.module.css";

type ParentControlPanelProps = {
  childProfileID: string;
  initialControls: ParentalControls;
  messages: LocalizedMessages["parents"]["controlPanel"];
};

type ToggleField = "autoplay" | "chat_enabled" | "external_links";

type ToggleMeta = {
  label: string;
  detail: string;
};

export function ParentControlPanel(props: ParentControlPanelProps) {
  const { childProfileID, initialControls, messages } = props;
  const toggleMeta: Record<ToggleField, ToggleMeta> = {
    autoplay: {
      label: messages.autoplayLabel,
      detail: messages.autoplayDetail
    },
    chat_enabled: {
      label: messages.chatLabel,
      detail: messages.chatDetail
    },
    external_links: {
      label: messages.linksLabel,
      detail: messages.linksDetail
    }
  };
  const [controls, setControls] = useState(initialControls);
  const [message, setMessage] = useState(messages.strictModeDefault);
  const [saving, setSaving] = useState(false);

  function toggle(field: ToggleField) {
    setControls((previous) => ({
      ...previous,
      [field]: !previous[field]
    }));
  }

  async function saveControls() {
    setSaving(true);
    setMessage(messages.savingMessage);

    try {
      const response = await fetch(`/api/parents/controls/${childProfileID}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify(controls)
      });

      if (!response.ok) {
        setMessage(messages.saveFailed);
        return;
      }

      setMessage(messages.saveSuccess);
    } catch {
      setMessage(messages.saveFailed);
    } finally {
      setSaving(false);
    }
  }

  return (
    <section className={styles.panel} aria-label="Parent control panel" id="controls">
      <header className={styles.header}>
        <div>
          <h3>{messages.title}</h3>
          <p>{messages.description}</p>
        </div>
        <div className={styles.badges}>
          <span>
            {messages.badgeSafety}: {controls.safety_mode}
          </span>
          <span>
            {messages.badgeSession}: {controls.session_limit_minutes} min
          </span>
        </div>
      </header>

      <div className={styles.grid}>
        <label>
          {messages.safetyMode}
          <select
            onChange={(event) =>
              setControls((prev) => ({
                ...prev,
                safety_mode: event.target.value === "balanced" ? "balanced" : "strict"
              }))
            }
            value={controls.safety_mode}
          >
            <option value="strict">{messages.strict}</option>
            <option value="balanced">{messages.balanced}</option>
          </select>
        </label>

        <label>
          {messages.sessionLimit}
          <input
            max={120}
            min={15}
            onChange={(event) =>
              setControls((prev) => ({
                ...prev,
                session_limit_minutes: Number(event.target.value)
              }))
            }
            type="range"
            value={controls.session_limit_minutes}
          />
          <strong>{controls.session_limit_minutes} min</strong>
        </label>

        <label>
          {messages.bedtimeWindow}
          <input
            onChange={(event) =>
              setControls((prev) => ({
                ...prev,
                bedtime_window: event.target.value
              }))
            }
            placeholder="20:00-06:30"
            type="text"
            value={controls.bedtime_window}
          />
        </label>
      </div>

      <ul className={styles.toggles}>
        {(["autoplay", "chat_enabled", "external_links"] as ToggleField[]).map((field) => (
          <li key={field}>
            <button
              aria-pressed={controls[field]}
              className={controls[field] ? styles.enabled : styles.disabled}
              onClick={() => toggle(field)}
              type="button"
            >
              <div>
                <span>{toggleMeta[field].label}</span>
                <small>{toggleMeta[field].detail}</small>
              </div>
              <strong>{controls[field] ? messages.on : messages.off}</strong>
            </button>
          </li>
        ))}
      </ul>

      <footer className={styles.footer}>
        <p role="status">{message}</p>
        <button disabled={saving} onClick={saveControls} type="button">
          {saving ? messages.saving : messages.saveControls}
        </button>
      </footer>
    </section>
  );
}
