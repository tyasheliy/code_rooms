package entity

import (
	"context"
	"github.com/google/uuid"
)

type Entry struct {
	Id      uuid.UUID `json:"id"`
	Session uuid.UUID `json:"session"`
}

type EntryRepo interface {
	GetById(ctx context.Context, id uuid.UUID) (*Entry, error)
	Create(ctx context.Context, entry *Entry) error
	Delete(ctx context.Context, id uuid.UUID) error
}
