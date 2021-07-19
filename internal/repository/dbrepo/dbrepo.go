package dbrepo

import (
	"database/sql"

	"github.com/ArmanurRahman/booking/internal/config"
	"github.com/ArmanurRahman/booking/internal/repository"
)

type postgressDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgressDBRepo{
		App: a,
		DB:  conn,
	}
}
