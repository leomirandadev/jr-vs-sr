package swagger

import (
	docs "github.com/leomirandadev/capsulas/docs"
	"github.com/leomirandadev/capsulas/internal/handlers/middlewares"
	"github.com/leomirandadev/capsulas/pkg/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Init(mid *middlewares.Middleware, router httprouter.Router) {

	docs.SwaggerInfo.Title = "Swagger"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	router.GET("/swagger/*", mid.Auth.BasicAuth(
		router.ParseHandler(
			httpSwagger.Handler(
				httpSwagger.URL("doc.json"),
			),
		),
	))
}
