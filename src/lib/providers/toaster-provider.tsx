import { Toaster } from "sonner";

import { useTheme } from "@/lib/hooks/useTheme";

interface ToasterProviderProps {
  children: React.ReactNode;
}

const ToasterProvider = ({ children }: ToasterProviderProps) => {
  const { theme } = useTheme();

  return (
    <>
      <Toaster
        richColors
        theme={theme}
        position="bottom-right"
        pauseWhenPageIsHidden
      />
      {children}
    </>
  );
};

export default ToasterProvider;
