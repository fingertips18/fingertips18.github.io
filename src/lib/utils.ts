import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

/**
 * Merges multiple class name values into a single string, handling conditional and array values,
 * and resolving Tailwind CSS class conflicts.
 *
 * @param inputs - A list of class values (strings, arrays, or objects) to be merged.
 * @returns A single string with merged and deduplicated class names.
 *
 * @example
 * ```ts
 * cn('btn', { 'btn-active': isActive }, ['extra-class']);
 * // => "btn btn-active extra-class" (if isActive is true)
 * ```
 */
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

/**
 * Returns a new array with the elements of the input array shuffled in random order.
 * The original array is not mutated.
 *
 * @typeParam T - The type of elements in the array.
 * @param array - The array to shuffle.
 * @returns A new array with the elements shuffled.
 */
export function shuffleArray<T>(array: T[]): T[] {
  const arr = array.slice(); // clone the array to not mutate the original
  for (let i = arr.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i - 1)); // random index from 0 to i
    [arr[i], arr[j]] = [arr[j], arr[i]]; // swap
  }

  return arr;
}
