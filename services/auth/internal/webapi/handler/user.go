package handler

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/tyasheliy/code_rooms/services/auth/internal/usecase"
	"github.com/tyasheliy/code_rooms/services/auth/internal/webapi/handler/response"
	"net/http"
	"strconv"
)

type UserHandler struct {
	user *usecase.UserUseCase
}

func NewUser(user *usecase.UserUseCase) *UserHandler {
	return &UserHandler{user: user}
}

func (h *UserHandler) GetById(ctx echo.Context) error {
	rawId := ctx.Param("id")
	id, err := strconv.ParseInt(rawId, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.user.GetById(context.Background(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, response.UserResponse{
		Id:     user.Id,
		Login:  user.Login,
		RoleId: user.RoleId,
	})
}
