"use client";

import { useState } from "react";
import { api } from "@/lib/api-client";

export type AnalysisStep = "idle" | "preparing" | "analyzing" | "ready";

export function useAnalysis() {
  const [isAnalyzing, setIsAnalyzing] = useState(false);
  const [repoId, setRepoId] = useState<number | null>(null);
  const [step, setStep] = useState<AnalysisStep>("idle");
  const [error, setError] = useState<string | null>(null);

  async function startAnalysis(repoUrl: string) {
    setIsAnalyzing(true);
    setStep("preparing");
    setError(null);

    try {
      const result = await api.analyzeRepository(repoUrl);
      setRepoId(result.repoId);
      setStep("analyzing");

      // After a delay, allow navigation
      setTimeout(() => setStep("ready"), 4000);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Analysis failed");
      setStep("idle");
      setIsAnalyzing(false);
    }
  }

  function reset() {
    setIsAnalyzing(false);
    setRepoId(null);
    setStep("idle");
    setError(null);
  }

  return { isAnalyzing, repoId, step, error, startAnalysis, reset };
}
