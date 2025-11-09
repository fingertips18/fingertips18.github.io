import { toast as sonnerToast } from 'sonner';

interface ToastProps {
  level: 'success' | 'info' | 'warning' | 'error' | 'loading';
  title: string;
  description: string;
}

export function toast({ level, title, description }: ToastProps) {
  return sonnerToast[level](title, {
    description,
  });
}
