export type DateIdeaStatus = "planned" | "done";

export type DateIdea = {
  id: number;
  author_id: number;
  author_username: string;
  title: string;
  description: string | null;
  status: DateIdeaStatus;
  created_at: string;
  updated_at: string;
};

export type CreateDateIdeaPayload = {
  author_id: number;
  title: string;
  description?: string;
};

export type UpdateDateIdeaPayload = {
  description?: string;
  status?: DateIdeaStatus;
};
