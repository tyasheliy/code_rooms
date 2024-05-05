package repo

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/tyasheliy/code_rooms/services/editor/internal/entity"
	"github.com/tyasheliy/code_rooms/services/editor/pkg/v1/cache"
)

type CacheSessionRepo struct {
	cache cache.AppCache
}

func NewCacheSessionRepo(cache cache.AppCache) *CacheSessionRepo {
	return &CacheSessionRepo{cache: cache}
}

func (r *CacheSessionRepo) GetById(ctx context.Context, id uuid.UUID) (*entity.Session, error) {
	key := r.cache.BuildCacheKey("session", id.String())

	res, ok := r.cache.Get(ctx, key)
	if !ok {
		return nil, errors.New("session with given key not found")
	}

	session, ok := res.(*entity.Session)
	if !ok {
		return nil, errors.New("invalid cached entity")
	}

	return session, nil
}

func (r *CacheSessionRepo) Create(ctx context.Context, session *entity.Session) error {
	key := r.cache.BuildCacheKey("session", session.Id.String())

	return r.cache.Set(ctx, key, session)
}

func (r *CacheSessionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	key := r.cache.BuildCacheKey("session", id.String())

	return r.cache.Delete(ctx, key)
}
