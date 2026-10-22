import { useState } from 'react';
import type { UserEntity } from '../../user/user.entity';
import { getUserLabel } from '../../user/user.util';
import type { FormEvent } from 'react';
import { categories } from '../spending.entity';
import type { Category, SpendingFilters } from '../spending.entity';

interface SpendingFiltersFormProps {
  users: UserEntity[];
  usersLoading: boolean;
  onApply: (filters: SpendingFilters) => void;
}

export const SpendingFiltersForm = ({
  users,
  usersLoading,
  onApply,
}: SpendingFiltersFormProps) => {
  const [error, setError] = useState('');

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    const data = new FormData(event.currentTarget);
    const userID = String(data.get('user_id') ?? '');
    const category = String(data.get('category') ?? '');
    const from = String(data.get('from') ?? '');
    const to = String(data.get('to') ?? '');

    if (from && to && from > to) {
      setError('The start date must not be after the end date.');
      return;
    }

    const filters: SpendingFilters = {};

    if (userID) {
      filters.user_id = Number(userID);
    }

    if (category) {
      filters.category = category as Category;
    }

    if (from) {
      filters.from = new Date(`${from}T00:00:00`).toISOString();
    }

    if (to) {
      // The API excludes the end date, so use midnight of the next day.
      const endDate = new Date(`${to}T00:00:00`);
      endDate.setDate(endDate.getDate() + 1);
      filters.to = endDate.toISOString();
    }

    setError('');
    onApply(filters);
  }

  function handleReset() {
    setError('');
    onApply({});
  }

  return (
    <>
      <form className="filters" onSubmit={handleSubmit} onReset={handleReset}>
        <label>
          Член семьи
          <select name="user_id" disabled={usersLoading}>
            <option value="">
              {usersLoading ? 'Загружаем участников…' : 'Все участники'}
            </option>
            {users.map(user => (
              <option key={user.id} value={user.id}>
                {getUserLabel(user)}
              </option>
            ))}
          </select>
        </label>

        <label>
          Категория
          <select name="category">
            <option value="">Все категории</option>
            {Object.entries(categories).map(([value, label]) => (
              <option key={value} value={value}>
                {label}
              </option>
            ))}
          </select>
        </label>

        <label>
          С даты
          <input name="from" type="date" />
        </label>
        <label>
          По дату
          <input name="to" type="date" />
        </label>

        <div className="filter-actions">
          <button type="submit">Применить</button>
          <button type="reset" className="text-button">
            Сбросить
          </button>
        </div>
      </form>

      {error && (
        <div className="error" role="alert">
          {error}
        </div>
      )}
    </>
  );
};
