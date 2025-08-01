package url

import (
	"time"

	"github.com/google/uuid"
)

type Link struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Original  string    `json:"original_url"`
	ShortCode string    `json:"short_code"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"`
}

type Click struct {
	ID        uuid.UUID `json:"id"`
	LinkID    uuid.UUID `json:"link_id"`
	Timestamp time.Time `json:"timestamp"`
	IP        string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Referrer  string    `json:"referrer"`
}
