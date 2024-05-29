package entity

import (
	"context"
	"github.com/google/uuid"
)

type Source struct {
	Id      uuid.UUID
	Name    string
	Data    []byte
	Session uuid.UUID
}

type SourceRepo interface {
	GetById(ctx context.Context, id uuid.UUID) (*Source, error)
	GetBySession(ctx context.Context, sessionId uuid.UUID) ([]*Source, error)
	GetByFilename(ctx context.Context, sessionId uuid.UUID, filename string) (*Source, error)
	Update(ctx context.Context, source Source) error
	Create(ctx context.Context, source Source) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteBySession(ctx context.Context, sessionId uuid.UUID) error
}
