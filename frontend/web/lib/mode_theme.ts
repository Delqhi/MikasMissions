import type { KidsMode } from "./experience_types";

export const modeClassName: Record<KidsMode, string> = {
  early: "mode-early",
  core: "mode-core",
  teen: "mode-teen"
};

export const modeLabel: Record<KidsMode, string> = {
  early: "Kids Early (3-5)",
  core: "Kids Core (6-11)",
  teen: "Teens (12-16)"
};

export function toMode(input: string | null): KidsMode {
  if (input === "early" || input === "core" || input === "teen") {
    return input;
  }

  return "core";
}
