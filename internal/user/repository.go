package user

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewAdRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Create - сохраняет пользователя в БД
func (r *Repository) Create(ctx context.Context, u *User) (*User, error) {
	query := `
		INSERT INTO users (email, password)
		VALUES ($1, $2)
		RETURNING id, created_at
	`
	err := r.pool.QueryRow(ctx, query, u.Email, u.PasswordHash).Scan(&u.ID, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// GetByEmail - возвращает пользователя по email (или nil, если не найден)
func (r *Repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, password, created_at
		FROM users
		WHERE email = $1
	`
	row := r.pool.QueryRow(ctx, query, strings.ToLower(email))

	var u User
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
