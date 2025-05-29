-- +goose Up
CREATE TABLE IF NOT EXISTS messages(
    id VARCHAR(36) PRIMARY KEY,
    message VARCHAR(200) NOT NULL,
    photo_url VARCHAR(200) NOT NULL,
    email VARCHAR(200) NOT NULL,
    capsula_id VARCHAR(36),
    FOREIGN KEY (capsula_id) REFERENCES capsulas(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose Down
DROP TABLE IF EXISTS messages;