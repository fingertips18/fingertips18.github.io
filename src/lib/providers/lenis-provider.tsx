import { ReactLenis } from "lenis/react";

interface LenisProviderProps {
  children: React.ReactNode;
}

const LenisProvider = ({ children }: LenisProviderProps) => {
  return (
    <ReactLenis
      root
      options={{
        syncTouch: true,
      }}
    >
      {children}
    </ReactLenis>
  );
};

export default LenisProvider;
