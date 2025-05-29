package services

import (
	"github.com/leomirandadev/capsulas/internal/repositories"
	"github.com/leomirandadev/capsulas/internal/services/capsulas"
	"github.com/leomirandadev/capsulas/internal/services/messages"
	"github.com/leomirandadev/capsulas/pkg/mail"
)

type Container struct {
	Capsula  capsulas.IService
	Messages messages.IService
}

type Options struct {
	Repo    *repositories.DB
	Mailing mail.MailSender
}

func New(opts Options) *Container {
	return &Container{
		Capsula:  capsulas.New(opts.Repo, opts.Mailing),
		Messages: messages.New(opts.Repo, opts.Mailing),
	}
}
