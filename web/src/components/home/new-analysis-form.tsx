"use client";

import { useRef, useState } from "react";
import { Loader2, Zap } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useRepoSearch } from "@/lib/hooks/use-repo-search";

interface NewAnalysisFormProps {
  onAnalyze: (repoUrl: string) => void;
  isAnalyzing: boolean;
}

export function NewAnalysisForm({ onAnalyze, isAnalyzing }: NewAnalysisFormProps) {
  const [value, setValue] = useState("");
  const [showSuggestions, setShowSuggestions] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);
  const { results, isLoading } = useRepoSearch(value);

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (value.trim()) {
      setShowSuggestions(false);
      onAnalyze(value.trim());
    }
  }

  function handleSelect(repoUrl: string) {
    setValue(repoUrl);
    setShowSuggestions(false);
    onAnalyze(repoUrl);
  }

  function handleBlur(e: React.FocusEvent) {
    // Only close if focus moved outside the container
    if (!containerRef.current?.contains(e.relatedTarget as Node)) {
      setShowSuggestions(false);
    }
  }

  return (
    <form onSubmit={handleSubmit} className="flex items-center gap-3">
      <div ref={containerRef} className="relative flex-1" onBlur={handleBlur}>
        <Input
          type="text"
          placeholder="Search your repos or paste a URL..."
          value={value}
          onChange={(e) => {
            setValue(e.target.value);
            setShowSuggestions(true);
          }}
          onFocus={() => {
            if (value.length >= 2) setShowSuggestions(true);
          }}
          disabled={isAnalyzing}
        />

        {showSuggestions && value.length >= 2 && (
          <div className="absolute top-full left-0 z-50 mt-1 w-full rounded-md border bg-popover shadow-md">
            {isLoading && (
              <div className="flex items-center justify-center py-3">
                <Loader2 className="h-4 w-4 animate-spin text-muted-foreground" />
              </div>
            )}
            {!isLoading && results.length === 0 && (
              <p className="py-3 text-center text-sm text-muted-foreground">
                No repositories found
              </p>
            )}
            {!isLoading && results.length > 0 && (
              <div className="max-h-60 overflow-y-auto py-1">
                {results.map((repo) => (
                  <button
                    key={repo.id}
                    type="button"
                    onMouseDown={(e) => e.preventDefault()}
                    onClick={() => handleSelect(repo.url)}
                    className="flex w-full flex-col gap-0.5 px-3 py-2 text-left hover:bg-accent"
                  >
                    <span className="text-sm font-medium">{repo.name}</span>
                    <span className="text-xs text-muted-foreground truncate">
                      {repo.url}
                    </span>
                  </button>
                ))}
              </div>
            )}
          </div>
        )}
      </div>

      <Button type="submit" disabled={!value.trim() || isAnalyzing} className="gap-2">
        {isAnalyzing ? (
          <Loader2 className="h-4 w-4 animate-spin" />
        ) : (
          <Zap className="h-4 w-4" />
        )}
        Analyze
      </Button>
    </form>
  );
}
