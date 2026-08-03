import { useCallback, useEffect, useState } from "react";
import { CheckCircle2 } from "lucide-react";

export type SuccessToastState = {
  id: number;
  message: string;
} | null;

export function useSuccessToast() {
  const [toast, setToast] = useState<SuccessToastState>(null);

  const showSuccessToast = useCallback((message: string) => {
    setToast({
      id: Date.now(),
      message,
    });
  }, []);
  const dismissSuccessToast = useCallback(() => setToast(null), []);

  return {
    toast,
    showSuccessToast,
    dismissSuccessToast,
  };
}

export function SuccessToast({
  toast,
  onClose,
}: {
  toast: SuccessToastState;
  onClose: () => void;
}) {
  useEffect(() => {
    if (!toast) return;

    const timeout = window.setTimeout(onClose, 3200);
    return () => window.clearTimeout(timeout);
  }, [toast, onClose]);

  if (!toast) return null;

  return (
    <div
      className="pointer-events-none fixed left-1/2 top-5 z-[100] w-[calc(100%-2rem)] max-w-sm -translate-x-1/2"
      role="status"
      aria-live="polite"
    >
      <div
        key={toast.id}
        className="success-toast flex items-center gap-3 rounded-xl border border-emerald-200 bg-white px-4 py-3 text-sm font-medium text-gray-800 shadow-lg"
      >
        <CheckCircle2 className="h-5 w-5 shrink-0 text-emerald-600" />
        {toast.message}
      </div>
    </div>
  );
}
