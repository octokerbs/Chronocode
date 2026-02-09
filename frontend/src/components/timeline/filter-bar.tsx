"use client";

import { useState } from "react";
import { Search, X, Layers } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Toggle } from "@/components/ui/toggle";
import { SUBCOMMIT_TYPE_CONFIG } from "@/lib/constants";
import type { SubcommitType } from "@/lib/types";

const ALL_TYPES: SubcommitType[] = [
  "FEATURE",
  "BUG",
  "REFACTOR",
  "DOCS",
  "CHORE",
  "MILESTONE",
  "WARNING",
];

interface FilterBarProps {
  activeTypes: Set<SubcommitType>;
  onToggleType: (type: SubcommitType) => void;
  searchQuery: string;
  onSearchChange: (query: string) => void;
  epicGrouping: boolean;
  onToggleEpicGrouping: () => void;
  resultCount: number;
  totalCount: number;
}

export function FilterBar({
  activeTypes,
  onToggleType,
  searchQuery,
  onSearchChange,
  epicGrouping,
  onToggleEpicGrouping,
  resultCount,
  totalCount,
}: FilterBarProps) {
  const [showSearch, setShowSearch] = useState(false);

  return (
    <div className="sticky top-14 z-40 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="flex items-center gap-3 px-6 py-3">
        {/* Type toggles */}
        <div className="flex flex-wrap items-center gap-1.5">
          {ALL_TYPES.map((type) => {
            const config = SUBCOMMIT_TYPE_CONFIG[type];
            const isActive = activeTypes.has(type);
            return (
              <Toggle
                key={type}
                pressed={isActive}
                onPressedChange={() => onToggleType(type)}
                size="sm"
                className={`h-7 gap-1 rounded-full px-2.5 text-xs ${
                  isActive
                    ? `${config.bg} ${config.text} ${config.border} border`
                    : ""
                }`}
              >
                <span
                  className={`h-2 w-2 rounded-full ${
                    isActive ? config.text.replace("text-", "bg-") : "bg-muted-foreground/40"
                  }`}
                />
                {config.label}
              </Toggle>
            );
          })}
        </div>

        <div className="flex-1" />

        {/* Search */}
        {showSearch ? (
          <div className="flex items-center gap-2">
            <Input
              placeholder="Search subcommits..."
              value={searchQuery}
              onChange={(e) => onSearchChange(e.target.value)}
              className="h-8 w-56"
              autoFocus
            />
            <Button
              variant="ghost"
              size="icon"
              className="h-8 w-8"
              onClick={() => {
                setShowSearch(false);
                onSearchChange("");
              }}
            >
              <X className="h-4 w-4" />
            </Button>
          </div>
        ) : (
          <Button
            variant="ghost"
            size="icon"
            className="h-8 w-8"
            onClick={() => setShowSearch(true)}
          >
            <Search className="h-4 w-4" />
          </Button>
        )}

        {/* Epic grouping */}
        <Toggle
          pressed={epicGrouping}
          onPressedChange={onToggleEpicGrouping}
          size="sm"
          className="h-8 gap-1.5"
        >
          <Layers className="h-3.5 w-3.5" />
          Epics
        </Toggle>

        {/* Count */}
        <span className="text-xs text-muted-foreground">
          {resultCount === totalCount
            ? `${totalCount} subcommits`
            : `${resultCount} / ${totalCount}`}
        </span>
      </div>
    </div>
  );
}
