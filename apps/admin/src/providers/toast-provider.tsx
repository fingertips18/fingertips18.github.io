import type { ReactNode } from 'react';
import { Toaster } from 'sonner';

interface ToastProviderProps {
  children: ReactNode;
}

export default function ToastProvider({ children }: ToastProviderProps) {
  return (
    <>
      <Toaster richColors position='top-right' />
      {children}
    </>
  );
}
