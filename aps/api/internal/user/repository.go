package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrAlreadyExists = errors.New("user already exists")
var ErrNotFound = errors.New("user not found")

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) List(ctx context.Context) ([]User, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, telegram_user_id, username, name, created_at
		FROM users
		ORDER BY username NULLS LAST, name NULLS LAST, id
	`)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.TelegramUserID, &user.Username, &user.Name, &user.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read users: %w", err)
	}
	return users, nil
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

func (r *PostgresRepository) GetByTelegramID(ctx context.Context, telegramID int64) (User, error) {
	var user User
	err := r.db.QueryRow(ctx, `
		SELECT id, telegram_user_id, username, name, created_at
		FROM users WHERE telegram_user_id = $1
	`, telegramID).Scan(&user.ID, &user.TelegramUserID, &user.Username, &user.Name, &user.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}
	if err != nil {
		return User{}, fmt.Errorf("get user by Telegram ID: %w", err)
	}
	return user, nil
}
