"use client";

import { useMemo, useState } from "react";
import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area";
import { DayGroup } from "./day-group";
import { SubcommitDetailPanel } from "./subcommit-detail-panel";
import { FilterBar } from "./filter-bar";
import { AnalysisStatus } from "./analysis-status";
import { Loader2 } from "lucide-react";
import type { Subcommit, SubcommitType } from "@/lib/types";

interface TimelineContainerProps {
  subcommits: Subcommit[];
  isLoading: boolean;
  isAnalyzing?: boolean;
  repoUrl?: string;
}

function groupByDay(subcommits: Subcommit[]): Map<string, Subcommit[]> {
  const groups = new Map<string, Subcommit[]>();
  for (const sc of subcommits) {
    const day = sc.createdAt
      ? new Date(sc.createdAt).toISOString().split("T")[0]
      : "unknown";
    const existing = groups.get(day) ?? [];
    groups.set(day, [...existing, sc]);
  }
  return new Map([...groups.entries()].sort().reverse());
}

function repoNameFromUrl(url?: string): string {
  if (!url) return "";
  try {
    const parts = new URL(url).pathname.split("/").filter(Boolean);
    if (parts.length >= 2) return `${parts[0]}/${parts[1]}`;
  } catch {
    // fallback: try splitting raw string
    const segments = url.split("/").filter(Boolean);
    if (segments.length >= 2) return `${segments[segments.length - 2]}/${segments[segments.length - 1]}`;
  }
  return url;
}

export function TimelineContainer({
  subcommits,
  isLoading,
  isAnalyzing = false,
  repoUrl,
}: TimelineContainerProps) {
  const [activeTypes, setActiveTypes] = useState<Set<SubcommitType>>(
    new Set([
      "FEATURE",
      "BUG",
      "REFACTOR",
      "DOCS",
      "CHORE",
      "MILESTONE",
      "WARNING",
    ]),
  );
  const [searchQuery, setSearchQuery] = useState("");
  const [activeEpics, setActiveEpics] = useState<Set<string> | null>(null);
  const [selectedSubcommit, setSelectedSubcommit] = useState<Subcommit | null>(
    null,
  );

  const repoName = useMemo(() => repoNameFromUrl(repoUrl), [repoUrl]);

  const allEpics = useMemo(() => {
    const epics = new Set<string>();
    for (const sc of subcommits) {
      if (sc.epic) epics.add(sc.epic);
    }
    return [...epics].sort();
  }, [subcommits]);

  // Initialize activeEpics to all epics on first data load
  const resolvedActiveEpics = useMemo(() => {
    if (activeEpics !== null) return activeEpics;
    return new Set(allEpics);
  }, [activeEpics, allEpics]);

  function toggleType(type: SubcommitType) {
    setActiveTypes((prev) => {
      const next = new Set(prev);
      if (next.has(type)) {
        next.delete(type);
      } else {
        next.add(type);
      }
      return next;
    });
  }

  function toggleEpic(epic: string) {
    setActiveEpics((prev) => {
      const current = prev ?? new Set(allEpics);
      const next = new Set(current);
      if (next.has(epic)) {
        next.delete(epic);
      } else {
        next.add(epic);
      }
      return next;
    });
  }

  const filtered = useMemo(() => {
    return subcommits.filter((sc) => {
      if (!activeTypes.has(sc.type)) return false;
      if (sc.epic && !resolvedActiveEpics.has(sc.epic)) return false;
      if (searchQuery) {
        const q = searchQuery.toLowerCase();
        const matchesTitle = sc.title?.toLowerCase().includes(q);
        const matchesDesc = sc.description?.toLowerCase().includes(q);
        const matchesIdea = sc.idea?.toLowerCase().includes(q);
        if (!matchesTitle && !matchesDesc && !matchesIdea) return false;
      }
      return true;
    });
  }, [subcommits, activeTypes, resolvedActiveEpics, searchQuery]);

  const groups = useMemo(() => groupByDay(filtered), [filtered]);

  const siblings = useMemo(() => {
    if (!selectedSubcommit) return [];
    return subcommits.filter(
      (sc) =>
        sc.commitSha === selectedSubcommit.commitSha &&
        sc.id !== selectedSubcommit.id,
    );
  }, [subcommits, selectedSubcommit]);

  const filterBarProps = {
    repoName,
    activeTypes,
    onToggleType: toggleType,
    epics: allEpics,
    activeEpics: resolvedActiveEpics,
    onToggleEpic: toggleEpic,
    searchQuery,
    onSearchChange: setSearchQuery,
  };

  if (isLoading && subcommits.length === 0) {
    return (
      <>
        <FilterBar
          {...filterBarProps}
          resultCount={0}
          totalCount={0}
        />
        <div className="flex items-center justify-center py-32">
          <div className="flex flex-col items-center gap-4">
            <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
            <p className="text-sm text-muted-foreground">
              Loading timeline...
            </p>
          </div>
        </div>
      </>
    );
  }

  if (subcommits.length === 0 && !isLoading) {
    return (
      <>
        <FilterBar
          {...filterBarProps}
          resultCount={0}
          totalCount={0}
        />
        <div className="flex items-center justify-center py-32">
          <div className="flex flex-col items-center gap-4">
            {isAnalyzing ? (
              <>
                <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
                <p className="text-sm text-muted-foreground">
                  Analyzing commits... subcommits will appear as they are
                  processed.
                </p>
              </>
            ) : (
              <p className="text-sm text-muted-foreground">
                No subcommits found for this repository.
              </p>
            )}
          </div>
        </div>
        <AnalysisStatus isAnalyzing={isAnalyzing} />
      </>
    );
  }

  return (
    <>
      <FilterBar
        {...filterBarProps}
        resultCount={filtered.length}
        totalCount={subcommits.length}
      />

      <ScrollArea className="w-full">
        <div className="flex gap-6 p-6">
          {[...groups.entries()].map(([date, scs]) => (
            <DayGroup
              key={date}
              date={date}
              subcommits={scs}
              onCardClick={setSelectedSubcommit}
            />
          ))}
        </div>
        <ScrollBar orientation="horizontal" />
      </ScrollArea>

      <SubcommitDetailPanel
        subcommit={selectedSubcommit}
        siblings={siblings}
        repoUrl={repoUrl}
        onClose={() => setSelectedSubcommit(null)}
        onSiblingClick={setSelectedSubcommit}
      />

      <AnalysisStatus isAnalyzing={isAnalyzing} />
    </>
  );
}
