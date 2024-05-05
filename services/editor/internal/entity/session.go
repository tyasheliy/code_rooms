package entity

import (
	"context"
	"github.com/google/uuid"
)

type Session struct {
	Id uuid.UUID
}

type SessionRepo interface {
	GetById(ctx context.Context, id uuid.UUID) (*Session, error)
	Create(ctx context.Context, session *Session) error
	Delete(ctx context.Context, id uuid.UUID) error
}
