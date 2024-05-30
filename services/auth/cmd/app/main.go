package main

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/tyasheliy/code_rooms/services/auth/internal/config"
	"github.com/tyasheliy/code_rooms/services/auth/internal/repo/factory"
	"github.com/tyasheliy/code_rooms/services/auth/internal/storage"
	"github.com/tyasheliy/code_rooms/services/auth/internal/usecase"
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

	db, err := storage.BuildDbConn(&cfg.StorageConfig)
	if err != nil {
		log.Fatal(err)
	}

	s := storage.New(db)
	defer s.Close()

	keyBuilder := cache.NewKeyBuilder(cache.KeyBuilderOptions{Separator: "_"})

	c := cache.NewInMemoryAppCache(keyBuilder, cfg.Expiration, cfg.Cleanup)

	h := &hasher.BcryptHasher{}

	jwtBuilder := jwtutils.NewBuilder(cfg.Secret, cfg.TokenConfig.Expiration, jwt.SigningMethodHS256)

	userRepo, err := factory.CreateUserRepo(cfg.Driver, s)
	if err != nil {
		log.Fatal(err)
	}

	auth := usecase.NewAuth(l, jwtBuilder, h, userRepo)
	_ = usecase.NewUser(l, c, userRepo)

	//_, _ = auth.Register(context.Background(), "test1", "test")
	t, _, err := auth.Authenticate(context.Background(), "test1", "test")
	fmt.Println(t, err)
}
