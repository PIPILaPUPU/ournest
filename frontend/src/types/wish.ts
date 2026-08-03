export type Wish = {
  id: number;
  owner_id: number;
  owner_username: string;
  title: string;
  description: string | null;
  url: string | null;
  price: number | null;
  status: "wanted" | "reserved" | "done";
  group_name: string;
  group_color: GroupColor;
  created_at: string;
  updated_at: string;
};

export type GroupColor = "slate" | "blue" | "green" | "violet" | "rose";

export type CreateWishPayload = {
  title: string;
  description?: string;
  url?: string;
  price?: number;
  status?: "wanted" | "reserved" | "done";
  group_name?: string;
  group_color?: GroupColor;
};

export type UpdateWishPayload = {
  description?: string;
  url?: string;
  price?: number;
};
