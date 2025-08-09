import { ReactNode } from 'react';
import { Toaster } from 'sonner';

import { useTheme } from '@/lib/hooks/useTheme';

interface ToasterProviderProps {
  children: ReactNode;
}

const ToasterProvider = ({ children }: ToasterProviderProps) => {
  const { theme } = useTheme();

  return (
    <>
      <Toaster richColors theme={theme} position='bottom-right' />
      {children}
    </>
  );
};

export default ToasterProvider;
