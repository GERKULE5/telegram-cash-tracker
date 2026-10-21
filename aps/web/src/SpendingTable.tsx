import { categories } from "./spending.entity";
import type { SpendingWithUsernameEntity } from "./spending.entity";
import { rubleFormatter, dateFormatter } from "./spending.format";

interface SpendingTableProps {
  spendings: SpendingWithUsernameEntity[];
  deletingID: number | null;
  offset: number;
  pageSize: number;
  hasNextPage: boolean;
  onEdit: (spending: SpendingWithUsernameEntity) => void;
  onDelete: (spending: SpendingWithUsernameEntity) => void;
  onPageChange: (offset: number) => void;
}

export function SpendingTable({
  spendings,
  deletingID,
  offset,
  pageSize,
  hasNextPage,
  onEdit,
  onDelete,
  onPageChange,
}: SpendingTableProps) {
  if (spendings.length === 0) {
    return (
      <div className="empty">
        <span className="empty-icon">₽</span>
        <h3>Здесь пока нет расходов</h3>
        <p>Добавьте первую трату или измените фильтры.</p>
      </div>
    );
  }

  const pageTotal = spendings.reduce(
    (sum, spending) => sum + spending.amount,
    0,
  );
  const pageNumber = offset / pageSize + 1;

  return (
    <>
      <div className="table-scroll">
        <table>
          <thead>
            <tr>
              <th>Трата</th>
              <th>Член семьи</th>
              <th>Дата</th>
              <th className="amount">Сумма</th>
              <th>
                <span className="sr-only">Действия</span>
              </th>
            </tr>
          </thead>

          <tbody>
            {spendings.map((spending) => (
              <tr key={spending.id}>
                <td>
                  <div className="spending-title">
                    {spending.comment || categories[spending.category]}
                  </div>
                  <span className={`category ${spending.category}`}>
                    {categories[spending.category] ?? spending.category}
                  </span>
                </td>

                <td>
                  <span className="member-avatar" aria-hidden="true">
                    {spending.user_id}
                  </span>
                  {spending.username
                    ? `@${spending.username}`
                    : `Участник #${spending.user_id}`}
                </td>

                <td className="date">
                  {dateFormatter.format(new Date(spending.spent_at))}
                </td>

                <td className="amount spending-amount">
                  {rubleFormatter.format(spending.amount)}
                </td>

                <td>
                  <div className="row-actions">
                    <button
                      className="text-button"
                      onClick={() => onEdit(spending)}
                    >
                      Изменить
                      <span className="sr-only"> трату #{spending.id}</span>
                    </button>
                    <button
                      className="text-button danger"
                      disabled={deletingID !== null}
                      onClick={() => onDelete(spending)}
                    >
                      {deletingID === spending.id ? "Удаляем…" : "Удалить"}
                      <span className="sr-only"> трату #{spending.id}</span>
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <div className="table-footer">
        <span>
          На этой странице: <strong>{rubleFormatter.format(pageTotal)}</strong>
        </span>
        <div className="pagination">
          <button
            disabled={offset === 0}
            onClick={() => onPageChange(offset - pageSize)}
            aria-label="Предыдущая страница"
          >
            ←
          </button>
          <span>Страница {pageNumber}</span>
          <button
            disabled={!hasNextPage}
            onClick={() => onPageChange(offset + pageSize)}
            aria-label="Следующая страница"
          >
            →
          </button>
        </div>
      </div>
    </>
  );
}
