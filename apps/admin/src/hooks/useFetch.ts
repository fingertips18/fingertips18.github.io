import { useCallback, useEffect, useRef, useState } from 'react';

import { toast } from '@/lib/toast';

interface FetchProps extends Omit<RequestInit, 'signal'> {
  url: string | string[];
  toastOptions?: {
    errorTitle?: string;
    errorMessage?: string;
  };
}

/**
 * Custom hook for fetching data from one or multiple URLs with error handling and loading states.
 *
 * @template T - The expected shape of the response data. Can be an array or object.
 *
 * @param {FetchProps} props - Configuration object for the fetch request.
 * @param {string | string[]} props.url - Single URL or array of URLs to fetch from.
 * @param {ToastOptions} [props.toastOptions] - Optional configuration for error toast notifications.
 * @param {...any} props - Additional fetch options (headers, method, etc.) passed to the fetch API.
 *
 * @returns {Object} Fetch state object.
 * @returns {T | null} data - The fetched data, or null if not yet loaded or an error occurred.
 * @returns {boolean} loading - Loading state indicator.
 * @returns {string | null} error - Error message if the fetch failed, or null if successful.
 *
 * @example
 * const { data, loading, error } = useFetch<User>({
 *   url: '/api/user',
 *   toastOptions: { errorTitle: 'Failed to load user' }
 * });
 *
 * @example
 * const { data, loading, error } = useFetch<[Users[], Posts[]]>({
 *   url: ['/api/users', '/api/posts']
 * });
 *
 * @remarks
 * - When multiple URLs are provided, responses are returned as a tuple.
 * - Automatically aborts previous requests when component unmounts or dependencies change.
 * - Refs are used to prevent unnecessary re-fetches due to object identity changes.
 * - AbortErrors are silently ignored as they represent intentional request cancellations.
 */
export function useFetch<T extends unknown[] | object>({
  url,
  toastOptions,
  ...props
}: FetchProps): { data: T | null; loading: boolean; error: string | null } {
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Store url, latest props and toastOptions in refs to avoid re-fetches on identity changes
  const urlRef = useRef(url);
  const propsRef = useRef(props);
  const toastOptionsRef = useRef(toastOptions);

  useEffect(() => {
    urlRef.current = url;
    propsRef.current = props;
    toastOptionsRef.current = toastOptions;
  }, [url, props, toastOptions]);

  const handleFetch = useCallback(
    async (signal: AbortSignal) => {
      const url = urlRef.current;

      try {
        setLoading(true); // Reset loading on each fetch
        let responses: T[] = [];

        if (Array.isArray(url)) {
          const fetchPromises: Promise<T>[] = url.map((u) =>
            fetch(u, { signal, ...propsRef.current }).then(async (res) => {
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
          const response = await fetch(url, { signal, ...propsRef.current });
          if (!response.ok) {
            const errorData = await response.json().catch(() => null);
            const message =
              (errorData?.message as string) ||
              `Server error: ${response.status} ${response.statusText}`;
            throw new Error(message);
          }
          responses = [(await response.json()) as T];
        }

        let mergedData: T;
        if (Array.isArray(url)) {
          // Return all responses as a tuple when multiple URLs provided
          mergedData = responses as T;
        } else {
          mergedData = responses[0];
        }

        setData(mergedData);
        setError(null);
      } catch (error) {
        if (error instanceof Error && error.name === 'AbortError') return;

        const message = (error as Error).message || 'Something went wrong';
        setError(message);
        const currentToastOptions = toastOptionsRef.current;
        if (currentToastOptions) {
          toast({
            level: 'error',
            title: currentToastOptions.errorTitle || 'Unable to fetch data',
            description:
              currentToastOptions.errorMessage ||
              'There was a problem retrieving the requested information. Please try again.',
          });
        }
      } finally {
        setLoading(false);
      }
    },
    [], // Only depend on url - props changes won't trigger re-fetches
  );

  useEffect(() => {
    const abortController = new AbortController();
    void handleFetch(abortController.signal);
    return () => abortController.abort();
  }, [handleFetch]);

  return { data, loading, error };
}
