import useSWR from "swr";
import { api } from "@/lib/api-client";

export function useSubcommits(repoId: string | null) {
  const { data, error, isLoading, mutate } = useSWR(
    repoId ? `/subcommits-timeline/${repoId}` : null,
    () => api.getSubcommitsTimeline(repoId!),
    { refreshInterval: 5000 },
  );

  return {
    subcommits: data?.subcommits ?? [],
    isAnalyzing: data?.isAnalyzing ?? false,
    repoUrl: data?.repoUrl ?? "",
    isLoading,
    error,
    refresh: mutate,
  };
}
