package repo

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/tyasheliy/code_rooms/services/editor/internal/entity"
	"github.com/tyasheliy/code_rooms/services/editor/pkg/v1/cache"
)

type CacheEntryRepo struct {
	cache cache.AppCache
}

func (r *CacheEntryRepo) GetById(ctx context.Context, id uuid.UUID) (*entity.Entry, error) {
	key := r.cache.BuildCacheKey("entry", id.String())

	res, ok := r.cache.Get(ctx, key)
	if !ok {
		return nil, errors.New("entry with given cache key not found")
	}

	entry, ok := res.(*entity.Entry)
	if !ok {
		return nil, errors.New("invalid cached entity")
	}

	return entry, nil
}

func (r *CacheEntryRepo) Create(ctx context.Context, entry *entity.Entry) error {
	key := r.cache.BuildCacheKey("entry", entry.Id.String())

	return r.cache.Set(ctx, key, entry)
}

func (r *CacheEntryRepo) Delete(ctx context.Context, id uuid.UUID) error {
	key := r.cache.BuildCacheKey("entry", id.String())

	return r.cache.Delete(ctx, key)
}
