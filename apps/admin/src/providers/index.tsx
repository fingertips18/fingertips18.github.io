import type { ReactNode } from 'react';

import { SidebarProvider } from '@/components/shadcn/sidebar';

import ToastProvider from './toast-provider';

interface ProvidersProps {
  children: ReactNode;
}

export default function Providers({ children }: ProvidersProps) {
  return (
    <SidebarProvider>
      <ToastProvider>{children}</ToastProvider>
    </SidebarProvider>
  );
}
