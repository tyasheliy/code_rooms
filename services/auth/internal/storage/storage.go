package storage

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/tyasheliy/code_rooms/services/auth/internal/config"
)

const (
	PG_DRIVER = "postgres"
)

type Storage struct {
	*sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{
		DB: db,
	}
}

func BuildDbConn(cfg *config.StorageConfig) (*sql.DB, error) {
	var db *sql.DB

	switch cfg.Driver {
	case PG_DRIVER:
		url := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Database,
		)
		openedDb, err := sql.Open("postgres", url)
		if err != nil {
			return nil, err
		}

		db = openedDb
	default:
		return nil, errors.New("unknown database driver")
	}

	return db, nil
}
