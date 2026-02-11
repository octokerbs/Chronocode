import useSWR from "swr";
import { useEffect, useState } from "react";
import { api } from "@/lib/api-client";

function useDebounce(value: string, delay: number) {
  const [debounced, setDebounced] = useState(value);
  useEffect(() => {
    const timer = setTimeout(() => setDebounced(value), delay);
    return () => clearTimeout(timer);
  }, [value, delay]);
  return debounced;
}

export function useRepoSearch(query: string) {
  const debouncedQuery = useDebounce(query, 300);

  const { data, error, isLoading } = useSWR(
    debouncedQuery.length >= 2
      ? `/user/repos/search?q=${debouncedQuery}`
      : null,
    () => api.searchRepos(debouncedQuery),
  );

  return {
    results: data?.repositories ?? [],
    isLoading,
    error,
  };
}
