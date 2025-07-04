package route

import (
	"github.com/sakhaclothing/config"
	"github.com/sakhaclothing/controller"
	"github.com/sakhaclothing/middlewares"

	"net/http"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func URL(w http.ResponseWriter, r *http.Request) {
	app := fiber.New()

	// koneksi DB
	config.ConnectDB()

	// CORS middleware jika perlu
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if c.Method() == "OPTIONS" {
			return c.SendStatus(http.StatusOK)
		}
		return c.Next()
	})

	// routes
	app.Post("/auth/:action", controller.AuthHandler)
	app.Use("/ws", middlewares.Protected())
	app.Get("/ws", websocket.New(controller.WebSocketHandler()))
	app.Get("/users", middlewares.Protected(), controller.GetAllUsers)
	app.Put("/user/:id", middlewares.Protected(), controller.UpdateProfileByID)
	app.Get("/user/:id", middlewares.Protected(), controller.GetProfile)
	app.Get("/debug/token", middlewares.Protected(), controller.DebugToken)

	adaptor.FiberApp(app).ServeHTTP(w, r)
}

func SetupRoutes(app *fiber.App) {
	app.Post("/auth/:action", controller.AuthHandler)

	app.Use("/ws", middlewares.Protected())
	app.Get("/ws", websocket.New(controller.WebSocketHandler()))
	app.Get("/users", middlewares.Protected(), controller.GetAllUsers)
	app.Put("/user/:id", middlewares.Protected(), controller.UpdateProfileByID)
	app.Get("/user/:id", middlewares.Protected(), controller.GetProfile)
	app.Get("/debug/token", middlewares.Protected(), controller.DebugToken)
}
