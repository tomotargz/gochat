package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Room is hub
type Room struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

func newRoom() *Room {
	return &Room{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

func auth(r *http.Request) (string, bool) {
	s, err := r.Cookie("SESSION")
	if err != nil {
		return "", false
	}

	name, ok := sessions[s.Value]
	if !ok {
		return "", false
	}
	return name, true
}

func (room *Room) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name, ok := auth(r)
	if !ok {
		http.Redirect(w, r, oauth2Conf.AuthCodeURL("state"), http.StatusFound)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := newClient(conn, room, name)
	go c.readMsg()
	go c.writeMsg()
	room.register <- c
}

func (room *Room) run() {
	for {
		select {
		case c := <-room.register:
			room.clients[c] = true
		case c := <-room.unregister:
			delete(room.clients, c)
		case b := <-room.broadcast:
			for c := range room.clients {
				c.write <- b
			}
		}
	}
}
