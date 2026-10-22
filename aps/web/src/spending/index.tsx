import { useCallback, useEffect, useRef, useState } from 'react';
import type { UserEntity } from '../user/user.entity';
import { findAllUsers } from '../user/user.service';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlus, faRotateRight } from '@fortawesome/free-solid-svg-icons';
import type {
  CreateSpendingDTO,
  SpendingFilters,
  SpendingWithUsernameEntity,
} from './spending.entity';
import {
  findAllSpendings,
  createSpending,
  updateSpending,
  deleteSpending,
} from './spending.service';
import { rubleFormatter } from './spending.util';
import { SpendingFiltersForm } from './component/spending-filters.component';
import { SpendingModal } from './component/spending-modal.component';
import { SpendingTable } from './component/spending-table.component';

const PAGE_SIZE = 20;

const SpendingModule = () => {
  const [users, setUsers] = useState<UserEntity[]>([]);
  const [usersLoading, setUsersLoading] = useState(true);
  const [usersError, setUsersError] = useState('');
  const [spendings, setSpendings] = useState<SpendingWithUsernameEntity[]>([]);
  const [filters, setFilters] = useState<SpendingFilters>({});
  const [offset, setOffset] = useState(0);
  const [hasNextPage, setHasNextPage] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [loadError, setLoadError] = useState('');
  const activeRequest = useRef<AbortController | null>(null);

  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingSpending, setEditingSpending] =
    useState<SpendingWithUsernameEntity | null>(null);
  const [deletingID, setDeletingID] = useState<number | null>(null);
  const [deleteError, setDeleteError] = useState('');

  const loadUsers = useCallback(async (signal?: AbortSignal) => {
    setUsersLoading(true);
    setUsersError('');

    try {
      const response = await findAllUsers(signal);
      if (!signal?.aborted) setUsers(response.data);
    } catch (error) {
      if (!signal?.aborted) {
        setUsersError(
          error instanceof Error ? error.message : 'Could not load users.',
        );
      }
    } finally {
      if (!signal?.aborted) setUsersLoading(false);
    }
  }, []);

  useEffect(() => {
    const controller = new AbortController();
    void loadUsers(controller.signal);
    return () => controller.abort();
  }, [loadUsers]);

  const loadSpendings = useCallback(async () => {
    // Cancel the previous load so an older response cannot replace newer data.
    activeRequest.current?.abort();
    const controller = new AbortController();
    activeRequest.current = controller;

    setIsLoading(true);
    setLoadError('');

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
          error instanceof Error ? error.message : 'Could not load spendings.',
        );
      }
    } finally {
      if (!controller.signal.aborted) {
        setIsLoading(false);
      }
    }
  }, [filters, offset]);

  useEffect(() => {
    void loadSpendings();

    return () => activeRequest.current?.abort();
  }, [loadSpendings]);

  function handleRefresh() {
    setDeleteError('');
    void loadSpendings();
  }

  function handleApplyFilters(nextFilters: SpendingFilters) {
    setFilters(nextFilters);
    setOffset(0);
    setDeleteError('');
  }

  function handleCreate() {
    setEditingSpending(null);
    setIsFormOpen(true);
  }

  function handleEdit(spending: SpendingWithUsernameEntity) {
    setEditingSpending(spending);
    setIsFormOpen(true);
  }

  function handleClose() {
    setIsFormOpen(false);
    setEditingSpending(null);
  }

  async function handleSave(values: CreateSpendingDTO) {
    if (editingSpending) {
      const { user_id, ...changes } = values;
      await updateSpending(editingSpending.id, changes);
    } else {
      await createSpending(values);
    }

    handleClose();
    handleRefresh();
  }

  async function handleDelete(spending: SpendingWithUsernameEntity) {
    const confirmed = window.confirm(
      `Удалить трату на ${rubleFormatter.format(spending.amount)}?`,
    );

    if (!confirmed) {
      return;
    }

    setDeletingID(spending.id);
    setDeleteError('');

    try {
      await deleteSpending(spending.id);

      if (spendings.length === 1 && offset > 0) {
        setOffset(offset - PAGE_SIZE);
      } else {
        handleRefresh();
      }
    } catch (error) {
      setDeleteError(
        error instanceof Error ? error.message : 'Could not delete spending.',
      );
    } finally {
      setDeletingID(null);
    }
  }

  return (
    <main>
      <div className="page-heading">
        <h1>Расходы семьи</h1>
        <button
          className="primary"
          onClick={handleCreate}
          disabled={usersLoading || !!usersError || users.length === 0}
        >
          <FontAwesomeIcon icon={faPlus} /> Добавить трату
        </button>
      </div>

      {usersError && (
        <div className="error" role="alert">
          {usersError}{' '}
          <button type="button" onClick={() => void loadUsers()}>
            Повторить
          </button>
        </div>
      )}

      {!usersLoading && !usersError && users.length === 0 && (
        <p>
          Участников пока нет. Сначала напишите Telegram-боту, чтобы
          зарегистрироваться.
        </p>
      )}

      <SpendingFiltersForm
        users={users}
        usersLoading={usersLoading}
        onApply={handleApplyFilters}
      />

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
            onClick={handleRefresh}
          >
            <FontAwesomeIcon icon={faRotateRight} /> Обновить
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
            onEdit={handleEdit}
            onDelete={handleDelete}
            onPageChange={setOffset}
          />
        )}
      </section>

      {isFormOpen && (
        <SpendingModal
          users={users}
          spending={editingSpending}
          onSave={handleSave}
          onClose={handleClose}
        />
      )}
    </main>
  );
};

export default SpendingModule;
