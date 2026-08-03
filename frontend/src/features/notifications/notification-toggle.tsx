import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Bell, BellOff } from "lucide-react";
import { Button } from "@/components/ui/button";
import { getPushConfig } from "@/features/notifications/api";
import {
  disablePush,
  enablePush,
  getPushStatus,
  type PushStatus,
} from "@/features/notifications/push";

const pushConfigKey = ["push-config"];
const pushStatusKey = ["push-status"];

export function NotificationToggle() {
  const queryClient = useQueryClient();
  const config = useQuery({
    queryKey: pushConfigKey,
    queryFn: getPushConfig,
    staleTime: 5 * 60 * 1000,
  });
  const status = useQuery({
    queryKey: pushStatusKey,
    queryFn: getPushStatus,
  });

  const toggle = useMutation({
    mutationFn: async () => {
      if (status.data === "enabled") {
        await disablePush();
        return;
      }
      if (!config.data?.enabled || !config.data.public_key) {
        throw new Error("push is not configured");
      }
      await enablePush(config.data.public_key);
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: pushStatusKey }),
  });

  const currentStatus: PushStatus = status.data ?? "disabled";
  const enabled = currentStatus === "enabled";
  const unavailable =
    currentStatus === "unsupported" ||
    currentStatus === "denied" ||
    config.isError ||
    (!enabled && config.data?.enabled === false);

  return (
    <Button
      type="button"
      variant="outline"
      onClick={() => toggle.mutate()}
      disabled={status.isLoading || config.isLoading || toggle.isPending || unavailable}
      title={statusTitle(currentStatus, config.data?.enabled ?? false, toggle.isError)}
    >
      {enabled ? (
        <Bell className="mr-2 h-4 w-4" />
      ) : (
        <BellOff className="mr-2 h-4 w-4" />
      )}
      <span className="hidden sm:inline">
        {enabled ? "Уведомления включены" : "Включить уведомления"}
      </span>
      <span className="sm:hidden">Push</span>
    </Button>
  );
}

function statusTitle(status: PushStatus, configured: boolean, failed: boolean): string {
  if (failed) return "Не удалось изменить настройки уведомлений";
  if (!configured) return "Push-уведомления не настроены на сервере";
  if (status === "unsupported") return "Браузер не поддерживает Web Push";
  if (status === "denied") return "Уведомления запрещены в настройках браузера";
  if (status === "enabled") return "Нажмите, чтобы отключить уведомления";
  return "Нажмите, чтобы получать уведомления";
}
