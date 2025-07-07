package route

import (
	"net/http"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/sakhaclothing/config"
	"github.com/sakhaclothing/controller"
)

func URL(w http.ResponseWriter, r *http.Request) {
	// Configure Fiber with optimized settings to prevent header size issues
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ServerHeader:          "Sakha API",
		AppName:               "Sakha Clothing API",
		ReadTimeout:           30,
		WriteTimeout:          30,
		IdleTimeout:           120,
		ReadBufferSize:        4096,
		WriteBufferSize:       4096,
	})

	// koneksi DB
	config.ConnectDB()

	// CORS middleware - simplified to prevent header size issues
	app.Use(func(c *fiber.Ctx) error {
		// Handle preflight requests
		if c.Method() == "OPTIONS" {
			c.Set("Access-Control-Allow-Origin", "*")
			c.Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			return c.SendStatus(http.StatusNoContent)
		}

		// Set minimal CORS headers for actual requests
		c.Set("Access-Control-Allow-Origin", "*")
		return c.Next()
	})

	// routes
	app.Post("/auth/:action", controller.AuthHandler)
	adaptor.FiberApp(app).ServeHTTP(w, r)
}

func SetupRoutes(app *fiber.App) {
	app.Post("/auth/:action", controller.AuthHandler)
}
