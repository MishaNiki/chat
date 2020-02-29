package hub

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	"github.com/MishaNiki/chat/internal/app/model"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// ClientHub ...
type ClientHub struct {
	hub *Hub
}

// WritePump ...
func (ch *ClientHub) WritePump(c *model.Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			msg, err := json.Marshal(message)
			if err != nil {
				log.Println("clientHub::55 error:", err)
			} else {
				w.Write(msg)
			}

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				msg, err = json.Marshal(message)
				if err != nil {
					log.Println("clientHub::65 error:", err)
				} else {
					w.Write(msg)
				}
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ReadPump ...
func (ch *ClientHub) ReadPump(c *model.Client) {
	defer func() {
		ch.hub.Unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(
		func(string) error {
			c.Conn.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		},
	)
	var message model.Message // bagovoe mesto
	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("clientHub::103 error: %v", err)
			}
			break
		}

		msg = bytes.TrimSpace(bytes.Replace(msg, []byte{'\n'}, []byte{' '}, -1))
		err = json.Unmarshal(msg, &message)
		if err != nil {
			log.Printf("clientHub::109 error: %v", err)
		} else {
			ch.hub.Broadcast <- &message
		}
	}
}
