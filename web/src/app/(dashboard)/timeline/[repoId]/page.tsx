"use client";

import { use } from "react";
import { useSubcommits } from "@/lib/hooks/use-subcommits";
import { TimelineContainer } from "@/components/timeline/timeline-container";

export default function TimelinePage({
  params,
}: {
  params: Promise<{ repoId: string }>;
}) {
  const { repoId } = use(params);
  const { subcommits, isLoading } = useSubcommits(repoId);

  return (
    <div className="flex h-[calc(100vh-3.5rem)] flex-col">
      <TimelineContainer
        subcommits={subcommits}
        isLoading={isLoading}
      />
    </div>
  );
}
