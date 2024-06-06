package webapi

import (
	"fmt"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tyasheliy/code_rooms/services/auth/internal/config"
	"github.com/tyasheliy/code_rooms/services/auth/internal/webapi/handler"
)

type WebApiApp struct {
	cfg  *config.AppConfig
	echo *echo.Echo
}

type Handlers struct {
	User *handler.UserHandler
	Auth *handler.AuthHandler
}

func NewWebApiApp(cfg *config.AppConfig, h *Handlers) *WebApiApp {
	e := echo.New()

	api := e.Group("/api/v1/", middleware.CORSWithConfig(middleware.CORSConfig{}))

	users := api.Group("users", echojwt.JWT([]byte(cfg.Secret)))
	users.GET("/:id", h.User.GetById)

	auth := api.Group("auth")
	auth.POST("/signin", h.Auth.Authenticate)
	auth.POST("/register", h.Auth.Register)

	return &WebApiApp{
		cfg:  cfg,
		echo: e,
	}
}

func (a *WebApiApp) Run() error {
	return a.echo.Start(fmt.Sprintf(":%s", a.cfg.WebApiPort))
}
