package handler

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/tyasheliy/code_rooms/services/editor/internal/usecase"
	"net/http"
)

type SessionHandler struct {
	service *usecase.SessionUseCase
}

func NewSessionHandler(service *usecase.SessionUseCase) *SessionHandler {
	return &SessionHandler{
		service: service,
	}
}

func (h *SessionHandler) Create(ctx echo.Context) error {
	sess, err := h.service.Create(context.Background())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusCreated, sess)
}
