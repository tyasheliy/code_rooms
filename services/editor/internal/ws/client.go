package ws

import (
	"context"
	"github.com/gorilla/websocket"
)

type client struct {
	hub  *sessionHub
	conn *websocket.Conn
	send chan *message
}

func newClient(hub *sessionHub, conn *websocket.Conn) *client {
	return &client{
		hub:  hub,
		conn: conn,
		send: make(chan *message),
	}
}

func (c *client) process(ctx context.Context) {
	c.hub.join <- c
	defer func() {
		c.hub.leave <- c
		close(c.send)
		c.conn.Close()
	}()

	go func() {
		for {
			select {
			case msg := <-c.send:
				if msg == nil {
					// TODO: понять почему при отключении передает нуловое сообщение
					return
				}

				rawMsg := map[string]interface{}{
					"filename": msg.filename,
					"data":     string(msg.data),
				}
				err := c.conn.WriteJSON(&rawMsg)
				if err != nil {
					return
				}

				break
			}
		}
	}()

	for {
		rawMsg := make(map[string]interface{})
		err := c.conn.ReadJSON(&rawMsg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				return
			}

			continue
		}

		msg := message{
			filename: rawMsg["filename"].(string),
			data:     []byte(rawMsg["data"].(string)),
		}

		c.hub.broadcast <- &msg
	}
}
