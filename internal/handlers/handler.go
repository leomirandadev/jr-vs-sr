package handlers

import (
	"github.com/leomirandadev/capsulas/internal/handlers/capsulas"
	"github.com/leomirandadev/capsulas/internal/handlers/middlewares"
	"github.com/leomirandadev/capsulas/internal/handlers/swagger"
	"github.com/leomirandadev/capsulas/internal/services"
	"github.com/leomirandadev/capsulas/pkg/httprouter"
)

type Options struct {
	Router            httprouter.Router
	Srv               *services.Container
	BasicAuthUser     string
	BasicAuthPassword string
}

func New(opts Options) {
	mid := middlewares.New(opts.BasicAuthUser, opts.BasicAuthPassword)

	capsulas.Init(opts.Router, opts.Srv)
	swagger.Init(mid, opts.Router)
}
