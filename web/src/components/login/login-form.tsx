"use client";

import { Github } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/providers/auth-provider";

export function LoginForm() {
  const { login } = useAuth();

  return (
    <div className="flex h-full flex-col items-center justify-center px-12">
      <div className="w-full max-w-sm space-y-8">
        <div className="space-y-2 text-center">
          <h2 className="text-2xl font-bold tracking-tight">
            Sign in to Chronocode
          </h2>
          <p className="text-sm text-muted-foreground">
            Analyze your repositories with AI
          </p>
        </div>

        <Button
          onClick={login}
          size="lg"
          className="w-full gap-2"
        >
          <Github className="h-5 w-5" />
          Sign in with GitHub
        </Button>

        <p className="text-center text-xs text-muted-foreground">
          By signing in, you grant Chronocode read access to your repositories
          for analysis.
        </p>
      </div>
    </div>
  );
}
