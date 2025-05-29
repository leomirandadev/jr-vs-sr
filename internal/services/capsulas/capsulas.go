package capsulas

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/leomirandadev/capsulas/internal/models"
	"github.com/leomirandadev/capsulas/internal/repositories"
	"github.com/leomirandadev/capsulas/pkg/customerr"
	"github.com/leomirandadev/capsulas/pkg/mail"
)

type IService interface {
	Create(ctx context.Context, req models.CreateCapsulaReq) (*models.Capsula, error)
	GetByID(ctx context.Context, ID string) (*models.Capsula, error)
}

type serviceImpl struct {
	repos   *repositories.DB
	mailing mail.MailSender
}

func New(repos *repositories.DB, mailing mail.MailSender) IService {
	return &serviceImpl{repos, mailing}
}

func (srv serviceImpl) Create(ctx context.Context, req models.CreateCapsulaReq) (*models.Capsula, error) {
	openDate, _ := time.Parse(time.DateOnly, req.OpenDate)

	dbCapsula := models.Capsula{
		Name:      req.Name,
		OpenDate:  openDate,
		CreatedAt: time.Now(),
	}
	var err error
	dbCapsula.ID, err = srv.repos.Capsula.Create(ctx, dbCapsula)
	if err != nil {
		return nil, customerr.WithStatus(http.StatusInternalServerError, "error to create capsula", err.Error())
	}

	return &dbCapsula, nil
}

func (srv serviceImpl) GetByID(ctx context.Context, ID string) (*models.Capsula, error) {
	result, err := srv.repos.Capsula.GetByID(ctx, ID)
	if err != nil {
		slog.WarnContext(ctx, "repositories.user.get_by_id", "err", err)
		return nil, customerr.WithStatus(http.StatusNotFound, "user not found", err)
	}

	return result, nil
}
