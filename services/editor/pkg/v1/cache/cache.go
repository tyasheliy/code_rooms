package cache

import (
	"context"
	"strings"
)

type AppCache interface {
	Get(ctx context.Context, key string) (interface{}, bool)
	Set(ctx context.Context, key string, value interface{}) error
	Delete(ctx context.Context, key string) error
	BuildCacheKey(args ...string) string
}

type KeyBuilder struct {
	KeyBuilderOptions
	pre    []string
	values []string
	post   []string
}

type KeyBuilderOptions struct {
	Separator string
}

func NewKeyBuilder(options KeyBuilderOptions) *KeyBuilder {
	return &KeyBuilder{
		KeyBuilderOptions: options,
		pre:               make([]string, 0),
		post:              make([]string, 0),
	}
}

func (b *KeyBuilder) Pre(value string) *KeyBuilder {
	b.pre = append(b.pre, value)
	return b
}

func (b *KeyBuilder) Post(value string) *KeyBuilder {
	b.post = append(b.post, value)
	return b
}

func (b *KeyBuilder) BuildCacheKey(args ...string) string {
	elems := make([]string, 0, len(b.pre)+len(args)+len(b.post))

	elems = append(elems, b.pre...)
	elems = append(elems, args...)
	elems = append(elems, b.post...)

	return strings.Join(elems, b.Separator)
}
