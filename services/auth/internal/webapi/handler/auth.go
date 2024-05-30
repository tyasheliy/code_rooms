package handler

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/tyasheliy/code_rooms/services/auth/internal/usecase"
	"github.com/tyasheliy/code_rooms/services/auth/internal/webapi/handler/request"
	"net/http"
)

type AuthHandler struct {
	service *usecase.AuthUseCase
}

func NewAuth(service *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) Authenticate(ctx echo.Context) error {
	var req request.AuthenticateRequest

	err := ctx.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	h.service.Authenticate(context.Background(), req.Login, req.Password)
}

func (h *AuthHandler) Register(ctx echo.Context) error {

}
