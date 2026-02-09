"use client";

import { FileText } from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";
import { TypeBadge } from "./type-badge";
import type { Subcommit } from "@/lib/types";

interface SubcommitCardProps {
  subcommit: Subcommit;
  onClick: () => void;
}

export function SubcommitCard({ subcommit, onClick }: SubcommitCardProps) {
  return (
    <Card
      onClick={onClick}
      className="w-64 shrink-0 cursor-pointer transition-all hover:border-foreground/20 hover:shadow-md"
    >
      <CardContent className="space-y-2 p-4">
        <div className="flex items-center justify-between">
          <TypeBadge type={subcommit.type} />
          {subcommit.files?.length > 0 && (
            <span className="flex items-center gap-1 text-[10px] text-muted-foreground">
              <FileText className="h-3 w-3" />
              {subcommit.files.length}
            </span>
          )}
        </div>
        <h4 className="text-sm font-semibold leading-tight line-clamp-2">
          {subcommit.title}
        </h4>
        {subcommit.idea && (
          <p className="text-xs text-muted-foreground line-clamp-2">
            {subcommit.idea}
          </p>
        )}
      </CardContent>
    </Card>
  );
}
