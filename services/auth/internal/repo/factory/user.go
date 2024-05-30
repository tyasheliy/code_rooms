package factory

import (
	"errors"
	"github.com/tyasheliy/code_rooms/services/auth/internal/entity"
	"github.com/tyasheliy/code_rooms/services/auth/internal/repo"
	"github.com/tyasheliy/code_rooms/services/auth/internal/storage"
)

func CreateUserRepo(driver string, s *storage.Storage) (entity.UserRepo, error) {
	switch driver {
	case storage.PG_DRIVER:
		return repo.NewPostgresUserRepo(s), nil
	default:
		return nil, errors.New("unknown storage driver")
	}
}
