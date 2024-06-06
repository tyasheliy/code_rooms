package main

import (
	"context"
	"github.com/tyasheliy/code_rooms/services/auth/internal/config"
	"github.com/tyasheliy/code_rooms/services/auth/internal/entity"
	"github.com/tyasheliy/code_rooms/services/auth/internal/repo/factory"
	"github.com/tyasheliy/code_rooms/services/auth/internal/storage"
	"github.com/tyasheliy/code_rooms/services/auth/pkg/v1/hasher"
	"log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

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

	userRepo, err := factory.CreateUserRepo(cfg.Driver, s)
	if err != nil {
		log.Fatal(err)
	}

	h := hasher.BcryptHasher{}

	p, err := h.Hash("admin")
	if err != nil {
		log.Fatal(err)
	}

	_, _ = userRepo.Create(context.Background(), &entity.User{
		RoleId:       2,
		Login:        "admin",
		PasswordHash: p,
	})
}
