package auth

import (
	"github.com/leomirandadev/capsulas/pkg/httprouter"
)

type AuthMiddleware interface {
	BasicAuth(next httprouter.HandlerFunc) httprouter.HandlerFunc
}
