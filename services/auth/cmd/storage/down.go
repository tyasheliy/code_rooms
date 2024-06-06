package main

import (
	"github.com/tyasheliy/code_rooms/services/auth/internal/config"
	"github.com/tyasheliy/code_rooms/services/auth/internal/storage"
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

	err = s.Down()
	if err != nil {
		log.Fatal(err)
	}
}
