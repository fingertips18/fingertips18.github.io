import { useCallback, useEffect, useState } from 'react';

import { toast } from '@/lib/toast';

interface FetchProps extends Omit<RequestInit, 'signal'> {
  url: string | string[]; // allow single or multiple URLs
  toastOptions?: {
    errorTitle?: string;
    errorMessage?: string;
  };
}

export function useFetch<T extends unknown[] | object>({
  url,
  toastOptions,
  ...props
}: FetchProps) {
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const handleFetch = useCallback(
    async (signal: AbortSignal) => {
      try {
        let responses: T[] = [];

        if (Array.isArray(url)) {
          const fetchPromises: Promise<T>[] = url.map((u) =>
            fetch(u, { signal, ...props }).then(async (res) => {
              if (!res.ok) {
                const errorData = await res.json().catch(() => null);
                const message =
                  (errorData?.message as string) ||
                  `Server error: ${res.status} ${res.statusText}`;
                throw new Error(message);
              }
              return res.json() as Promise<T>;
            }),
          );

          responses = await Promise.all(fetchPromises);
        } else {
          const response = await fetch(url, { signal, ...props });
          if (!response.ok) {
            const errorData = await response.json().catch(() => null);
            const message =
              (errorData?.message as string) ||
              `Server error: ${response.status} ${response.statusText}`;
            throw new Error(message);
          }
          responses = [(await response.json()) as T];
        }

        // merge arrays only if T is an array type
        let mergedData: T;
        if (Array.isArray(responses[0])) {
          mergedData = ([] as unknown[]).concat(
            ...(responses as unknown as unknown[][]),
          ) as T;
        } else {
          mergedData = responses[0];
        }

        setData(mergedData);
        setError(null);
      } catch (error) {
        if (error instanceof Error && error.name === 'AbortError') return;

        const message = (error as Error).message || 'Something went wrong';
        setError(message);
        if (toastOptions) {
          toast({
            level: 'error',
            title: toastOptions.errorTitle || 'Unable to fetch data',
            description:
              toastOptions.errorMessage ||
              'There was a problem retrieving the requested information. Please try again.',
          });
        }
      } finally {
        setLoading(false);
      }
    },
    [url, props, toastOptions],
  );

  useEffect(() => {
    const abortController = new AbortController();
    void handleFetch(abortController.signal);
    return () => abortController.abort();
  }, [handleFetch]);

  return { data, loading, error };
}
