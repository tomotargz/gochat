package main

import (
	"github.com/gorilla/websocket"
	"log"
)

// Client is
type Client struct {
	conn  *websocket.Conn
	r     *Room
	write chan []byte
	name  string
}

func newClient(conn *websocket.Conn, r *Room, name string) *Client {
	return &Client{
		conn:  conn,
		r:     r,
		write: make(chan []byte),
		name:  name,
	}
}

func (c *Client) writeMsg() {
	defer func() {
		c.conn.Close()
		c.r.unregister <- c
	}()
	for {
		b, ok := <-c.write
		if !ok {
			return
		}
		c.conn.WriteMessage(websocket.TextMessage, b)
	}
}

func (c *Client) readMsg() {
	defer func() {
		c.conn.Close()
		c.r.unregister <- c
	}()
	for {
		_, m, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		c.r.broadcast <- []byte("[" + c.name + "] " + string(m))
	}
}
