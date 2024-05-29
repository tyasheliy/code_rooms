package handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/tyasheliy/code_rooms/services/editor/internal/usecase"
	"net/http"
)

type EntryHandler struct {
	service *usecase.EntryUseCase
}

func NewEntryHandler(service *usecase.EntryUseCase) *EntryHandler {
	return &EntryHandler{
		service: service,
	}
}

type createEntryRequest struct {
	SessionId string `json:"session_id"`
}

func (h *EntryHandler) Create(ctx echo.Context) error {
	var req createEntryRequest

	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	sessionId, err := uuid.Parse(req.SessionId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid session id")
	}

	entry, err := h.service.Create(context.Background(), sessionId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusCreated, entry)
}
