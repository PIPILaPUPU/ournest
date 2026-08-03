export type DateIdeaStatus = "planned" | "done";

export type DateIdea = {
  id: number;
  author_id: number;
  author_username: string;
  title: string;
  description: string | null;
  status: DateIdeaStatus;
  is_secret: boolean;
  created_at: string;
  updated_at: string;
};

export type CreateDateIdeaPayload = {
  author_id: number;
  title: string;
  description?: string;
  is_secret?: boolean;
};

export type UpdateDateIdeaPayload = {
  description?: string;
  status?: DateIdeaStatus;
};
