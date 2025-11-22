export function ensureString(value: unknown, name: string): string {
  if (typeof value === 'string') return value;
  // Coerce numbers/booleans to strings for flexibility with JSON responses
  if (typeof value === 'number' || typeof value === 'boolean') {
    return String(value);
  }
  throw new Error(`Expected property '${name}' to be a string`);
}
