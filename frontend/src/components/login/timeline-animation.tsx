"use client";

import { motion } from "framer-motion";
import { SUBCOMMIT_DOT_COLORS } from "@/lib/constants";
import type { SubcommitType } from "@/lib/types";

const mockCards: { type: SubcommitType; title: string }[] = [
  { type: "FEATURE", title: "Add OAuth2 authentication" },
  { type: "REFACTOR", title: "Extract service layer" },
  { type: "BUG", title: "Fix race condition in worker" },
  { type: "DOCS", title: "Update API documentation" },
  { type: "FEATURE", title: "Implement timeline view" },
  { type: "CHORE", title: "Upgrade dependencies" },
  { type: "MILESTONE", title: "v1.0 release preparation" },
  { type: "WARNING", title: "Deprecate legacy endpoint" },
];

export function TimelineAnimation() {
  return (
    <div className="relative h-80 overflow-hidden">
      {/* Timeline line */}
      <div className="absolute left-8 top-0 bottom-0 w-px bg-border" />

      {/* Animated cards */}
      <div className="relative pl-16 space-y-3">
        {mockCards.map((card, i) => (
          <motion.div
            key={i}
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{
              delay: i * 0.4,
              duration: 0.5,
              repeat: Infinity,
              repeatDelay: mockCards.length * 0.4 + 2,
            }}
            className="flex items-center gap-3"
          >
            {/* Dot */}
            <div className="absolute left-8 -translate-x-1/2">
              <div
                className={`h-2.5 w-2.5 rounded-full ${SUBCOMMIT_DOT_COLORS[card.type]}`}
              />
            </div>

            {/* Card */}
            <div className="rounded-lg border bg-card/50 px-4 py-2.5 backdrop-blur-sm">
              <span className="text-xs font-medium text-muted-foreground">
                {card.type}
              </span>
              <p className="text-sm font-medium">{card.title}</p>
            </div>
          </motion.div>
        ))}
      </div>
    </div>
  );
}
