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
CREATE TRIGGER update_messages_updated_on BEFORE
UPDATE ON messages FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();
-- +goose Down
DROP TRIGGER update_messages_updated_on ON messages;
DROP TABLE IF EXISTS messages;