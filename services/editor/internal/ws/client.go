package ws

import (
	"context"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type client struct {
	hub        *sessionHub
	conn       *websocket.Conn
	send       chan *message
	pingTicker *time.Ticker
}

func newClient(hub *sessionHub, conn *websocket.Conn) *client {
	t := time.NewTicker(30 * time.Second)

	return &client{
		hub:        hub,
		conn:       conn,
		send:       make(chan *message),
		pingTicker: t,
	}
}

func (c *client) process(ctx context.Context) {
	c.hub.join <- c
	defer func() {
		c.hub.leave <- c
		close(c.send)
		err := c.conn.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	go func() {
		for {
			select {
			case msg := <-c.send:
				if msg == nil {
					log.Println("message is nil")
					return
				}

				rawMsg := map[string]interface{}{
					"filename": msg.filename,
					"data":     string(msg.data),
				}
				err := c.conn.WriteJSON(&rawMsg)
				if err != nil {
					log.Println(err)
					return
				}
				break
			case <-c.pingTicker.C:
				err := c.conn.WriteMessage(websocket.PingMessage, []byte("ping"))
				if err != nil {
					log.Println(err)
					return
				}
				break
			}
		}
	}()

	for {
		_ = c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		c.conn.SetPongHandler(func(message string) error {
			_ = c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
			return nil
		})

		rawMsg := make(map[string]interface{})
		err := c.conn.ReadJSON(&rawMsg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg := message{
			filename: rawMsg["filename"].(string),
			data:     []byte(rawMsg["data"].(string)),
		}

		c.hub.broadcast <- &msg
	}
}
