package messages

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/leomirandadev/capsulas/internal/models"
)

type IRepository interface {
	Create(ctx context.Context, newMessage models.Message) (string, error)
	ListMessages(ctx context.Context, capsulaID string) ([]models.Message, error)
}

func NewSqlx(writer, reader *sqlx.DB) IRepository {
	return &repoSqlx{writer: writer, reader: reader}
}

type repoSqlx struct {
	writer *sqlx.DB
	reader *sqlx.DB
}

func (repo repoSqlx) Create(ctx context.Context, newMessage models.Message) (string, error) {
	newMessage.ID = uuid.New().String()

	_, err := repo.writer.NamedExecContext(ctx, `
		INSERT INTO messages (id, message, email, capsula_id, photo_url) 
		VALUES (:id, :message, :email, :capsula_id, :photo_url) 
	`, newMessage)

	if err != nil {
		slog.WarnContext(ctx, "writer.exec_context", "err", err)
		return "", err
	}

	return newMessage.ID, nil
}

func (repo repoSqlx) ListMessages(ctx context.Context, capsulaID string) ([]models.Message, error) {
	var capsula []models.Message

	err := repo.reader.SelectContext(ctx, &capsula, `
		SELECT * FROM messages WHERE capsula_id = $1
	`, capsulaID)

	if err != nil {
		slog.WarnContext(ctx, "reader.get_context", "err", err)
		return nil, err
	}

	return capsula, nil
}
