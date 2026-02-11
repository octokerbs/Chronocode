"use client";

import { useState } from "react";
import { Check, Search, Loader2 } from "lucide-react";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useRepoSearch } from "@/lib/hooks/use-repo-search";

interface RepoSearchComboboxProps {
  onSelect: (repoUrl: string) => void;
}

export function RepoSearchCombobox({ onSelect }: RepoSearchComboboxProps) {
  const [open, setOpen] = useState(false);
  const [query, setQuery] = useState("");
  const { results, isLoading } = useRepoSearch(query);

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button variant="outline" size="sm" className="gap-2">
          <Search className="h-3.5 w-3.5" />
          Search your repos
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-80 p-0" align="end">
        <div className="p-2">
          <Input
            placeholder="Search repositories..."
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            className="h-8"
          />
        </div>
        <div className="max-h-60 overflow-y-auto">
          {isLoading && (
            <div className="flex items-center justify-center py-4">
              <Loader2 className="h-4 w-4 animate-spin text-muted-foreground" />
            </div>
          )}
          {!isLoading && results.length === 0 && query.length >= 2 && (
            <p className="py-4 text-center text-sm text-muted-foreground">
              No repositories found
            </p>
          )}
          {results.map((repo) => (
            <button
              key={repo.id}
              onClick={() => {
                onSelect(repo.url);
                setOpen(false);
                setQuery("");
              }}
              className="flex w-full items-center gap-2 px-3 py-2 text-left text-sm hover:bg-accent"
            >
              <Check className="h-3.5 w-3.5 opacity-0" />
              <span className="truncate">{repo.name}</span>
            </button>
          ))}
        </div>
      </PopoverContent>
    </Popover>
  );
}
