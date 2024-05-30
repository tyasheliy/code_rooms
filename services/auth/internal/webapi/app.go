package webapi

import "github.com/labstack/echo/v4"

type WebApiApp struct {
	echo *echo.Echo
}

func NewWebApiApp() *WebApiApp {
	e := echo.New()

	return &WebApiApp{
		echo: e,
	}
}
