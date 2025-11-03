import { useSyncExternalStore } from 'react';

const useResize = () => {
  const subscribe = (callback: () => void) => {
    window.addEventListener('resize', callback);
    return () => window.removeEventListener('resize', callback);
  };

  const width = useSyncExternalStore(
    subscribe,
    () => window.innerWidth,
    () => 0, // Server snapshot
  );

  const height = useSyncExternalStore(
    subscribe,
    () => window.innerHeight,
    () => 0, // Server snapshot
  );

  return { width, height };
};

export { useResize };
