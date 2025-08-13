package url

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewURLRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// CreateLink Сохраняет короткую ссылку в базу данных
func (r *Repository) CreateLink(ctx context.Context, userID uuid.UUID, originalURL, code string) (*Link, error) {
	query := `INSERT INTO links (user_id, original_url, short_code)
	          VALUES ($1, $2, $3)
			  RETURNING id, created_at`

	var link Link
	link.UserID = userID
	link.Original = originalURL
	link.ShortCode = code
	link.IsActive = true

	err := r.pool.QueryRow(ctx, query, userID, originalURL, code).
		Scan(&link.ID, &link.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &link, nil
}

// GetLink возращает оригинальную ссылку по коду
func (r *Repository) GetLink(ctx context.Context, code string) (*Link, error) {
	query := `SELECT id, user_id, original_url, short_code, created_at, is_active
	          FROM links
	          WHERE short_code = $1`

	var link Link
	err := r.pool.QueryRow(ctx, query, code).
		Scan(&link.ID, &link.UserID, &link.Original, &link.ShortCode, &link.CreatedAt, &link.IsActive)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &link, nil
}

// GetUserLinks возращает все созданные ссылки пользователя
func (r *Repository) GetUserLinks(ctx context.Context, userID uuid.UUID) (*[]Link, error) {
	query := `SELECT id, user_id, original_url, short_code, created_at, is_active
	          FROM links
	          WHERE user_id = $1`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []Link
	for rows.Next() {
		var l Link
		if err := rows.Scan(&l.ID, &l.UserID, &l.Original, &l.ShortCode, &l.CreatedAt, &l.IsActive); err != nil {
			return nil, err
		}
		links = append(links, l)
	}
	return &links, nil
}

// LogClick логирует переход по ссылке
func (r *Repository) LogClick(ctx context.Context, linkID uuid.UUID, ip, userAgent, referrer string) error {
	query := `INSERT INTO clicks (link_id, ip_address, user_agent, referrer)
	          VALUES ($1, $2, $3, $4)`
	_, err := r.pool.Exec(ctx, query, linkID, ip, userAgent, referrer)
	return err
}

// GetClicks возращает статистику по переходам
func (r *Repository) GetClicks(ctx context.Context, linkID uuid.UUID) (*[]Click, error) {
	query := `SELECT id, timestamp, ip_address, user_agent, referrer
	          FROM clicks
	          WHERE link_id = $1`

	rows, err := r.pool.Query(ctx, query, linkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clicks []Click
	for rows.Next() {
		var c Click
		if err := rows.Scan(&c.ID, &c.Timestamp, &c.IP, &c.UserAgent, &c.Referrer); err != nil {
			return nil, err
		}
		c.LinkID = linkID
		clicks = append(clicks, c)
	}
	return &clicks, nil
}
