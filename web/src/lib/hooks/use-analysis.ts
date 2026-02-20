"use client";

import { useState } from "react";
import { api } from "@/lib/api-client";

export function useAnalysis() {
  const [isAnalyzing, setIsAnalyzing] = useState(false);
  const [repoId, setRepoId] = useState<number | null>(null);
  const [error, setError] = useState<string | null>(null);

  async function startAnalysis(repoUrl: string) {
    setIsAnalyzing(true);
    setError(null);

    try {
      const result = await api.analyzeRepository(repoUrl);
      setRepoId(result.repoId);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Analysis failed");
      setIsAnalyzing(false);
    }
  }

  function reset() {
    setIsAnalyzing(false);
    setRepoId(null);
    setError(null);
  }

  return { isAnalyzing, repoId, error, startAnalysis, reset };
}
