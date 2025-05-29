-- +goose Up
CREATE TABLE IF NOT EXISTS capsulas(
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    open_date TIMESTAMPTZ NOT NULL,
    sent BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TRIGGER update_capsulas_updated_on BEFORE
UPDATE ON capsulas FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();
-- +goose Down
DROP TRIGGER update_capsulas_updated_on ON capsulas;
DROP TABLE IF EXISTS capsulas;