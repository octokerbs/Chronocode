"use client";

import { Monitor } from "lucide-react";
import { useMediaQuery } from "@/lib/hooks/use-media-query";

export function DesktopOnlyGate({ children }: { children: React.ReactNode }) {
  const isDesktop = useMediaQuery("(min-width: 1024px)");

  if (!isDesktop) {
    return (
      <div className="flex h-screen flex-col items-center justify-center gap-4 bg-background p-8 text-center">
        <Monitor className="h-12 w-12 text-muted-foreground" />
        <h1 className="text-2xl font-bold">Best viewed on desktop</h1>
        <p className="max-w-sm text-muted-foreground">
          Chronocode&apos;s timeline visualization requires a larger screen for
          the best experience.
        </p>
      </div>
    );
  }

  return <>{children}</>;
}
