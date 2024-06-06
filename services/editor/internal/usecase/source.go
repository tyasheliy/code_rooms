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

func (s *SourceUseCase) UpdateDataByFilename(ctx context.Context, sessionId uuid.UUID, filename string, data []byte) error {
	source, err := s.repo.GetByFilename(ctx, sessionId, filename)
	if err != nil {
		id := uuid.New()

		sourceVal := entity.Source{
			Id:      id,
			Name:    filename,
			Data:    data,
			Session: sessionId,
		}

		err = s.repo.Create(ctx, sourceVal)

		if err != nil {
			return err
		}

		source = &sourceVal

		s.logger.Debug(ctx, "created source",
			"id", id.String(),
			"name", filename,
			"session_id", sessionId.String(),
		)
	}

	source.Data = data

	err = s.repo.Update(ctx, *source)
	if err != nil {
		s.logger.Error(ctx, "error updating source",
			"filename", filename,
			"error", err,
		)
		return errors.New("failed to update source")
	}
	s.logger.Debug(ctx, "updated source",
		"filename", filename,
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

func (s *SourceUseCase) DeleteBySession(ctx context.Context, sessionId uuid.UUID) error {
	err := s.repo.DeleteBySession(ctx, sessionId)
	if err != nil {
		s.logger.Error(ctx, "error deleting source by session",
			"sessionId", sessionId,
			"error", err,
		)
		return errors.New("failed to delete sources")
	}
	s.logger.Debug(ctx, "deleted sources by session", "sessionId", sessionId)

	return nil
}

func (s *SourceUseCase) GetByFilename(ctx context.Context, sessionId uuid.UUID, filename string) (*entity.Source, error) {
	source, err := s.repo.GetByFilename(ctx, sessionId, filename)
	if err != nil {
		s.logger.Error(ctx, "error getting source by filename",
			"session", sessionId,
			"filename", filename,
			"error", err,
		)
		return nil, errors.New("failed to get source")
	}
	s.logger.Debug(ctx, "got source by filename",
		"sessionId", sessionId,
		"filename", filename,
		"source", source,
	)

	return source, nil
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
