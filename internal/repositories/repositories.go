package repositories

import (
	"github.com/jmoiron/sqlx"
	"github.com/leomirandadev/capsulas/internal/repositories/capsula"
	"github.com/leomirandadev/capsulas/internal/repositories/messages"
)

type DB struct {
	Capsula  capsula.IRepository
	Messages messages.IRepository
}

type Options struct {
	WriterSqlx *sqlx.DB
	ReaderSqlx *sqlx.DB
}

func New(opts Options) *DB {
	return &DB{
		Capsula:  capsula.NewSqlx(opts.WriterSqlx, opts.ReaderSqlx),
		Messages: messages.NewSqlx(opts.WriterSqlx, opts.ReaderSqlx),
	}
}
