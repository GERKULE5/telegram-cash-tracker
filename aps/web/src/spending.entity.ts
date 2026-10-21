export const categories = {
  food: "Еда",
  transport: "Транспорт",
  housing: "Жильё",
  entertainment: "Развлечения",
  health: "Здоровье",
  purchases: "Покупки",
  other: "Другое",
};

export type Category = keyof typeof categories;

export interface SpendingEntity {
  id: number;
  user_id: number;
  amount: number;
  category: Category;
  spent_at: string;
  comment: string | null;
  created_at: string;
  updated_at: string | null;
}

export interface SpendingWithUsernameEntity extends SpendingEntity {
  username: string | null;
}

export interface CreateSpendingDTO {
  user_id: number;
  amount: number;
  category: Category;
  spent_at: string;
  comment: string;
}

export interface UpdateSpendingDTO {
  amount?: number;
  category?: Category;
  spent_at?: string;
  comment?: string;
}

export interface SpendingFilters {
  user_id?: number;
  category?: Category;
  from?: string;
  to?: string;
}
