package capsula

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/leomirandadev/capsulas/internal/models"
)

type IRepository interface {
	Create(ctx context.Context, newCapsula models.Capsula) (string, error)
	GetByID(ctx context.Context, ID string) (*models.Capsula, error)
	ToSendToday(ctx context.Context) ([]models.Capsula, error)
	SetAsSent(ctx context.Context, id string) error
}

func NewSqlx(writer, reader *sqlx.DB) IRepository {
	return &repoSqlx{writer: writer, reader: reader}
}

type repoSqlx struct {
	writer *sqlx.DB
	reader *sqlx.DB
}

func (repo repoSqlx) Create(ctx context.Context, newCapsula models.Capsula) (string, error) {
	newCapsula.ID = uuid.New().String()

	_, err := repo.writer.NamedExecContext(ctx, `
		INSERT INTO capsulas (id, name, open_date, sent) 
		VALUES (:id, :name, :open_date, FALSE) 
	`, newCapsula)

	if err != nil {
		slog.WarnContext(ctx, "writer.exec_context", "err", err)
		return "", err
	}

	return newCapsula.ID, nil
}

func (repo repoSqlx) GetByID(ctx context.Context, ID string) (*models.Capsula, error) {
	var capsula models.Capsula

	err := repo.reader.GetContext(ctx, &capsula, `
		SELECT * FROM capsulas WHERE id = $1
	`, ID)

	if err != nil {
		slog.WarnContext(ctx, "reader.get_context", "err", err)
		return nil, err
	}

	return &capsula, nil
}

func (repo repoSqlx) ToSendToday(ctx context.Context) ([]models.Capsula, error) {
	var capsula []models.Capsula

	if err := repo.reader.SelectContext(ctx, &capsula, `
		SELECT * FROM capsulas WHERE open_date > NOW() AND sent = FALSE
	`); err != nil {
		slog.WarnContext(ctx, "reader.get_context", "err", err)
		return nil, err
	}

	return capsula, nil
}

func (repo repoSqlx) SetAsSent(ctx context.Context, id string) error {
	if _, err := repo.writer.ExecContext(ctx, `
		UPDATE capsulas SET sent = TRUE WHERE id = $1
	`, id); err != nil {
		slog.WarnContext(ctx, "reader.get_context", "err", err)
		return err
	}

	return nil
}
