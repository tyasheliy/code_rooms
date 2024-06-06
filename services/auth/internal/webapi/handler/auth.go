package handler

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/tyasheliy/code_rooms/services/auth/internal/usecase"
	"github.com/tyasheliy/code_rooms/services/auth/internal/webapi/handler/request"
	"github.com/tyasheliy/code_rooms/services/auth/internal/webapi/handler/response"
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

	access, refresh, err := h.service.Authenticate(context.Background(), req.Login, req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, response.AuthenticateResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

func (h *AuthHandler) Register(ctx echo.Context) error {
	var req request.RegisterRequest

	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	user, err := h.service.Register(context.Background(), req.Login, req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusCreated, response.UserResponse{
		Id:     user.Id,
		Login:  user.Login,
		RoleId: user.RoleId,
	})
}
