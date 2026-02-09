"use client";

import { useRouter } from "next/navigation";
import { Check, Loader2, ArrowRight } from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import type { AnalysisStep } from "@/lib/hooks/use-analysis";

interface AnalysisProgressModalProps {
  open: boolean;
  step: AnalysisStep;
  repoId: number | null;
  error: string | null;
  onClose: () => void;
}

const steps = [
  { key: "preparing" as const, label: "Preparing repository..." },
  { key: "analyzing" as const, label: "Analyzing commits with AI..." },
  { key: "ready" as const, label: "Timeline ready" },
];

function getStepStatus(
  stepKey: string,
  currentStep: AnalysisStep,
): "done" | "active" | "pending" {
  const order = ["preparing", "analyzing", "ready"];
  const currentIdx = order.indexOf(currentStep);
  const stepIdx = order.indexOf(stepKey);

  if (stepIdx < currentIdx) return "done";
  if (stepIdx === currentIdx) return "active";
  return "pending";
}

export function AnalysisProgressModal({
  open,
  step,
  repoId,
  error,
  onClose,
}: AnalysisProgressModalProps) {
  const router = useRouter();

  return (
    <Dialog open={open} onOpenChange={() => onClose()}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Analyzing Repository</DialogTitle>
        </DialogHeader>

        {error ? (
          <div className="space-y-4">
            <p className="text-sm text-destructive">{error}</p>
            <Button variant="outline" onClick={onClose}>
              Close
            </Button>
          </div>
        ) : (
          <div className="space-y-6 py-4">
            {steps.map((s) => {
              const status = getStepStatus(s.key, step);
              return (
                <div key={s.key} className="flex items-center gap-3">
                  {status === "done" && (
                    <div className="flex h-6 w-6 items-center justify-center rounded-full bg-primary">
                      <Check className="h-3.5 w-3.5 text-primary-foreground" />
                    </div>
                  )}
                  {status === "active" && (
                    <div className="flex h-6 w-6 items-center justify-center">
                      <Loader2 className="h-5 w-5 animate-spin" />
                    </div>
                  )}
                  {status === "pending" && (
                    <div className="h-6 w-6 rounded-full border-2 border-muted" />
                  )}
                  <span
                    className={
                      status === "pending"
                        ? "text-sm text-muted-foreground"
                        : "text-sm font-medium"
                    }
                  >
                    {s.label}
                  </span>
                </div>
              );
            })}

            {step === "ready" && repoId && (
              <Button
                className="w-full gap-2"
                onClick={() => {
                  onClose();
                  router.push(`/timeline/${repoId}`);
                }}
              >
                View Timeline
                <ArrowRight className="h-4 w-4" />
              </Button>
            )}
          </div>
        )}
      </DialogContent>
    </Dialog>
  );
}
