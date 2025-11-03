import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

/**
 * Combines multiple class name values into a single string, filtering out falsy values,
 * and merges Tailwind CSS classes intelligently to avoid conflicts.
 *
 * @param inputs - A list of class values (strings, arrays, or objects) to be combined.
 * @returns A single string of merged class names.
 *
 * @example
 * ```typescript
 * cn('btn', { 'btn-active': isActive }, ['extra-class'])
 * // => 'btn btn-active extra-class'
 * ```
 */
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}
