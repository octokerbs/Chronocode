import useSWR from "swr";
import { api } from "@/lib/api-client";

export function useRepositories() {
  const { data, error, isLoading, mutate } = useSWR(
    "/repositories",
    () => api.getRepositories(),
  );

  return {
    repositories: data?.repositories ?? [],
    isLoading,
    error,
    refresh: mutate,
  };
}
