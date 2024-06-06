package webapi

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tyasheliy/code_rooms/services/editor/internal/usecase"
	"github.com/tyasheliy/code_rooms/services/editor/internal/webapi/handler"
	"github.com/tyasheliy/code_rooms/services/editor/pkg/v1/logger"
)

type WebApiApp struct {
	logger logger.AppLogger
	echo   *echo.Echo
}

func NewWebApiApp(logger logger.AppLogger,
	entry *usecase.EntryUseCase,
	session *usecase.SessionUseCase,
) *WebApiApp {
	e := echo.New()

	api := e.Group("/api/v1", middleware.CORSWithConfig(middleware.CORSConfig{}))

	entries := api.Group("/entries")
	entryHandler := handler.NewEntryHandler(entry)
	entries.POST("", entryHandler.Create)

	sessions := api.Group("/sessions")
	sessionHandler := handler.NewSessionHandler(session)
	sessions.POST("", sessionHandler.Create)

	return &WebApiApp{
		logger: logger,
		echo:   e,
	}
}

func (a *WebApiApp) Run(port string) error {
	return a.echo.Start(fmt.Sprintf(":%s", port))
}
