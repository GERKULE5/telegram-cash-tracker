import { api } from '../common/networking/axios';
import type {
  CreateSpendingDTO,
  SpendingEntity,
  SpendingFilters,
  SpendingWithUsernameEntity,
  UpdateSpendingDTO,
} from './spending.entity';

export const findAllSpendings = (
  offset: number,
  limit: number,
  filters: SpendingFilters = {},
  signal?: AbortSignal,
) =>
  api.get<SpendingWithUsernameEntity[]>('/admin/spendings', {
    params: { offset, limit, ...filters },
    signal,
  });

export const createSpending = (data: CreateSpendingDTO) =>
  api.post<SpendingEntity>('/admin/spendings', data);

export const updateSpending = (id: number, data: UpdateSpendingDTO) =>
  api.patch<SpendingEntity>(`/admin/spendings/${id}`, data);

export const deleteSpending = (id: number) =>
  api.delete<void>(`/admin/spendings/${id}`);
