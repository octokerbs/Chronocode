"use client";

import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";

interface AnalysisErrorDialogProps {
  open: boolean;
  error: string | null;
  onClose: () => void;
}

export function AnalysisErrorDialog({
  open,
  error,
  onClose,
}: AnalysisErrorDialogProps) {
  return (
    <Dialog open={open} onOpenChange={() => onClose()}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Analysis Failed</DialogTitle>
        </DialogHeader>
        <div className="space-y-4">
          <p className="text-sm text-destructive">{error}</p>
          <Button variant="outline" onClick={onClose}>
            Close
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
