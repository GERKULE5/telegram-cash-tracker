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

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
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
	query := `
		SELECT id, user_id, amount, category, spent_at, comment, created_at, updated_at
		FROM spendings
		WHERE user_id = $1
	`
	args := []any{params.UserID}
	if params.Category != "" {
		args = append(args, params.Category)
		query += fmt.Sprintf(" AND category = $%d", len(args))
	}
	if params.From != nil {
		args = append(args, *params.From)
		query += fmt.Sprintf(" AND spent_at >= $%d", len(args))
	}
	if params.To != nil {
		args = append(args, *params.To)
		query += fmt.Sprintf(" AND spent_at < $%d", len(args))
	}
	args = append(args, params.Limit, params.Offset)
	query += fmt.Sprintf(" ORDER BY created_at DESC, id DESC LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	rows, err := r.db.Query(ctx, query, args...)
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

func (r *PostgresRepository) ListAll(ctx context.Context, params ListParams) ([]SpendingWithUsername, error) {
	query := `
		SELECT s.id, s.user_id, s.amount, s.category, s.spent_at, s.comment,
		       s.created_at, s.updated_at, u.username
		FROM spendings s
		JOIN users u ON u.id = s.user_id
		WHERE TRUE
	`
	args := make([]any, 0)
	if params.UserID != 0 {
		args = append(args, params.UserID)
		query += " AND s.user_id = $1"
	}
	if params.Category != "" {
		args = append(args, params.Category)
		query += fmt.Sprintf(" AND s.category = $%d", len(args))
	}
	if params.From != nil {
		args = append(args, *params.From)
		query += fmt.Sprintf(" AND s.spent_at >= $%d", len(args))
	}
	if params.To != nil {
		args = append(args, *params.To)
		query += fmt.Sprintf(" AND s.spent_at < $%d", len(args))
	}
	args = append(args, params.Limit, params.Offset)
	query += fmt.Sprintf(" ORDER BY s.created_at DESC, s.id DESC LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list spendings: %w", err)
	}
	defer rows.Close()

	result := make([]SpendingWithUsername, 0)
	for rows.Next() {
		var item SpendingWithUsername
		err := rows.Scan(&item.ID, &item.UserID, &item.Amount, &item.Category,
			&item.SpentAt, &item.Comment, &item.CreatedAt, &item.UpdatedAt, &item.Username)
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

func (r *PostgresRepository) Update(ctx context.Context, id, userID int32, params UpdateParams) (Spending, error) {
	result, err := mapSpending(r.db.QueryRow(ctx, `
		UPDATE spendings
		SET amount = COALESCE($3::integer, amount),
		    category = COALESCE($4::category, category),
		    spent_at = COALESCE($5::timestamptz, spent_at),
		    comment = COALESCE($6::text, comment),
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND user_id = $2
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

func (r *PostgresRepository) UpdateByID(ctx context.Context, id int32, params UpdateParams) (Spending, error) {
	result, err := mapSpending(r.db.QueryRow(ctx, `
		UPDATE spendings
		SET amount = COALESCE($2::integer, amount),
		    category = COALESCE($3::category, category),
		    spent_at = COALESCE($4::timestamptz, spent_at),
		    comment = COALESCE($5::text, comment),
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, user_id, amount, category, spent_at, comment, created_at, updated_at
	`, id, params.Amount, params.Category, params.SpentAt, params.Comment))
	if errors.Is(err, pgx.ErrNoRows) {
		return Spending{}, ErrNotFound
	}
	if err != nil {
		return Spending{}, fmt.Errorf("update spending: %w", err)
	}
	return result, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id, userID int32) error {
	result, err := r.db.Exec(ctx, `DELETE FROM spendings WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("delete spending: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *PostgresRepository) DeleteByID(ctx context.Context, id int32) error {
	result, err := r.db.Exec(ctx, `DELETE FROM spendings WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete spending: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id int32) (Spending, error) {
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
