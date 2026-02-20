"use client";

import { AnimatePresence } from "framer-motion";
import { SubcommitCard } from "./subcommit-card";
import type { Subcommit } from "@/lib/types";

interface DayGroupProps {
  date: string;
  subcommits: Subcommit[];
  onCardClick: (subcommit: Subcommit) => void;
}

export function DayGroup({ date, subcommits, onCardClick }: DayGroupProps) {
  const formatted = new Date(date + "T00:00:00").toLocaleDateString("en-US", {
    weekday: "short",
    month: "short",
    day: "numeric",
  });

  return (
    <div className="flex shrink-0 flex-col gap-3">
      <div className="sticky top-0 z-10 bg-background/95 px-1 py-2 backdrop-blur">
        <h3 className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
          {formatted}
        </h3>
      </div>
      <div className="flex flex-col gap-3 px-1">
        <AnimatePresence>
          {subcommits.map((sc, i) => (
            <SubcommitCard
              key={sc.id}
              subcommit={sc}
              onClick={() => onCardClick(sc)}
              index={i}
            />
          ))}
        </AnimatePresence>
      </div>
    </div>
  );
}
