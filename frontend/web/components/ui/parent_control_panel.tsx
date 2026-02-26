"use client";

import { useState } from "react";
import type { ParentalControls } from "../../lib/experience_types";
import styles from "./parent_control_panel.module.css";

type ParentControlPanelProps = {
  childProfileID: string;
  initialControls: ParentalControls;
};

type ToggleField = "autoplay" | "chat_enabled" | "external_links";

type ToggleMeta = {
  label: string;
  detail: string;
};

const toggleMeta: Record<ToggleField, ToggleMeta> = {
  autoplay: {
    label: "Autoplay",
    detail: "Automatically continue to the next approved mission"
  },
  chat_enabled: {
    label: "Chat",
    detail: "Enable protected in-app conversation surfaces"
  },
  external_links: {
    label: "External links",
    detail: "Allow leaving MikasMissions to third-party websites"
  }
};

export function ParentControlPanel(props: ParentControlPanelProps) {
  const { childProfileID, initialControls } = props;
  const [controls, setControls] = useState(initialControls);
  const [message, setMessage] = useState("Strict mode active by default.");
  const [saving, setSaving] = useState(false);

  function toggle(field: ToggleField) {
    setControls((previous) => ({
      ...previous,
      [field]: !previous[field]
    }));
  }

  async function saveControls() {
    setSaving(true);
    setMessage("Saving parental controls...");

    try {
      const response = await fetch(`/api/parents/controls/${childProfileID}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify(controls)
      });

      if (!response.ok) {
        setMessage("Save failed. Please retry.");
        return;
      }

      setMessage("Controls updated and audit logged.");
    } catch {
      setMessage("Save failed. Please retry.");
    } finally {
      setSaving(false);
    }
  }

  return (
    <section className={styles.panel} aria-label="Parent control panel" id="controls">
      <header className={styles.header}>
        <div>
          <h3>Parent controls</h3>
          <p>Configure child-safe defaults without opening every advanced setting.</p>
        </div>
        <div className={styles.badges}>
          <span>Safety: {controls.safety_mode}</span>
          <span>Session: {controls.session_limit_minutes} min</span>
        </div>
      </header>

      <div className={styles.grid}>
        <label>
          Safety mode
          <select
            onChange={(event) =>
              setControls((prev) => ({
                ...prev,
                safety_mode: event.target.value === "balanced" ? "balanced" : "strict"
              }))
            }
            value={controls.safety_mode}
          >
            <option value="strict">Strict</option>
            <option value="balanced">Balanced</option>
          </select>
        </label>

        <label>
          Session limit (minutes)
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
          Bedtime window
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
              <strong>{controls[field] ? "On" : "Off"}</strong>
            </button>
          </li>
        ))}
      </ul>

      <footer className={styles.footer}>
        <p role="status">{message}</p>
        <button disabled={saving} onClick={saveControls} type="button">
          {saving ? "Saving..." : "Save controls"}
        </button>
      </footer>
    </section>
  );
}
