"use client";

import { GitBranch, Brain, BarChart3 } from "lucide-react";
import { TimelineAnimation } from "./timeline-animation";

const features = [
  {
    icon: Brain,
    title: "AI-Powered Analysis",
    description: "Understand your codebase evolution with AI-driven commit analysis",
  },
  {
    icon: GitBranch,
    title: "Visual Timeline",
    description: "See your project's story unfold commit by commit",
  },
  {
    icon: BarChart3,
    title: "Team Insights",
    description: "Understand what your team is building at a glance",
  },
];

export function MarketingPanel() {
  return (
    <div className="flex h-full flex-col justify-between p-12">
      <div>
        <h1 className="text-4xl font-bold tracking-tight">Chronocode</h1>
        <p className="mt-2 text-lg text-muted-foreground">
          The story behind every commit
        </p>
      </div>

      <TimelineAnimation />

      <div className="space-y-6">
        {features.map((feature) => (
          <div key={feature.title} className="flex items-start gap-4">
            <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg border bg-card">
              <feature.icon className="h-5 w-5" />
            </div>
            <div>
              <h3 className="font-semibold">{feature.title}</h3>
              <p className="text-sm text-muted-foreground">
                {feature.description}
              </p>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
