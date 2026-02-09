"use client";

import { NewAnalysisForm } from "@/components/home/new-analysis-form";
import { RepositoryGrid } from "@/components/home/repository-grid";
import { AnalysisProgressModal } from "@/components/home/analysis-progress-modal";
import { useRepositories } from "@/lib/hooks/use-repositories";
import { useAnalysis } from "@/lib/hooks/use-analysis";

export default function HomePage() {
  const { repositories, isLoading, refresh } = useRepositories();
  const { isAnalyzing, repoId, step, error, startAnalysis, reset } =
    useAnalysis();

  function handleAnalyze(repoUrl: string) {
    startAnalysis(repoUrl);
  }

  function handleModalClose() {
    reset();
    refresh();
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

      <AnalysisProgressModal
        open={isAnalyzing}
        step={step}
        repoId={repoId}
        error={error}
        onClose={handleModalClose}
      />
    </div>
  );
}
