import type { SubcommitType } from "./types";

export const SUBCOMMIT_TYPE_CONFIG: Record<
  SubcommitType,
  { bg: string; text: string; border: string; label: string }
> = {
  FEATURE: {
    bg: "bg-green-500/15",
    text: "text-green-600 dark:text-green-400",
    border: "border-green-500/30",
    label: "Feature",
  },
  BUG: {
    bg: "bg-red-500/15",
    text: "text-red-600 dark:text-red-400",
    border: "border-red-500/30",
    label: "Bug",
  },
  REFACTOR: {
    bg: "bg-purple-500/15",
    text: "text-purple-600 dark:text-purple-400",
    border: "border-purple-500/30",
    label: "Refactor",
  },
  DOCS: {
    bg: "bg-blue-500/15",
    text: "text-blue-600 dark:text-blue-400",
    border: "border-blue-500/30",
    label: "Docs",
  },
  WARNING: {
    bg: "bg-orange-500/15",
    text: "text-orange-600 dark:text-orange-400",
    border: "border-orange-500/30",
    label: "Warning",
  },
  CHORE: {
    bg: "bg-gray-500/15",
    text: "text-gray-600 dark:text-gray-400",
    border: "border-gray-500/30",
    label: "Chore",
  },
  MILESTONE: {
    bg: "bg-yellow-500/15",
    text: "text-yellow-600 dark:text-yellow-400",
    border: "border-yellow-500/30",
    label: "Milestone",
  },
};

export const SUBCOMMIT_DOT_COLORS: Record<SubcommitType, string> = {
  FEATURE: "bg-green-500",
  BUG: "bg-red-500",
  REFACTOR: "bg-purple-500",
  DOCS: "bg-blue-500",
  WARNING: "bg-orange-500",
  CHORE: "bg-gray-400",
  MILESTONE: "bg-yellow-500",
};
