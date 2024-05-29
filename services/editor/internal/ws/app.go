package ws

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/tyasheliy/code_rooms/services/editor/internal/usecase"
	"github.com/tyasheliy/code_rooms/services/editor/pkg/v1/logger"
	"net/http"
)

type WsApp struct {
	logger   logger.AppLogger
	upgrader *websocket.Upgrader
	entry    *usecase.EntryUseCase
	session  *usecase.SessionUseCase
	source   *usecase.SourceUseCase
	hubs     map[*sessionHub]bool
}

func NewWsApp(
	logger logger.AppLogger,
	entry *usecase.EntryUseCase,
	session *usecase.SessionUseCase,
	source *usecase.SourceUseCase,
) *WsApp {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return &WsApp{
		hubs:     make(map[*sessionHub]bool),
		upgrader: &upgrader,
		logger:   logger,
		entry:    entry,
		session:  session,
		source:   source,
	}
}

func (a *WsApp) handler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	a.logger.Info(ctx, "received request",
		"request", r,
	)
	defer a.logger.Info(ctx, "request handled",
		"request", r,
	)

	rawEntryId := r.PathValue("entry")
	entryId, err := uuid.Parse(rawEntryId)
	if err != nil {
		a.logger.Warn(ctx, "failed to parse entry id",
			"request", r,
			"error", err,
		)
		return
	}

	entry, err := a.entry.GetById(ctx, entryId)
	if err != nil {
		return
	}

	session, err := a.session.GetById(ctx, entry.Session)
	if err != nil {
		return
	}

	conn, err := a.upgrader.Upgrade(w, r, nil)
	if err != nil {
		a.logger.Error(ctx, "error upgrading connection to websocket",
			"request", r,
			"error", err,
		)
		return
	}
	a.logger.Debug(ctx, "upgrading connection to websocket",
		"request", r,
	)

	var hub *sessionHub
	ok := false
	for h := range a.hubs {
		if h.sessionId == session.Id {
			hub = h
			ok = true
			break
		}
	}

	if !ok {
		hub = newSessionHub(a, session.Id)
		a.hubs[hub] = true
		go hub.serve(ctx)
		a.logger.Debug(ctx, "created new hub for session",
			"session", session,
			"hub", hub,
		)
	} else {
		a.logger.Debug(ctx, "found existing hub for session",
			"session", session,
			"hub", hub,
		)
	}

	cl := newClient(hub, conn)

	go cl.process(ctx)
}

func (a *WsApp) Run(port string) error {
	http.HandleFunc("/app/{entry}", a.handler)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)

	return err
}
