package ws

import (
	"context"
	"github.com/google/uuid"
)

type sessionHub struct {
	sessionId uuid.UUID
	app       *WsApp
	clients   map[*client]bool
	broadcast chan *message
	join      chan *client
	leave     chan *client
}

func newSessionHub(app *WsApp, sessionId uuid.UUID) *sessionHub {
	return &sessionHub{
		app:       app,
		sessionId: sessionId,
		clients:   make(map[*client]bool),
		broadcast: make(chan *message),
		join:      make(chan *client),
		leave:     make(chan *client),
	}
}

func (h *sessionHub) serve(ctx context.Context) {
	// TODO: logging
	// TODO: protobuf
	for {
		select {
		case cl := <-h.join:
			h.clients[cl] = true

			sources, err := h.app.source.GetBySession(ctx, h.sessionId)
			if err != nil {
				continue
			}

			for _, source := range sources {
				cl.send <- &message{
					filename: source.Name,
					data:     source.Data,
				}
			}
		case cl := <-h.leave:
			delete(h.clients, cl)

			if len(h.clients) == 0 {
				h.shutdown()
				return
			}
		case msg := <-h.broadcast:
			err := h.app.source.UpdateDataByFilename(ctx, h.sessionId, msg.filename, msg.data)
			if err != nil {
				continue
			}

			for cl := range h.clients {
				cl.send <- msg
			}
		}
	}
}

func (h *sessionHub) shutdown() {
	close(h.broadcast)
	close(h.leave)
	close(h.join)

	ctx := context.Background()

	_ = h.app.source.DeleteBySession(ctx, h.sessionId)
	_ = h.app.session.Delete(ctx, h.sessionId)

	delete(h.app.hubs, h)
}
