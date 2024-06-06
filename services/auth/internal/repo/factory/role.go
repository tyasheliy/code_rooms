package factory

import (
	"errors"
	"github.com/tyasheliy/code_rooms/services/auth/internal/entity"
	"github.com/tyasheliy/code_rooms/services/auth/internal/repo"
	"github.com/tyasheliy/code_rooms/services/auth/internal/storage"
)

func CreateRoleRepo(driver string, s *storage.Storage) (entity.RoleRepo, error) {
	switch driver {
	case storage.PG_DRIVER:
		return repo.NewPostgresRoleRepo(s), nil
	default:
		return nil, errors.New("unknown storage driver")
	}
}
