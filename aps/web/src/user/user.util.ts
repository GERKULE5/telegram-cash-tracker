import type { UserEntity } from './user.entity';

export const getUserLabel = (user: UserEntity) => {
  if (user.username) return `@${user.username}`;
  if (user.name) return `${user.name} (#${user.id})`;
  return `Участник #${user.id}`;
};
