"use client";

import { useEffect, useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { Loader2, Check } from "lucide-react";

interface AnalysisStatusProps {
  isAnalyzing: boolean;
}

export function AnalysisStatus({ isAnalyzing }: AnalysisStatusProps) {
  const [showUpToDate, setShowUpToDate] = useState(false);
  const [wasAnalyzing, setWasAnalyzing] = useState(false);

  useEffect(() => {
    if (isAnalyzing) {
      setWasAnalyzing(true);
      setShowUpToDate(false);
    } else if (wasAnalyzing) {
      setShowUpToDate(true);
      const timer = setTimeout(() => setShowUpToDate(false), 4000);
      return () => clearTimeout(timer);
    }
  }, [isAnalyzing, wasAnalyzing]);

  return (
    <div className="fixed bottom-4 right-4 z-50">
      <AnimatePresence mode="wait">
        {isAnalyzing && (
          <motion.div
            key="analyzing"
            initial={{ opacity: 0, y: 20, scale: 0.95 }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            exit={{ opacity: 0, y: 10, scale: 0.95 }}
            transition={{ duration: 0.2 }}
            className="flex items-center gap-2.5 rounded-lg border bg-background/95 px-4 py-2.5 shadow-lg backdrop-blur"
          >
            <Loader2 className="h-4 w-4 animate-spin text-muted-foreground" />
            <span className="text-sm text-muted-foreground">
              Analyzing commits...
            </span>
            <span className="relative flex h-2 w-2">
              <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-orange-400 opacity-75" />
              <span className="relative inline-flex h-2 w-2 rounded-full bg-orange-500" />
            </span>
          </motion.div>
        )}

        {showUpToDate && !isAnalyzing && (
          <motion.div
            key="up-to-date"
            initial={{ opacity: 0, y: 20, scale: 0.95 }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            exit={{ opacity: 0, y: 10, scale: 0.95 }}
            transition={{ duration: 0.2 }}
            className="flex items-center gap-2.5 rounded-lg border bg-background/95 px-4 py-2.5 shadow-lg backdrop-blur"
          >
            <div className="flex h-4 w-4 items-center justify-center rounded-full bg-green-500">
              <Check className="h-2.5 w-2.5 text-white" />
            </div>
            <span className="text-sm text-muted-foreground">Up to date</span>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
