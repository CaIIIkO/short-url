-- +goose Up
-- +goose StatementBegin
CREATE TABLE links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INTEGER NOT NULL REFERENCES users(id),
    original_url TEXT NOT NULL,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    is_active BOOLEAN DEFAULT true
);

CREATE TABLE clicks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    link_id INTEGER NOT NULL REFERENCES links(id),
    timestamp TIMESTAMP DEFAULT now(),
    ip_address TEXT,
    user_agent TEXT,
    referrer TEXT
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
