"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { NewAnalysisForm } from "@/components/home/new-analysis-form";
import { RepositoryGrid } from "@/components/home/repository-grid";
import { AnalysisErrorDialog } from "@/components/home/analysis-error-dialog";
import { useRepositories } from "@/lib/hooks/use-repositories";
import { useAnalysis } from "@/lib/hooks/use-analysis";

export default function HomePage() {
  const router = useRouter();
  const { repositories, isLoading, refresh } = useRepositories();
  const { isAnalyzing, repoId, error, startAnalysis, reset } = useAnalysis();

  useEffect(() => {
    if (repoId) {
      refresh();
      router.push(`/timeline/${repoId}`);
    }
  }, [repoId, router, refresh]);

  function handleAnalyze(repoUrl: string) {
    startAnalysis(repoUrl);
  }

  function handleErrorClose() {
    reset();
  }

  return (
    <div className="mx-auto max-w-5xl px-6 py-8">
      <div className="space-y-8">
        <div className="space-y-4">
          <h1 className="text-2xl font-bold tracking-tight">
            Your Repositories
          </h1>
          <NewAnalysisForm
            onAnalyze={handleAnalyze}
            isAnalyzing={isAnalyzing}
          />
        </div>

        <RepositoryGrid repositories={repositories} isLoading={isLoading} />
      </div>

      <AnalysisErrorDialog
        open={!!error}
        error={error}
        onClose={handleErrorClose}
      />
    </div>
  );
}
