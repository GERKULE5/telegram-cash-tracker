import { useEffect, useRef, useState } from "react";
import type { FormEvent, SyntheticEvent } from "react";
import {
  categories,
  type SpendingEntity,
  type CreateSpendingDTO,
} from "./spending.entity";

function toLocalDateTime(value: Date) {
  // datetime-local expects local time without a timezone suffix.
  const local = new Date(value.getTime() - value.getTimezoneOffset() * 60_000);
  return local.toISOString().slice(0, 19);
}

interface SpendingFormProps {
  spending: SpendingEntity | null;
  onSave: (values: CreateSpendingDTO) => Promise<void>;
  onClose: () => void;
}

export function SpendingForm({ spending, onSave, onClose }: SpendingFormProps) {
  const dialog = useRef<HTMLDialogElement>(null);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    dialog.current?.showModal();
  }, []);

  function handleCancel(event: SyntheticEvent<HTMLDialogElement>) {
    event.preventDefault();

    if (!saving) {
      onClose();
    }
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    const data = new FormData(event.currentTarget);

    setSaving(true);
    setError("");

    try {
      await onSave({
        user_id: spending?.user_id ?? Number(data.get("user_id")),
        amount: Number(data.get("amount")),
        category: data.get("category") as CreateSpendingDTO["category"],
        spent_at: new Date(String(data.get("spent_at"))).toISOString(),
        comment: String(data.get("comment") ?? "").trim(),
      });
    } catch (err) {
      setError(err instanceof Error ? err.message : "Could not save spending.");
    } finally {
      setSaving(false);
    }
  }

  return (
    <dialog ref={dialog} aria-labelledby="form-title" onCancel={handleCancel}>
      <div className="modal-heading">
        <div>
          <p className="eyebrow">СЕМЕЙНЫЙ СЧЁТ</p>
          <h2 id="form-title">{spending ? "Изменить трату" : "Новая трата"}</h2>
        </div>
        <button
          className="text-button"
          aria-label="Закрыть"
          disabled={saving}
          onClick={onClose}
        >
          ✕
        </button>
      </div>

      <form onSubmit={handleSubmit}>
        <fieldset disabled={saving}>
          <label>
            Член семьи · внутренний ID
            <input
              name="user_id"
              type="number"
              required
              min="1"
              max="2147483647"
              defaultValue={spending?.user_id}
              disabled={!!spending}
              placeholder="Например, 1"
            />
          </label>

          <div className="form-grid">
            <label>
              Сумма, ₽
              <input
                autoFocus
                name="amount"
                type="number"
                required
                min="1"
                max="2147483647"
                step="1"
                defaultValue={spending?.amount}
                placeholder="0"
              />
            </label>

            <label>
              Категория
              <select
                name="category"
                defaultValue={spending?.category ?? "food"}
              >
                {Object.entries(categories).map(([key, label]) => (
                  <option key={key} value={key}>
                    {label}
                  </option>
                ))}
              </select>
            </label>
          </div>

          <label>
            Дата и время
            <input
              name="spent_at"
              type="datetime-local"
              step="1"
              required
              defaultValue={toLocalDateTime(
                spending ? new Date(spending.spent_at) : new Date(),
              )}
            />
          </label>

          <label>
            Комментарий <span className="optional">необязательно</span>
            <textarea
              name="comment"
              rows={3}
              defaultValue={spending?.comment ?? ""}
              placeholder="На что потратили?"
            />
          </label>
        </fieldset>

        {error && (
          <div className="error" role="alert">
            {error}
          </div>
        )}

        <div className="modal-actions">
          <button type="button" disabled={saving} onClick={onClose}>
            Отмена
          </button>
          <button className="primary" disabled={saving}>
            {saving ? "Сохраняем…" : "Сохранить"}
          </button>
        </div>
      </form>
    </dialog>
  );
}
