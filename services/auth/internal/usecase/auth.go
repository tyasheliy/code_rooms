package usecase

import (
	"context"
	"errors"
	"github.com/tyasheliy/code_rooms/services/auth/internal/entity"
	"github.com/tyasheliy/code_rooms/services/auth/pkg/v1/hasher"
	"github.com/tyasheliy/code_rooms/services/auth/pkg/v1/jwtutils"
	"github.com/tyasheliy/code_rooms/services/editor/pkg/v1/logger"
)

const default_role = "user"

type AuthUseCase struct {
	logger     logger.AppLogger
	jwtBuilder *jwtutils.Builder
	hasher     hasher.Hasher
	user       entity.UserRepo
	role       entity.RoleRepo
}

func NewAuth(logger logger.AppLogger, jwtBuilder *jwtutils.Builder, hasher hasher.Hasher, user entity.UserRepo, role entity.RoleRepo) *AuthUseCase {
	return &AuthUseCase{
		logger:     logger,
		jwtBuilder: jwtBuilder,
		hasher:     hasher,
		user:       user,
		role:       role,
	}
}

func (s *AuthUseCase) Authenticate(ctx context.Context, login string, password string) (accessToken string, refreshToken string, err error) {
	user, err := s.user.GetByLogin(ctx, login)
	if err != nil {
		s.logger.Error(ctx, "error.auth.authenticate.user.get_by_login",
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

	accessToken, err = s.jwtBuilder.Claim("id", user.Id).Claim("role", user.RoleId).BuildRaw()
	if err != nil {
		s.logger.Error(ctx, "error.auth.authenticate.jwt_builder.build_raw",
			"error", err,
		)
		return "", "", errors.New("failed to build token")
	}

	return accessToken, accessToken, nil
}

func (s *AuthUseCase) Register(ctx context.Context, login string, password string) (*entity.User, error) {
	user, err := s.user.GetByLogin(ctx, login)
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

	defaultRole, err := s.role.GetByName(ctx, default_role)
	if err != nil {
		s.logger.Error(ctx, "error.auth.register.role.get_by_name",
			"role_name", default_role,
			"error", err,
		)
		return nil, errors.New("failed to register user")
	}

	user = &entity.User{
		RoleId:       defaultRole.Id,
		Login:        login,
		PasswordHash: passwordHash,
	}

	saved, err := s.user.Create(ctx, user)
	if err != nil {
		s.logger.Error(ctx, "error.auth.register.user.create",
			"login", login,
			"error", err,
		)

		return nil, errors.New("failed to register user")
	}

	return saved, nil
}
