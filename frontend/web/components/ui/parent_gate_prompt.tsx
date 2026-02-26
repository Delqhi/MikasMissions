import { BigActionButton } from "./big_action_button";
import type { LocalizedMessages } from "../../lib/messages/types";
import styles from "./parent_gate_prompt.module.css";

type ParentGatePromptProps = {
  actionLabel: string;
  challengeType: "pin" | "device_confirm";
  labels?: LocalizedMessages["kids"]["gate"];
};

const defaultLabels = {
  title: "Parent Gate",
  requestVerification: "Request verification",
  verificationHint: "Adult verification is required",
  requiresApproval: "The action requires parent approval via"
} as const;

export function ParentGatePrompt({ actionLabel, challengeType, labels }: ParentGatePromptProps) {
  const text = labels ?? defaultLabels;

  return (
    <aside className={styles.card} aria-label="Parental gate required">
      <h3>{text.title}</h3>
      <p>
        {text.requiresApproval} <strong>{actionLabel}</strong> via {challengeType.replace("_", " ")}.
      </p>
      <BigActionButton hint={text.verificationHint} label={text.requestVerification} variant="secondary" />
    </aside>
  );
}
