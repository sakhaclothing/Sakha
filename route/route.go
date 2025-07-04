package route

import (
	"net/http"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/sakhaclothing/config"
	"github.com/sakhaclothing/controller"
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

	adaptor.FiberApp(app).ServeHTTP(w, r)
}

func SetupRoutes(app *fiber.App) {
	app.Post("/auth/:action", controller.AuthHandler)
}
