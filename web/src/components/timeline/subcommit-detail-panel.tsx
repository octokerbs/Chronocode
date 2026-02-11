"use client";

import { ExternalLink, FileText, GitCommit } from "lucide-react";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { TypeBadge } from "./type-badge";
import type { Subcommit } from "@/lib/types";

interface SubcommitDetailPanelProps {
  subcommit: Subcommit | null;
  siblings: Subcommit[];
  repoUrl?: string;
  onClose: () => void;
  onSiblingClick: (sc: Subcommit) => void;
}

export function SubcommitDetailPanel({
  subcommit,
  siblings,
  repoUrl,
  onClose,
  onSiblingClick,
}: SubcommitDetailPanelProps) {
  if (!subcommit) return null;

  const commitUrl = repoUrl
    ? `${repoUrl}/commit/${subcommit.commitSha}`
    : null;

  return (
    <Sheet open={!!subcommit} onOpenChange={() => onClose()}>
      <SheetContent className="w-full overflow-y-auto sm:max-w-lg">
        <SheetHeader>
          <div className="flex items-center gap-2">
            <TypeBadge type={subcommit.type} size="md" />
            {subcommit.epic && (
              <Badge variant="outline" className="text-xs">
                {subcommit.epic}
              </Badge>
            )}
          </div>
          <SheetTitle className="text-left text-lg">
            {subcommit.title}
          </SheetTitle>
        </SheetHeader>

        <div className="mt-6 space-y-6">
          {/* Idea */}
          {subcommit.idea && (
            <div>
              <h4 className="mb-1 text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                Idea
              </h4>
              <p className="text-sm leading-relaxed">{subcommit.idea}</p>
            </div>
          )}

          {/* Description */}
          {subcommit.description && (
            <div>
              <h4 className="mb-1 text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                Description
              </h4>
              <p className="text-sm leading-relaxed whitespace-pre-wrap">
                {subcommit.description}
              </p>
            </div>
          )}

          <Separator />

          {/* Commit SHA */}
          <div className="flex items-center gap-2">
            <GitCommit className="h-4 w-4 text-muted-foreground" />
            {commitUrl ? (
              <a
                href={commitUrl}
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center gap-1 font-mono text-xs hover:underline"
              >
                {subcommit.commitSha.slice(0, 7)}
                <ExternalLink className="h-3 w-3" />
              </a>
            ) : (
              <span className="font-mono text-xs">
                {subcommit.commitSha.slice(0, 7)}
              </span>
            )}
          </div>

          {/* Files */}
          {subcommit.files?.length > 0 && (
            <div>
              <h4 className="mb-2 text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                Files ({subcommit.files.length})
              </h4>
              <div className="space-y-1">
                {subcommit.files.map((file, i) => (
                  <div
                    key={i}
                    className="flex items-center gap-2 rounded px-2 py-1 text-xs font-mono hover:bg-accent"
                  >
                    <FileText className="h-3 w-3 shrink-0 text-muted-foreground" />
                    <span className="truncate">{file}</span>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Date */}
          {subcommit.createdAt && (
            <div className="text-xs text-muted-foreground">
              {new Date(subcommit.createdAt).toLocaleString("en-US", {
                dateStyle: "medium",
                timeStyle: "short",
              })}
            </div>
          )}

          {/* Sibling subcommits */}
          {siblings.length > 0 && (
            <>
              <Separator />
              <div>
                <h4 className="mb-2 text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                  From the same commit ({siblings.length})
                </h4>
                <div className="space-y-2">
                  {siblings.map((sib) => (
                    <button
                      key={sib.id}
                      onClick={() => onSiblingClick(sib)}
                      className="flex w-full items-center gap-2 rounded-lg border p-3 text-left transition-colors hover:bg-accent"
                    >
                      <TypeBadge type={sib.type} />
                      <span className="text-sm font-medium truncate">
                        {sib.title}
                      </span>
                    </button>
                  ))}
                </div>
              </div>
            </>
          )}
        </div>
      </SheetContent>
    </Sheet>
  );
}
