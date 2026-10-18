package spending

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func mapSpending(row pgx.Row) (Spending, error) {
	var result Spending
	err := row.Scan(&result.ID, &result.UserID, &result.Amount, &result.Category,
		&result.SpentAt, &result.Comment, &result.CreatedAt, &result.UpdatedAt)
	return result, err
}

func (r *PostgresRepository) Create(ctx context.Context, params CreateParams) (Spending, error) {
	result, err := mapSpending(r.db.QueryRow(ctx, `
		INSERT INTO spendings (user_id, amount, category, spent_at, comment)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, amount, category, spent_at, comment, created_at, updated_at
	`, params.UserID, params.Amount, params.Category, params.SpentAt, params.Comment))
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return Spending{}, ErrUserNotFound
		}
		return Spending{}, fmt.Errorf("create spending: %w", err)
	}
	return result, nil
}

func (r *PostgresRepository) List(ctx context.Context, params ListParams) ([]Spending, error) {
	return r.list(ctx, params, &params.UserID)
}

func (r *PostgresRepository) ListAll(ctx context.Context, params ListParams) ([]Spending, error) {
	if params.UserID != 0 {
		return r.list(ctx, params, &params.UserID)
	}
	return r.list(ctx, params, nil)
}

func (r *PostgresRepository) list(ctx context.Context, params ListParams, userID *int) ([]Spending, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, amount, category, spent_at, comment, created_at, updated_at
		FROM spendings
		WHERE ($1::integer IS NULL OR user_id = $1)
		  AND ($2::text = '' OR category::text = $2)
		  AND ($3::timestamptz IS NULL OR spent_at >= $3)
		  AND ($4::timestamptz IS NULL OR spent_at < $4)
		ORDER BY created_at DESC, id DESC
		LIMIT $5 OFFSET $6
	`, userID, params.Category, params.From, params.To, params.Limit, params.Offset)
	if err != nil {
		return nil, fmt.Errorf("list spendings: %w", err)
	}
	defer rows.Close()

	result := make([]Spending, 0)
	for rows.Next() {
		item, err := mapSpending(rows)
		if err != nil {
			return nil, fmt.Errorf("scan spending: %w", err)
		}
		result = append(result, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read spendings: %w", err)
	}
	return result, nil
}

func (r *PostgresRepository) Update(ctx context.Context, id, userID int, params UpdateParams) (Spending, error) {
	return r.update(ctx, id, &userID, params)
}

func (r *PostgresRepository) UpdateByID(ctx context.Context, id int, params UpdateParams) (Spending, error) {
	return r.update(ctx, id, nil, params)
}

func (r *PostgresRepository) update(ctx context.Context, id int, userID *int, params UpdateParams) (Spending, error) {
	result, err := mapSpending(r.db.QueryRow(ctx, `
		UPDATE spendings
		SET amount = COALESCE($3::integer, amount),
		    category = COALESCE($4::category, category),
		    spent_at = COALESCE($5::timestamptz, spent_at),
		    comment = COALESCE($6::text, comment),
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND ($2::integer IS NULL OR user_id = $2)
		RETURNING id, user_id, amount, category, spent_at, comment, created_at, updated_at
	`, id, userID, params.Amount, params.Category, params.SpentAt, params.Comment))
	if errors.Is(err, pgx.ErrNoRows) {
		return Spending{}, ErrNotFound
	}
	if err != nil {
		return Spending{}, fmt.Errorf("update spending: %w", err)
	}
	return result, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id, userID int) error {
	return r.delete(ctx, id, &userID)
}

func (r *PostgresRepository) DeleteByID(ctx context.Context, id int) error {
	return r.delete(ctx, id, nil)
}

func (r *PostgresRepository) delete(ctx context.Context, id int, userID *int) error {
	result, err := r.db.Exec(ctx, `DELETE FROM spendings WHERE id = $1 AND ($2::integer IS NULL OR user_id = $2)`, id, userID)
	if err != nil {
		return fmt.Errorf("delete spending: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id int) (Spending, error) {
	result, err := mapSpending(r.db.QueryRow(ctx, `
		SELECT id, user_id, amount, category, spent_at, comment, created_at, updated_at
		FROM spendings WHERE id = $1
	`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return Spending{}, ErrNotFound
	}
	if err != nil {
		return Spending{}, fmt.Errorf("get spending: %w", err)
	}
	return result, nil
}
