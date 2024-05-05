package cache

import (
	"context"
	"github.com/patrickmn/go-cache"
	"time"
)

type InMemoryAppCache struct {
	cache *cache.Cache
	*KeyBuilder
}

func NewInMemoryAppCache(keyBuilder *KeyBuilder, expiration time.Duration, cleanup time.Duration) *InMemoryAppCache {
	c := cache.New(expiration, cleanup)

	return &InMemoryAppCache{
		KeyBuilder: keyBuilder,
		cache:      c,
	}
}

func (c *InMemoryAppCache) Get(ctx context.Context, key string) (interface{}, bool) {
	return c.cache.Get(key)
}

func (c *InMemoryAppCache) Set(ctx context.Context, key string, value interface{}) error {
	c.cache.Set(key, value, cache.DefaultExpiration)
	return nil
}

func (c *InMemoryAppCache) Delete(ctx context.Context, key string) error {
	c.cache.Delete(key)
	return nil
}
