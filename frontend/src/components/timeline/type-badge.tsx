import { SUBCOMMIT_TYPE_CONFIG } from "@/lib/constants";
import type { SubcommitType } from "@/lib/types";

interface TypeBadgeProps {
  type: SubcommitType;
  size?: "sm" | "md";
}

export function TypeBadge({ type, size = "sm" }: TypeBadgeProps) {
  const config = SUBCOMMIT_TYPE_CONFIG[type] ?? SUBCOMMIT_TYPE_CONFIG.CHORE;

  return (
    <span
      className={`inline-flex items-center rounded-full border font-medium ${config.bg} ${config.text} ${config.border} ${
        size === "sm" ? "px-2 py-0.5 text-[10px]" : "px-2.5 py-0.5 text-xs"
      }`}
    >
      {config.label}
    </span>
  );
}
