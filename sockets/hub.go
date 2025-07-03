package sockets

import (
	"github.com/gofiber/websocket/v2"
)

var Clients = make(map[*websocket.Conn]bool)
var Broadcast = make(chan string)

func HandleWS(c *websocket.Conn) {
	defer func() {
		delete(Clients, c)
		c.Close()
	}()
	Clients[c] = true

	for {
		var msg string
		if err := c.ReadJSON(&msg); err != nil {
			break
		}
		Broadcast <- msg
	}
}

func StartHub() {
	for {
		msg := <-Broadcast
		for client := range Clients {
			client.WriteJSON(msg)
		}
	}
}
