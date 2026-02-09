"use client";

import { useMemo, useState } from "react";
import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area";
import { DayGroup } from "./day-group";
import { SubcommitDetailPanel } from "./subcommit-detail-panel";
import { FilterBar } from "./filter-bar";
import { Skeleton } from "@/components/ui/skeleton";
import { Loader2 } from "lucide-react";
import type { Subcommit, SubcommitType } from "@/lib/types";

interface TimelineContainerProps {
  subcommits: Subcommit[];
  isLoading: boolean;
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

export function TimelineContainer({
  subcommits,
  isLoading,
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
  const [epicGrouping, setEpicGrouping] = useState(false);
  const [selectedSubcommit, setSelectedSubcommit] = useState<Subcommit | null>(
    null,
  );

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

  const filtered = useMemo(() => {
    return subcommits.filter((sc) => {
      if (!activeTypes.has(sc.type)) return false;
      if (searchQuery) {
        const q = searchQuery.toLowerCase();
        const matchesTitle = sc.title?.toLowerCase().includes(q);
        const matchesDesc = sc.description?.toLowerCase().includes(q);
        const matchesIdea = sc.idea?.toLowerCase().includes(q);
        if (!matchesTitle && !matchesDesc && !matchesIdea) return false;
      }
      return true;
    });
  }, [subcommits, activeTypes, searchQuery]);

  const groups = useMemo(() => groupByDay(filtered), [filtered]);

  const siblings = useMemo(() => {
    if (!selectedSubcommit) return [];
    return subcommits.filter(
      (sc) =>
        sc.commitSha === selectedSubcommit.commitSha &&
        sc.id !== selectedSubcommit.id,
    );
  }, [subcommits, selectedSubcommit]);

  if (isLoading && subcommits.length === 0) {
    return (
      <>
        <FilterBar
          activeTypes={activeTypes}
          onToggleType={toggleType}
          searchQuery={searchQuery}
          onSearchChange={setSearchQuery}
          epicGrouping={epicGrouping}
          onToggleEpicGrouping={() => setEpicGrouping(!epicGrouping)}
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

  if (subcommits.length === 0) {
    return (
      <>
        <FilterBar
          activeTypes={activeTypes}
          onToggleType={toggleType}
          searchQuery={searchQuery}
          onSearchChange={setSearchQuery}
          epicGrouping={epicGrouping}
          onToggleEpicGrouping={() => setEpicGrouping(!epicGrouping)}
          resultCount={0}
          totalCount={0}
        />
        <div className="flex items-center justify-center py-32">
          <div className="flex flex-col items-center gap-4">
            <Skeleton className="h-6 w-48" />
            <p className="text-sm text-muted-foreground">
              Analysis in progress... subcommits will appear as they are
              processed.
            </p>
          </div>
        </div>
      </>
    );
  }

  return (
    <>
      <FilterBar
        activeTypes={activeTypes}
        onToggleType={toggleType}
        searchQuery={searchQuery}
        onSearchChange={setSearchQuery}
        epicGrouping={epicGrouping}
        onToggleEpicGrouping={() => setEpicGrouping(!epicGrouping)}
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
    </>
  );
}
