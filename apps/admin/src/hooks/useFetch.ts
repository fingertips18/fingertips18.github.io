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

        const data = await response.json();

        if (!response.ok) {
          const message = 'Server error occured';
          setError(message);
          throw Error(message);
        }

        setData(data as T);
        setError(null);
      } catch (error) {
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
