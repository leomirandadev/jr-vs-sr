package messages

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/leomirandadev/capsulas/internal/models"
	"github.com/leomirandadev/capsulas/internal/repositories"
	"github.com/leomirandadev/capsulas/pkg/customerr"
	"github.com/leomirandadev/capsulas/pkg/mail"
)

type IService interface {
	Create(ctx context.Context, req models.CreateMessageReq) (*models.Message, error)
}

type serviceImpl struct {
	repos   *repositories.DB
	mailing mail.MailSender
}

func New(repos *repositories.DB, mailing mail.MailSender) IService {
	s := &serviceImpl{repos, mailing}

	go s.worker()

	return s
}

func (srv serviceImpl) Create(ctx context.Context, req models.CreateMessageReq) (*models.Message, error) {
	dbMessage := models.Message{
		CapsulaID: req.CapsulaID,
		Message:   req.Message,
		Email:     req.Email,
		PhotoURL:  req.PhotoURL,
		CreatedAt: time.Now(),
	}
	var err error
	dbMessage.ID, err = srv.repos.Messages.Create(ctx, dbMessage)
	if err != nil {
		return nil, customerr.WithStatus(http.StatusInternalServerError, "error to create capsula", err.Error())
	}

	return &dbMessage, nil
}

func (srv serviceImpl) worker() {
	for {
		time.Sleep(time.Hour)
		ctx := context.Background()
		capsulas, err := srv.repos.Capsula.ToSendToday(ctx)
		if err != nil {
			slog.Error("error to recover data", "err", err)
			continue
		}

		for _, capsula := range capsulas {
			messages, err := srv.repos.Messages.ListMessages(ctx, capsula.ID)
			if err != nil {
				slog.Error("error to recover data", "err", err)
				continue
			}
			for _, message := range messages {
				if err := srv.mailing.Send([]string{message.Email}, "Capsula do tempo", fmt.Sprintf("Your message is: %s", message.Message)); err != nil {
					slog.Error("user is not getting the message", "err", err)
					continue
				}
			}
			if err := srv.repos.Capsula.SetAsSent(ctx, capsula.ID); err != nil {
				slog.Error("error to set data", "err", err)
				continue
			}
		}

	}
}
