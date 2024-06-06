package main

import (
	"github.com/golang-jwt/jwt"
	"github.com/tyasheliy/code_rooms/services/auth/internal/config"
	"github.com/tyasheliy/code_rooms/services/auth/internal/repo/factory"
	"github.com/tyasheliy/code_rooms/services/auth/internal/storage"
	"github.com/tyasheliy/code_rooms/services/auth/internal/usecase"
	"github.com/tyasheliy/code_rooms/services/auth/internal/webapi"
	"github.com/tyasheliy/code_rooms/services/auth/internal/webapi/handler"
	"github.com/tyasheliy/code_rooms/services/auth/pkg/v1/hasher"
	"github.com/tyasheliy/code_rooms/services/auth/pkg/v1/jwtutils"
	"github.com/tyasheliy/code_rooms/services/editor/pkg/v1/cache"
	"github.com/tyasheliy/code_rooms/services/editor/pkg/v1/logger"
	"log"
	"os"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	l := logger.NewSlogAppLogger(logger.AppLoggerOptions{
		Format: cfg.Format,
		Level:  cfg.Level,
		Output: os.Stdout,
	})

	opts, closeConn, err := storage.BuildStorageOpts(&cfg.StorageConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer closeConn()

	s := storage.New(opts)

	err = s.Migrate()
	if err != nil {
		log.Fatal(err)
	}

	keyBuilder := cache.NewKeyBuilder(cache.KeyBuilderOptions{Separator: "_"})

	c := cache.NewInMemoryAppCache(keyBuilder, cfg.Expiration, cfg.Cleanup)

	h := &hasher.BcryptHasher{}

	jwtBuilder := jwtutils.NewBuilder(cfg.Secret, cfg.TokenConfig.Expiration, jwt.SigningMethodHS256)

	userRepo, err := factory.CreateUserRepo(cfg.Driver, s)
	if err != nil {
		log.Fatal(err)
	}

	roleRepo, err := factory.CreateRoleRepo(cfg.Driver, s)
	if err != nil {
		log.Fatal(err)
	}

	auth := usecase.NewAuth(l, jwtBuilder, h, userRepo, roleRepo)
	user := usecase.NewUser(l, c, userRepo)

	app := webapi.NewWebApiApp(&cfg.AppConfig, &webapi.Handlers{
		User: handler.NewUser(user),
		Auth: handler.NewAuth(auth),
	})

	if err = app.Run(); err != nil {
		log.Fatal(err)
	}
}
