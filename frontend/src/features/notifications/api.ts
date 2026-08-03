import { apiClient } from "@/lib/api-client";

export type PushConfig = {
  enabled: boolean;
  public_key: string;
};

export type PushSubscriptionPayload = {
  endpoint: string;
  keys: {
    p256dh: string;
    auth: string;
  };
};

export async function getPushConfig(): Promise<PushConfig> {
  const response = await apiClient.get<PushConfig>("/notifications/vapid-public-key");
  return response.data;
}

export async function savePushSubscription(
  subscription: PushSubscriptionPayload,
): Promise<void> {
  await apiClient.post("/notifications/subscriptions", subscription);
}

export async function removePushSubscription(endpoint: string): Promise<void> {
  await apiClient.delete("/notifications/subscriptions", {
    data: { endpoint },
  });
}
