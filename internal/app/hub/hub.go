package hub

import (
	"github.com/MishaNiki/chat/internal/app/model"
	"github.com/gorilla/websocket"
)

// Hub ...
type Hub struct {
	Clients    map[*model.Client]bool
	Broadcast  chan *model.Message
	Register   chan *model.Client
	Unregister chan *model.Client
	ClientHub  *ClientHub
	Upgrader   *websocket.Upgrader
}

// New ...
func New() *Hub {
	return &Hub{
		Clients:    make(map[*model.Client]bool),
		Broadcast:  make(chan *model.Message),
		Register:   make(chan *model.Client),
		Unregister: make(chan *model.Client),
		Upgrader:   newUpgrader(),
	}
}

// Run ...
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}

// Create Upgrader
func newUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
}

// Client ...
// templete : hub.Client().WritePump()
func (h *Hub) Client() *ClientHub {
	if h.ClientHub != nil {
		return h.ClientHub
	}
	h.ClientHub = &ClientHub{
		hub: h,
	}
	return h.ClientHub
}
