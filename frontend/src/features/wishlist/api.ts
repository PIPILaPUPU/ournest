import { apiClient } from "@/lib/api-client";
import type { CreateWishPayload, UpdateWishPayload, Wish } from "@/types/wish";

export async function getWishlist(): Promise<Wish[]> {
  const response = await apiClient.get<Wish[]>("/wishlist");
  return response.data;
}

export async function createWish(payload: CreateWishPayload): Promise<Wish> {
  const response = await apiClient.post<Wish>("/wishlist", payload);
  return response.data;
}

export async function updateWish(id: number, payload: UpdateWishPayload): Promise<Wish> {
  const response = await apiClient.patch<Wish>(`/wishlist/${id}`, payload);
  return response.data;
}

export async function deleteWish(id: number): Promise<void> {
  await apiClient.delete(`/wishlist/${id}`);
}
