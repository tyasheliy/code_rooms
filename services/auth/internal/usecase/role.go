package usecase

import (
	"context"
	"errors"
	"github.com/tyasheliy/code_rooms/services/auth/internal/entity"
	"github.com/tyasheliy/code_rooms/services/editor/pkg/v1/cache"
	"github.com/tyasheliy/code_rooms/services/editor/pkg/v1/logger"
	"strconv"
)

type RoleUseCase struct {
	logger logger.AppLogger
	cache  cache.AppCache
	repo   entity.RoleRepo
}

func NewRole(logger logger.AppLogger, cache cache.AppCache, repo entity.RoleRepo) *RoleUseCase {
	return &RoleUseCase{logger: logger, cache: cache, repo: repo}
}

func (s *RoleUseCase) GetById(ctx context.Context, id int) (*entity.Role, error) {
	cacheKey := s.cache.BuildCacheKey("role", strconv.Itoa(id))

	var role *entity.Role

	cached, exists := s.cache.Get(ctx, cacheKey)
	if exists {
		cachedRole, ok := cached.(entity.Role)
		if !ok {
			s.logger.Warn(ctx, "invalid role cache")

			_ = s.cache.Delete(ctx, cacheKey)
		} else {
			role = &cachedRole
		}
	} else {
		var err error
		role, err = s.repo.GetById(ctx, id)
		if err != nil {
			s.logger.Error(ctx, "error.role.get_by_id.repo.get_by_id",
				"id", id,
				"error", err,
			)
			return nil, errors.New("role not found")
		}

		err = s.cache.Set(ctx, cacheKey, *role)
		if err != nil {
			s.logger.Warn(ctx, "warn.role.get_by_id.cache.set",
				"cache_key", cacheKey,
				"role", *role,
				"error", err,
			)
		}
	}

	return role, nil
}
