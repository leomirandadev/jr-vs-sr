package main

import (
	"log"
	"log/slog"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/leomirandadev/capsulas/internal/handlers"
	"github.com/leomirandadev/capsulas/internal/repositories"
	"github.com/leomirandadev/capsulas/internal/services"
	"github.com/leomirandadev/capsulas/pkg/envs"
	"github.com/leomirandadev/capsulas/pkg/httprouter"
	"github.com/leomirandadev/capsulas/pkg/mail"
	"github.com/leomirandadev/capsulas/pkg/slogtint"
	"github.com/leomirandadev/capsulas/pkg/validator"
)

type Config struct {
	Port    string       `mapstructure:"PORT" validate:"required"`
	Env     string       `mapstructure:"ENV" validate:"required"`
	Mailing mail.Options `mapstructure:"MAILING" validate:"required"`

	Database struct {
		Reader string `mapstructure:"READER" validate:"required"`
		Writer string `mapstructure:"WRITER" validate:"required"`
	} `mapstructure:"DATABASE" validate:"required"`

	BasicAuth struct {
		User     string `mapstructure:"USER"`
		Password string `mapstructure:"PASSWORD"`
	} `mapstructure:"BASIC_AUTH"`
}

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" before paste the token
func main() {
	var cfg Config
	if err := envs.Load(".", &cfg); err != nil {
		log.Fatal("cfg not loaded", err)
	}

	if err := validator.Validate(cfg); err != nil {
		log.Fatal("missing required fields", err)
	}

	tools := toolsInit(cfg)

	writeDBConn := sqlx.MustConnect("pgx", cfg.Database.Reader)
	readDBConn := writeDBConn

	repo := repositories.New(repositories.Options{
		ReaderSqlx: readDBConn,
		WriterSqlx: writeDBConn,
	})

	srv := services.New(services.Options{
		Repo:    repo,
		Mailing: tools.mailing,
	})

	handlers.New(handlers.Options{
		Srv:               srv,
		Router:            tools.router,
		BasicAuthUser:     cfg.BasicAuth.User,
		BasicAuthPassword: cfg.BasicAuth.Password,
	})

	tools.router.Serve(cfg.Port)
}

type tools struct {
	router  httprouter.Router
	mailing mail.MailSender
}

func toolsInit(cfg Config) tools {

	slog.SetDefault(slog.New(
		slogtint.NewHandler(os.Stderr, &slogtint.Options{
			AddSource: true,
			Level:     slog.LevelDebug,
		}),
	))

	return tools{
		router:  httprouter.NewChiRouter(),
		mailing: mail.NewSMTP(cfg.Mailing),
	}
}
