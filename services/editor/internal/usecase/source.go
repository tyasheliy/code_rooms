package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/tyasheliy/code_rooms/services/editor/internal/entity"
	"github.com/tyasheliy/code_rooms/services/editor/pkg/v1/logger"
)

type SourceUseCase struct {
	logger logger.AppLogger
	repo   entity.SourceRepo
}

func NewSourceUseCase(logger logger.AppLogger, repo entity.SourceRepo) *SourceUseCase {
	return &SourceUseCase{
		logger: logger,
		repo:   repo,
	}
}

func (s *SourceUseCase) Create(ctx context.Context, sessionId uuid.UUID, filename string, data []byte) (*entity.Source, error) {
	id := uuid.New()

	source := entity.Source{
		Id:      id,
		Name:    filename,
		Data:    data,
		Session: sessionId,
	}

	err := s.repo.Create(ctx, source)
	if err != nil {
		s.logger.Error(ctx, "error creating source",
			"id", id,
			"sessionId", sessionId,
			"filename", filename,
			"error", err,
		)
		return nil, err
	}
	s.logger.Debug(ctx, "created source",
		"id", id,
		"sessionId", sessionId,
		"filename", filename,
	)

	return &source, nil
}

func (s *SourceUseCase) UpdateData(ctx context.Context, id uuid.UUID, data []byte) error {
	source, err := s.repo.GetById(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "error getting source",
			"id", id,
			"error", err,
		)
		return errors.New("source not found")
	}
	s.logger.Debug(ctx, "got source",
		"id", id,
	)

	source.Data = data

	err = s.repo.Update(ctx, *source)
	if err != nil {
		s.logger.Error(ctx, "error updating source",
			"id", id,
			"error", err,
		)
		return errors.New("failed to update source")
	}
	s.logger.Debug(ctx, "updated source",
		"id", id,
	)

	return nil
}

func (s *SourceUseCase) GetBySession(ctx context.Context, sessionId uuid.UUID) ([]*entity.Source, error) {
	sources, err := s.repo.GetBySession(ctx, sessionId)
	if err != nil {
		s.logger.Error(ctx, "error getting sources by session",
			"sessionId", sessionId,
			"error", err,
		)
		return nil, err
	}
	s.logger.Debug(ctx, "got sources by session",
		"sessionId", sessionId,
	)

	return sources, nil
}

func (s *SourceUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "error deleting source",
			"id", id,
			"error", err,
		)
		return errors.New("failed to delete source")
	}
	s.logger.Debug(ctx, "deleted source",
		"id", id,
	)

	return nil
}
