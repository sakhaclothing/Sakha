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
		origin := c.Get("Origin")
		allowedOrigin := "https://sakhaclothing.shop"
		if origin == allowedOrigin {
			c.Set("Access-Control-Allow-Origin", allowedOrigin)
			c.Set("Vary", "Origin")
		}
		if c.Method() == "OPTIONS" {
			c.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Set("Access-Control-Allow-Credentials", "true")
			return c.SendStatus(204)
		}
		return c.Next()
	})

	// routes
	app.Post("/auth/google-login", controller.GoogleLoginHandler)
	app.Post("/auth/:action", controller.AuthHandler)
	app.Post("/tracker", controller.TrackerHandler)
	app.Get("/config/google-client-id", controller.GetGoogleClientID)

	// Simple test endpoint
	app.Post("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
	adaptor.FiberApp(app).ServeHTTP(w, r)
}

func SetupRoutes(app *fiber.App) {
	app.Post("/auth/google-login", controller.GoogleLoginHandler)
	app.Post("/auth/:action", controller.AuthHandler)
	app.Post("/tracker", controller.TrackerHandler)
	app.Get("/config/google-client-id", controller.GetGoogleClientID)
}