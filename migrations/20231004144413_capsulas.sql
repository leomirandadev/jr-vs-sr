-- +goose Up
CREATE TABLE IF NOT EXISTS capsulas(
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    open_date TIMESTAMPTZ NOT NULL,
    sent BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose Down
DROP TABLE IF EXISTS capsulas;