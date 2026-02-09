"use client";

import Link from "next/link";
import { Card, CardContent } from "@/components/ui/card";
import { ArrowRight } from "lucide-react";
import type { UserRepository } from "@/lib/types";

interface RepositoryCardProps {
  repo: UserRepository;
}

export function RepositoryCard({ repo }: RepositoryCardProps) {
  return (
    <Link href={`/timeline/${repo.id}`}>
      <Card className="group cursor-pointer transition-all hover:border-foreground/20 hover:shadow-md">
        <CardContent className="flex items-center justify-between p-5">
          <div className="min-w-0 flex-1">
            <h3 className="truncate font-semibold">{repo.name}</h3>
            {repo.addedAt && (
              <p className="mt-1 text-xs text-muted-foreground">
                Analyzed{" "}
                {new Date(repo.addedAt).toLocaleDateString("en-US", {
                  month: "short",
                  day: "numeric",
                  year: "numeric",
                })}
              </p>
            )}
          </div>
          <ArrowRight className="h-4 w-4 shrink-0 text-muted-foreground transition-transform group-hover:translate-x-0.5" />
        </CardContent>
      </Card>
    </Link>
  );
}
