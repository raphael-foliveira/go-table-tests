package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func New(databaseUrl string) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", databaseUrl)
}
