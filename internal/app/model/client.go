package model

import (
	"github.com/gorilla/websocket"
)

// Client ...
type Client struct {
	Conn *websocket.Conn
	Send chan *Message
}

// NewClient ...
func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		Conn: conn,
		Send: make(chan *Message),
	}
}
