package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/tyasheliy/code_rooms/services/editor/internal/entity"
	"github.com/tyasheliy/code_rooms/services/editor/pkg/v1/logger"
)

type EntryUseCase struct {
	logger logger.AppLogger
	repo   entity.EntryRepo
}

func NewEntryUseCase(logger logger.AppLogger, repo entity.EntryRepo) *EntryUseCase {
	return &EntryUseCase{
		logger: logger,
		repo:   repo,
	}
}

func (s *EntryUseCase) Create(ctx context.Context, sessionId uuid.UUID) (*entity.Entry, error) {
	entry := &entity.Entry{
		Id:      uuid.New(),
		Session: sessionId,
	}

	err := s.repo.Create(ctx, entry)
	if err != nil {
		s.logger.Error(ctx, "error creating entry",
			"sessionId", sessionId,
			"error", err,
		)
		return nil, errors.New("failed to create entry")
	}
	s.logger.Debug(ctx, "created entry",
		"entry", entry,
	)

	return entry, nil
}

func (s *EntryUseCase) Check(ctx context.Context, id uuid.UUID) (bool, error) {
	_, err := s.repo.GetById(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "error checking entry",
			"id", id,
			"error", err,
		)
		return false, err
	}
	s.logger.Debug(ctx, "checked entry",
		"id", id,
	)

	return true, nil
}

func (s *EntryUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "error deleting entry",
			"id", id,
			"error", err,
		)
		return err
	}
	s.logger.Debug(ctx, "deleted entry",
		"id", id,
	)

	return nil
}
