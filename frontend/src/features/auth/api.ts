import { apiClient } from "@/lib/api-client";
import type { LoginPayload, LoginResponse } from "@/types/auth";

export async function login(payload: LoginPayload): Promise<LoginResponse> {
  const response = await apiClient.post<LoginResponse>("/auth/login", payload);
  return response.data;
}
