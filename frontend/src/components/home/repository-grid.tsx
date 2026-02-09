"use client";

import { RepositoryCard } from "./repository-card";
import { Skeleton } from "@/components/ui/skeleton";
import { GitBranch } from "lucide-react";
import type { UserRepository } from "@/lib/types";

interface RepositoryGridProps {
  repositories: UserRepository[];
  isLoading: boolean;
}

export function RepositoryGrid({ repositories, isLoading }: RepositoryGridProps) {
  if (isLoading) {
    return (
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {Array.from({ length: 6 }).map((_, i) => (
          <Skeleton key={i} className="h-24 rounded-xl" />
        ))}
      </div>
    );
  }

  if (repositories.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center gap-4 rounded-xl border border-dashed py-16">
        <div className="flex h-16 w-16 items-center justify-center rounded-full bg-muted">
          <GitBranch className="h-8 w-8 text-muted-foreground" />
        </div>
        <div className="text-center">
          <h3 className="font-semibold">No repositories analyzed yet</h3>
          <p className="mt-1 text-sm text-muted-foreground">
            Paste a GitHub URL above to analyze your first repository
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      {repositories.map((repo) => (
        <RepositoryCard key={repo.id} repo={repo} />
      ))}
    </div>
  );
}
