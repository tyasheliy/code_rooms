package entity

import (
	"context"
	"github.com/google/uuid"
)

type Entry struct {
	Id      uuid.UUID
	Session uuid.UUID
}

type EntryRepo interface {
	GetById(ctx context.Context, id uuid.UUID) (*Entry, error)
	Create(ctx context.Context, entry *Entry) error
	Delete(ctx context.Context, id uuid.UUID) error
}
