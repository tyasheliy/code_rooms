package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/tyasheliy/code_rooms/services/editor/internal/entity"
	"github.com/tyasheliy/code_rooms/services/editor/pkg/v1/logger"
)

type SessionUseCase struct {
	repo   entity.SessionRepo
	logger logger.AppLogger
}

func NewSessionUseCase(logger logger.AppLogger, repo entity.SessionRepo) *SessionUseCase {
	return &SessionUseCase{
		repo:   repo,
		logger: logger,
	}
}

func (s *SessionUseCase) GetById(ctx context.Context, id uuid.UUID) (*entity.Session, error) {
	session, err := s.repo.GetById(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "error getting session",
			"id", id,
			"error", err,
		)
		return nil, errors.New("session not found")
	}
	s.logger.Debug(ctx, "session from repo", "session", session)

	return session, nil
}

func (s *SessionUseCase) Create(ctx context.Context) (*entity.Session, error) {
	id := uuid.New()

	session := &entity.Session{
		Id: id,
	}

	err := s.repo.Create(ctx, session)
	if err != nil {
		s.logger.Error(ctx, "error creating session", "error", err)
		return nil, errors.New("failed to create session")
	}
	s.logger.Debug(ctx, "created session", "session", session)

	return session, nil
}

func (s *SessionUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "error deleting session",
			"id", id,
			"error", err,
		)
		return errors.New("failed to delete session")
	}
	s.logger.Debug(ctx, "deleted session", "id", id)

	return nil
}
