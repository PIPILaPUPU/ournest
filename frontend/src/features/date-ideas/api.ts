import { apiClient } from "@/lib/api-client";
import axios from "axios";
import type {
  CreateDateIdeaPayload,
  DateIdea,
  UpdateDateIdeaPayload,
} from "@/features/date-ideas/types";

export async function getDateIdeas(): Promise<DateIdea[]> {
  const response = await apiClient.get<DateIdea[]>("/date-ideas");
  return response.data;
}

export async function getRandomSecretDateIdea(): Promise<DateIdea | null> {
  try {
    const response = await apiClient.get<DateIdea>("/date-ideas/secret/random");
    return response.data;
  } catch (error) {
    if (axios.isAxiosError(error) && error.response?.status === 404) {
      return null;
    }
    throw error;
  }
}

export async function createDateIdea(payload: CreateDateIdeaPayload): Promise<DateIdea> {
  const response = await apiClient.post<DateIdea>("/date-ideas", payload);
  return response.data;
}

export async function updateDateIdea(
  id: number,
  payload: UpdateDateIdeaPayload,
): Promise<DateIdea> {
  const response = await apiClient.patch<DateIdea>(`/date-ideas/${id}`, payload);
  return response.data;
}

export async function deleteDateIdea(id: number): Promise<void> {
  await apiClient.delete(`/date-ideas/${id}`);
}
