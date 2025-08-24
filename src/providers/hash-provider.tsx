import { useLenis } from 'lenis/react';
import { ReactNode, useEffect } from 'react';
import { useLocation } from 'react-router-dom';

interface HashProviderProps {
  children: ReactNode;
}

function HashProvider({ children }: HashProviderProps) {
  const { hash } = useLocation();
  const lenis = useLenis();

  useEffect(() => {
    if (!hash || !lenis) return;

    // Remove "#" from the hash
    const id = hash.substring(1);
    const section = document.getElementById(id);

    if (!section) return;

    // When changing a route, the DOM tree changes height,
    // But its not aware of the change, so we need to resize it before scrolling
    lenis.resize();

    lenis.scrollTo(section);
  }, [hash, lenis]);

  return <>{children}</>;
}

export default HashProvider;
