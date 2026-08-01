import axios from "axios";
import type { LoginPayload, LoginResponse } from "@/types/auth";
import type { CreateWishPayload, UpdateWishPayload, Wish } from "@/types/wish";

const baseURL = import.meta.env.VITE_API_URL ?? "http://localhost:8080";

export const api = axios.create({
  baseURL,
  timeout: 10000,
});

export async function getWishlist(): Promise<Wish[]> {
  const response = await api.get<Wish[]>("/wishlist");
  return response.data;
}

export async function login(payload: LoginPayload): Promise<LoginResponse> {
  const response = await api.post<LoginResponse>("/auth/login", payload);
  return response.data;
}

export async function createWish(payload: CreateWishPayload): Promise<Wish> {
  const response = await api.post<Wish>("/wishlist", payload);
  return response.data;
}

export async function updateWish(id: number, payload: UpdateWishPayload): Promise<Wish> {
  const response = await api.patch<Wish>(`/wishlist/${id}`, payload);
  return response.data;
}

export async function deleteWish(id: number): Promise<void> {
  await api.delete(`/wishlist/${id}`);
}
