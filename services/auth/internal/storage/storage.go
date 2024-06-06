package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/tyasheliy/code_rooms/services/auth/internal/config"
)

const (
	PG_DRIVER = "postgres"
)

type Storage struct {
	*sql.DB
	migration *migrate.Migrate
}

type StorageOptions struct {
	db        *sql.DB
	migration *migrate.Migrate
}

func New(opts *StorageOptions) *Storage {
	return &Storage{
		DB:        opts.db,
		migration: opts.migration,
	}
}

func (s *Storage) Migrate() error {
	err := s.migration.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}

		return err
	}

	return nil
}

func (s *Storage) Down() error {
	err := s.migration.Down()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}

		return err
	}

	return nil
}

func BuildStorageOpts(cfg *config.StorageConfig) (*StorageOptions, func(), error) {
	var db *sql.DB
	var driver database.Driver

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
			return nil, nil, err
		}

		openedDriver, err := postgres.WithInstance(openedDb, &postgres.Config{})
		if err != nil {
			return nil, nil, err
		}

		driver = openedDriver
		db = openedDb
	default:
		return nil, nil, errors.New("unknown database driver")
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", cfg.MigrationDir),
		cfg.Database,
		driver,
	)
	if err != nil {
		return nil, nil, err
	}

	opts := &StorageOptions{
		db:        db,
		migration: m,
	}

	return opts, func() {
		opts.db.Close()
		opts.migration.Close()
	}, nil
}
