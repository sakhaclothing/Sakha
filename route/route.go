package route

import (
	"net/http"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/sakhaclothing/config"
	"github.com/sakhaclothing/controller"
)

func URL(w http.ResponseWriter, r *http.Request) {
	// Ultra-minimal Fiber config
	app := fiber.New(fiber.Config{
		DisableStartupMessage:     true,
		DisableDefaultDate:        true,
		DisableDefaultContentType: true,
	})

	// koneksi DB
	config.ConnectDB()

	// Minimal CORS - only for preflight
	app.Use(func(c *fiber.Ctx) error {
		if c.Method() == "OPTIONS" {
			c.Set("Access-Control-Allow-Origin", "*")
			c.Set("Access-Control-Allow-Methods", "POST")
			c.Set("Access-Control-Allow-Headers", "Content-Type")
			return c.SendStatus(http.StatusNoContent)
		}
		return c.Next()
	})

	// routes
	app.Post("/auth/:action", controller.AuthHandler)

	// Simple test endpoint
	app.Post("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
	adaptor.FiberApp(app).ServeHTTP(w, r)
}

func SetupRoutes(app *fiber.App) {
	app.Post("/auth/:action", controller.AuthHandler)
}
