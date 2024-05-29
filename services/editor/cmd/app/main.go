package main

import (
	"github.com/tyasheliy/code_rooms/services/editor/internal/config"
	"github.com/tyasheliy/code_rooms/services/editor/internal/repo"
	"github.com/tyasheliy/code_rooms/services/editor/internal/usecase"
	"github.com/tyasheliy/code_rooms/services/editor/internal/webapi"
	"github.com/tyasheliy/code_rooms/services/editor/internal/ws"
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

	kb := cache.NewKeyBuilder(cache.KeyBuilderOptions{
		Separator: "_",
	}).Pre("editor")

	c := cache.NewInMemoryAppCache(kb, cfg.Expiration, cfg.Cleanup)

	// session
	sessionRepo := repo.NewCacheSessionRepo(c)
	sessionService := usecase.NewSessionUseCase(l, sessionRepo)

	// entry
	entryRepo := repo.NewCacheEntryRepo(c)
	entryService := usecase.NewEntryUseCase(l, entryRepo)

	// source
	sourceRepo := repo.NewFileSourceRepo(cfg.SourceDir)
	sourceService := usecase.NewSourceUseCase(l, sourceRepo)

	wsApp := ws.NewWsApp(
		l,
		entryService,
		sessionService,
		sourceService,
	)

	go func() {
		err = wsApp.Run(cfg.SocketPort)
		if err != nil {
			log.Fatal(err)
		}
	}()

	webapiApp := webapi.NewWebApiApp(l, entryService, sessionService)

	err = webapiApp.Run(cfg.WebApiPort)
	if err != nil {
		log.Fatal(err)
	}
}
