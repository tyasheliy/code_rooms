package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/tyasheliy/code_rooms/services/editor/internal/repo"
	"github.com/tyasheliy/code_rooms/services/editor/internal/usecase"
	"github.com/tyasheliy/code_rooms/services/editor/pkg/v1/logger"
	"os"
)

func main() {
	l := logger.NewSlogAppLogger(logger.AppLoggerOptions{
		Format: logger.TEXT_FORMAT,
		Level:  logger.INFO_LEVEL,
		Output: os.Stdout,
	})

	//keyBuilder := cache.NewKeyBuilder(cache.KeyBuilderOptions{
	//	Separator: "_",
	//})
	//keyBuilder.Pre("editor")
	//c := cache.NewInMemoryAppCache(keyBuilder, 5*time.Minute, 15*time.Minute)

	os.Mkdir("./source", 0750)

	sourceRepo := repo.NewFileSourceRepo("./source")
	service := usecase.NewSourceUseCase(l, sourceRepo)

	ctx := context.Background()

	sessionId := uuid.New()

	source, err := service.Create(ctx,
		sessionId,
		"test.txt",
		[]byte("test data"),
	)

	fmt.Println(source, err)

	source2, err := service.Create(ctx,
		uuid.New(),
		"test.txt",
		[]byte("test data"),
	)

	fmt.Println(source2, err)

	err = service.UpdateData(ctx, source.Id, []byte("updated data"))

	fmt.Println(err)

	sources, err := service.GetBySession(ctx, sessionId)

	fmt.Println(len(sources), err)

	//err = service.Delete(ctx, source.Id)
	//err = service.Delete(ctx, source2.Id)

	fmt.Println(err)
}
