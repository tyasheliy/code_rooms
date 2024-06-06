package usecase

import (
	"context"
	"errors"
	"github.com/tyasheliy/code_rooms/services/auth/internal/entity"
	"github.com/tyasheliy/code_rooms/services/editor/pkg/v1/cache"
	"github.com/tyasheliy/code_rooms/services/editor/pkg/v1/logger"
	"strconv"
)

type UserUseCase struct {
	logger logger.AppLogger
	cache  cache.AppCache
	repo   entity.UserRepo
}

func NewUser(logger logger.AppLogger, cache cache.AppCache, repo entity.UserRepo) *UserUseCase {
	return &UserUseCase{
		logger: logger,
		cache:  cache,
		repo:   repo,
	}
}

func (s *UserUseCase) GetById(ctx context.Context, id int64) (*entity.User, error) {
	cacheKey := s.cache.BuildCacheKey("user", strconv.Itoa(int(id)))

	cached, exists := s.cache.Get(ctx, cacheKey)
	if exists {
		cachedUser, ok := cached.(entity.User)
		if ok {
			return &cachedUser, nil
		}

		s.logger.Warn(ctx, "warn.user.get_by_id.cache.get",
			"error", errors.New("invalid user cache"),
		)

		_ = s.cache.Delete(ctx, cacheKey)
	}

	user, err := s.repo.GetById(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "error.user.get_by_id.repo.get_by_id",
			"id", id,
			"error", err,
		)
		return nil, errors.New("user not found")
	}

	err = s.cache.Set(ctx, cacheKey, *user)
	if err != nil {
		s.logger.Warn(ctx, "warn.user.get_by_id.cache.set",
			"cache_key", cacheKey,
			"error", err,
		)
	}

	return user, nil
}
