package controller

import (
	"github.com/sakhaclothing/sockets"
	"github.com/gofiber/websocket/v2"
)

func WebSocketHandler() func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		sockets.HandleWS(c)
	}
}
