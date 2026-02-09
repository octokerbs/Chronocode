"use client";

import Link from "next/link";
import { GitBranch } from "lucide-react";
import { UserMenu } from "./user-menu";

export function Navbar() {
  return (
    <header className="sticky top-0 z-50 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="flex h-14 items-center justify-between px-6">
        <Link href="/home" className="flex items-center gap-2">
          <GitBranch className="h-5 w-5" />
          <span className="text-lg font-bold tracking-tight">Chronocode</span>
        </Link>
        <UserMenu />
      </div>
    </header>
  );
}
