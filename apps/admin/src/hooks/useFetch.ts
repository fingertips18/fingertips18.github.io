import { useCallback, useEffect, useState } from 'react';

import { toast } from '@/lib/toast';

interface FetchProps extends Omit<RequestInit, 'signal'> {
  url: string;
}

export function useFetch<T>({ url, ...props }: FetchProps) {
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const handleFetch = useCallback(
    async (signal: AbortSignal) => {
      try {
        const response = await fetch(url, {
          signal,
          ...props,
        });

        if (!response.ok) {
          const errorData = await response.json().catch(() => null);
          const message =
            (errorData?.message as string) ||
            `Server error: ${response.status} ${response.statusText}`;

          setError(message);
          throw Error(message);
        }

        const data = await response.json();

        setData(data as T);
        setError(null);
      } catch (error) {
        if (error instanceof Error && error.name === 'AbortError') {
          return;
        }
        const message = (error as Error).message || 'Something went wrong';
        setError(message);
        toast({
          level: 'error',
          title: 'Fetch error',
          description: message,
        });
      } finally {
        setLoading(false);
      }
    },
    [props, url],
  );

  useEffect(() => {
    const abortController = new AbortController();

    void handleFetch(abortController.signal);

    return () => abortController.abort();
  }, [handleFetch]);

  return { data, loading, error };
}
