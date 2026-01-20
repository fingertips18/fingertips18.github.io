/**
 * Ensures a value is a string, coercing numbers and booleans if necessary.
 * @param {Object} options - The options object
 * @param {unknown} options.value - The value to validate and coerce
 * @param {string} options.name - The property name for error messages
 * @returns {string} The validated or coerced string value
 * @throws {Error} If the value cannot be coerced to a string
 */
export function ensureString({
  value,
  name,
}: {
  value: unknown;
  name: string;
}): string {
  if (typeof value === 'string') return value;
  // Coerce numbers/booleans to strings for flexibility with JSON responses
  if (typeof value === 'number' || typeof value === 'boolean') {
    return String(value);
  }
  throw new Error(`Expected property '${name}' to be a string`);
}

/**
 * Ensures a value is a valid Date, coercing from string representation if necessary.
 * @param {Object} options - The options object
 * @param {unknown} options.value - The value to validate and coerce
 * @param {string} options.name - The property name for error messages
 * @returns {Date} The validated or coerced Date value
 * @throws {Error} If the value cannot be coerced to a valid Date
 */
export function ensureDate({
  value,
  name,
}: {
  value: unknown;
  name: string;
}): Date {
  const date = new Date(ensureString({ value, name }));
  if (isNaN(date.getTime())) {
    throw new Error(`Invalid date for property '${name}'`);
  }

  return date;
}

/**
 * Ensures a value is a valid number, coercing from string representation if necessary.
 * @param {Object} options - The options object
 * @param {unknown} options.value - The value to validate and coerce
 * @param {string} options.name - The property name for error messages
 * @returns {number} The validated or coerced number value
 * @throws {Error} If the value cannot be coerced to a valid number
 */
export function ensureNumber({
  value,
  name,
}: {
  value: unknown;
  name: string;
}): number {
  const num = Number(ensureString({ value, name }));
  if (isNaN(num)) {
    throw new Error(`Invalid number for property '${name}'`);
  }

  return num;
}
