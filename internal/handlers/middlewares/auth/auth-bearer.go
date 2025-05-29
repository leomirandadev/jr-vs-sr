package auth

import (
	"net/http"

	"github.com/leomirandadev/capsulas/pkg/customerr"
	"github.com/leomirandadev/capsulas/pkg/httprouter"
)

type middlewareJWT struct {
	basicAuthUser     string
	basicAuthPassword string
}

func New(basicAuthUser, basicAuthPassword string) AuthMiddleware {
	return &middlewareJWT{
		basicAuthUser:     basicAuthUser,
		basicAuthPassword: basicAuthPassword,
	}
}

func (m *middlewareJWT) BasicAuth(next httprouter.HandlerFunc) httprouter.HandlerFunc {
	return func(httpCtx httprouter.Context) error {
		username, password, ok := httpCtx.GetRequestReader().BasicAuth()
		if ok {
			if username == m.basicAuthUser && password == m.basicAuthPassword {
				return next(httpCtx)
			}
		}

		httpCtx.GetResponseWriter().Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		return customerr.WithStatus(http.StatusUnauthorized, "unauthorizated", "")
	}
}
