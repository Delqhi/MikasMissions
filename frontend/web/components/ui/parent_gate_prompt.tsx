import { BigActionButton } from "./big_action_button";
import styles from "./parent_gate_prompt.module.css";

type ParentGatePromptProps = {
  actionLabel: string;
  challengeType: "pin" | "device_confirm";
};

export function ParentGatePrompt({ actionLabel, challengeType }: ParentGatePromptProps) {
  return (
    <aside className={styles.card} aria-label="Parental gate required">
      <h3>Parent Gate</h3>
      <p>
        The action <strong>{actionLabel}</strong> requires parent approval via {challengeType.replace("_", " ")}.
      </p>
      <BigActionButton hint="Adult verification is required" label="Request verification" variant="secondary" />
    </aside>
  );
}
