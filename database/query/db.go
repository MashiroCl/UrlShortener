package database

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/mashirocl/urlshortener/config"
)

func NewDB(cfg config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open(cfg.Driver, cfg.DSN())
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, err
}
