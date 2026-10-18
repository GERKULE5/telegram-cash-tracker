package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrAlreadyExists = errors.New("user already exists")

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, params CreateUserParams) (User, error) {
	var user User
	err := r.db.QueryRow(ctx, `
		INSERT INTO users (telegram_user_id, username, name)
		VALUES ($1, $2, $3)
		RETURNING id, telegram_user_id, username, name, created_at
	`, params.TelegramUserID, params.Username, params.Name).Scan(
		&user.ID,
		&user.TelegramUserID,
		&user.Username,
		&user.Name,
		&user.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return User{}, ErrAlreadyExists
		}
		return User{}, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}
