package middlewares

import (
	"github.com/leomirandadev/capsulas/internal/handlers/middlewares/auth"
)

type Middleware struct {
	Auth auth.AuthMiddleware
}

func New(basicAuthUser, basicAuthPassword string) *Middleware {
	return &Middleware{
		Auth: auth.New(basicAuthUser, basicAuthPassword),
	}
}
