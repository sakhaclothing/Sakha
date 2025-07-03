package route

import (
	"net/http"

	"github.com/WeChat-Easy-Chat/controller"
	"github.com/WeChat-Easy-Chat/middlewares"
	"github.com/WeChat-Easy-Chat/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func URL(w http.ResponseWriter, r *http.Request) {
	if config.SetAccessControlHeaders(w, r) {
		return // If it's a preflight request, return early.
	}
	config.SetEnv()
}

func SetupRoutes(app *fiber.App) {
	app.Post("/auth/:action", controller.AuthHandler)

	app.Use("/ws", middlewares.Protected())
	app.Get("/ws", websocket.New(controller.WebSocketHandler()))
}

