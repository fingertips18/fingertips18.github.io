import { toast as sonnerToast } from 'sonner';

interface ToastProps {
  level: 'success' | 'info' | 'warning' | 'error' | 'loading';
  title: string;
  description: string;
}

/**
 * Displays a toast notification with the specified level, title, and description.
 * @param props - The toast properties
 * @param props.level - The severity level of the toast (e.g., 'success', 'error', 'warning', 'info')
 * @param props.title - The title of the toast notification
 * @param props.description - The description or message content of the toast
 * @returns The ID of the displayed toast notification
 */
export function toast({
  level,
  title,
  description,
}: ToastProps): string | number {
  return sonnerToast[level](title, {
    description,
  });
}
