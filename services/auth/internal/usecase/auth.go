package usecase

import (
	"context"
	"errors"
	"github.com/tyasheliy/code_rooms/services/auth/internal/entity"
	"github.com/tyasheliy/code_rooms/services/auth/pkg/v1/hasher"
	"github.com/tyasheliy/code_rooms/services/auth/pkg/v1/jwtutils"
	"github.com/tyasheliy/code_rooms/services/editor/pkg/v1/logger"
)

type AuthUseCase struct {
	logger     logger.AppLogger
	jwtBuilder *jwtutils.Builder
	hasher     hasher.Hasher
	repo       entity.UserRepo
}

func NewAuth(logger logger.AppLogger, jwtBuilder *jwtutils.Builder, hasher hasher.Hasher, repo entity.UserRepo) *AuthUseCase {
	return &AuthUseCase{
		logger:     logger,
		jwtBuilder: jwtBuilder,
		hasher:     hasher,
		repo:       repo,
	}
}

func (s *AuthUseCase) Authenticate(ctx context.Context, login string, password string) (accessToken string, refreshToken string, err error) {
	user, err := s.repo.GetByLogin(ctx, login)
	if err != nil {
		s.logger.Error(ctx, "error.auth.authenticate.repo.get_by_login",
			"login", login,
			"error", err,
		)
		return "", "", errors.New("user not found")
	}

	if !s.hasher.Check(password, user.PasswordHash) {
		err := errors.New("invalid password")

		s.logger.Warn(ctx, "warn.auth.authenticate.hasher.check",
			"error", err,
		)
		return "", "", err
	}

	accessToken, err = s.jwtBuilder.Claim("id", user.Id).BuildRaw()
	if err != nil {
		s.logger.Error(ctx, "error.auth.authenticate.jwt_builder.build_raw",
			"error", err,
		)
		return "", "", errors.New("failed to build token")
	}

	return accessToken, accessToken, nil
}

func (s *AuthUseCase) Register(ctx context.Context, login string, password string) (*entity.User, error) {
	user, err := s.repo.GetByLogin(ctx, login)
	if err == nil {
		err := errors.New("user already exists")

		s.logger.Error(ctx, "error.auth.register",
			"login", login,
			"error", err,
		)

		return nil, err
	}

	passwordHash, err := s.hasher.Hash(password)
	if err != nil {
		s.logger.Error(ctx, "error.auth.register.hasher.hash",
			"error", err,
		)
		return nil, errors.New("failed to register user")
	}

	user = &entity.User{
		Login:        login,
		PasswordHash: passwordHash,
	}

	saved, err := s.repo.Create(ctx, user)
	if err != nil {
		s.logger.Error(ctx, "error.auth.register.repo.create",
			"login", login,
			"error", err,
		)

		return nil, errors.New("failed to register user")
	}

	return saved, nil
}
