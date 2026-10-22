import { api } from '../common/networking/axios';
import type { UserEntity } from './user.entity';

export const findAllUsers = (signal?: AbortSignal) =>
  api.get<UserEntity[]>('/admin/users', { signal });
