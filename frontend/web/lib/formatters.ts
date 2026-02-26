export function formatDurationMS(durationMS: number): string {
  const totalMinutes = Math.max(1, Math.round(durationMS / 60000));
  return `${totalMinutes} min`;
}

export function clampPercent(input: number): number {
  if (input < 0) {
    return 0;
  }
  if (input > 100) {
    return 100;
  }

  return input;
}
