import { useEffect, useState } from "react";
import type {
  CreateSpendingDTO,
  SpendingFilters,
  SpendingWithUsernameEntity,
} from "./spending.entity";
import {
  findAllSpendings,
  createSpending,
  updateSpending,
  deleteSpending,
} from "./spending.service";
import { rubleFormatter } from "./spending.format";
import { SpendingFiltersForm } from "./SpendingFiltersForm";
import { SpendingForm } from "./SpendingForm";
import { SpendingTable } from "./SpendingTable";
import "./index.css";

const PAGE_SIZE = 20;

export function App() {
  const [spendings, setSpendings] = useState<SpendingWithUsernameEntity[]>([]);
  const [filters, setFilters] = useState<SpendingFilters>({});
  const [offset, setOffset] = useState(0);
  const [hasNextPage, setHasNextPage] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [loadError, setLoadError] = useState("");
  const [reloadCount, setReloadCount] = useState(0);

  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingSpending, setEditingSpending] =
    useState<SpendingWithUsernameEntity | null>(null);
  const [deletingID, setDeletingID] = useState<number | null>(null);
  const [deleteError, setDeleteError] = useState("");

  useEffect(() => {
    const controller = new AbortController();

    async function loadSpendings() {
      setIsLoading(true);
      setLoadError("");

      try {
        // One extra record tells us whether the next page exists.
        const response = await findAllSpendings(
          offset,
          PAGE_SIZE + 1,
          filters,
          controller.signal,
        );

        if (controller.signal.aborted) {
          return;
        }

        const records = response.data ?? [];
        setSpendings(records.slice(0, PAGE_SIZE));
        setHasNextPage(records.length > PAGE_SIZE);
      } catch (error) {
        if (!controller.signal.aborted) {
          setLoadError(
            error instanceof Error
              ? error.message
              : "Could not load spendings.",
          );
        }
      } finally {
        if (!controller.signal.aborted) {
          setIsLoading(false);
        }
      }
    }

    void loadSpendings();

    return () => controller.abort();
  }, [filters, offset, reloadCount]);

  function refreshSpendings() {
    setDeleteError("");
    setReloadCount((count) => count + 1);
  }

  function applyFilters(nextFilters: SpendingFilters) {
    setFilters(nextFilters);
    setOffset(0);
    setDeleteError("");
  }

  function openCreateForm() {
    setEditingSpending(null);
    setIsFormOpen(true);
  }

  function openEditForm(spending: SpendingWithUsernameEntity) {
    setEditingSpending(spending);
    setIsFormOpen(true);
  }

  function closeForm() {
    setIsFormOpen(false);
    setEditingSpending(null);
  }

  async function saveSpending(values: CreateSpendingDTO) {
    if (editingSpending) {
      const { user_id, ...changes } = values;
      await updateSpending(editingSpending.id, changes);
    } else {
      await createSpending(values);
    }

    closeForm();
    refreshSpendings();
  }

  async function removeSpending(spending: SpendingWithUsernameEntity) {
    const confirmed = window.confirm(
      `Удалить трату на ${rubleFormatter.format(spending.amount)}?`,
    );

    if (!confirmed) {
      return;
    }

    setDeletingID(spending.id);
    setDeleteError("");

    try {
      await deleteSpending(spending.id);

      if (spendings.length === 1 && offset > 0) {
        setOffset(offset - PAGE_SIZE);
      } else {
        refreshSpendings();
      }
    } catch (error) {
      setDeleteError(
        error instanceof Error ? error.message : "Could not delete spending.",
      );
    } finally {
      setDeletingID(null);
    }
  }

  return (
    <main>
      <div className="page-heading">
        <h1>Расходы семьи</h1>
        <button className="primary" onClick={openCreateForm}>
          <span aria-hidden="true">＋</span> Добавить трату
        </button>
      </div>

      <SpendingFiltersForm onApply={applyFilters} />

      <section
        className="ledger"
        aria-label="Список расходов"
        aria-busy={isLoading}
      >
        <div className="ledger-heading">
          <div>
            <h2>История трат</h2>
            <p>Сначала недавно добавленные</p>
          </div>
          <button
            className="text-button"
            disabled={isLoading}
            onClick={refreshSpendings}
          >
            Обновить ↻
          </button>
        </div>

        {deleteError && (
          <div className="error" role="alert">
            {deleteError}
          </div>
        )}

        {isLoading && (
          <div className="empty" role="status">
            Загружаем расходы…
          </div>
        )}

        {!isLoading && loadError && (
          <>
            <div className="error" role="alert">
              {loadError}
            </div>
            <div className="empty">
              Не удалось показать расходы. Попробуйте обновить список.
            </div>
          </>
        )}

        {!isLoading && !loadError && (
          <SpendingTable
            spendings={spendings}
            deletingID={deletingID}
            offset={offset}
            pageSize={PAGE_SIZE}
            hasNextPage={hasNextPage}
            onEdit={openEditForm}
            onDelete={removeSpending}
            onPageChange={setOffset}
          />
        )}
      </section>

      <footer>
        Семейный счёт
        <span>Суммы в рублях · Даты в вашем часовом поясе</span>
      </footer>

      {isFormOpen && (
        <SpendingForm
          spending={editingSpending}
          onSave={saveSpending}
          onClose={closeForm}
        />
      )}
    </main>
  );
}
