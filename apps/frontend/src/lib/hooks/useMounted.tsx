import { useSyncExternalStore } from 'react';

const emptySubscribe = () => () => {};

export function useMounted() {
  return useSyncExternalStore(
    emptySubscribe,
    () => true, // Client-side: always true
    () => false, // Server-side: always false
  );
}
