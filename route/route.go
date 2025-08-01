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

	// CORS configuration for production and development
	app.Use(func(c *fiber.Ctx) error {
		origin := c.Get("Origin")

		// Allowed origins
		allowedOrigins := []string{
			"https://sakhaclothing.shop",
			"http://127.0.0.1:5500",
			"http://localhost:5500",
			"http://127.0.0.1:3000",
			"http://localhost:3000",
			"http://127.0.0.1:8080",
			"http://localhost:8080",
			"http://127.0.0.1:5000",
			"http://localhost:5000",
		}

		// Check if origin is allowed
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				c.Set("Access-Control-Allow-Origin", allowedOrigin)
				break
			}
		}

		// Set CORS headers for all requests
		c.Set("Access-Control-Allow-Methods", "POST, GET, PUT, PATCH, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Set("Access-Control-Allow-Credentials", "true")
		c.Set("Vary", "Origin")

		// Handle preflight requests
		if c.Method() == "OPTIONS" {
			return c.SendStatus(204)
		}

		return c.Next()
	})

	// routes
	app.Post("/auth/google-login", controller.GoogleLoginHandler)
	app.Post("/auth/:action", controller.AuthHandler)
	app.Post("/tracker", controller.TrackerHandler)
	app.Get("/tracker/count", controller.TrackerCountHandler)
	app.Get("/config/google-client-id", controller.GetGoogleClientID)
	app.Put("/auth/change-password", controller.ChangePasswordHandler)
	app.Put("/auth/profile", controller.UpdateProfileHandler)
	app.Post("/auth/verify-email", controller.AuthHandler)
	app.Post("/auth/resend-otp", controller.AuthHandler)

	// Product routes
	app.Get("/products", controller.GetProducts)
	app.Get("/products/:id", controller.GetProduct)
	app.Post("/products", controller.CreateProduct)
	app.Put("/products/:id", controller.UpdateProduct)
	app.Delete("/products/:id", controller.DeleteProduct)
	app.Patch("/products/:id/featured", controller.ToggleFeatured)

	// Newsletter routes
	app.Post("/newsletter/subscribe", controller.SubscribeNewsletter)
	app.Get("/newsletter/unsubscribe", controller.UnsubscribeNewsletter)
	app.Get("/newsletter/subscribers", controller.GetSubscribers)
	app.Post("/newsletter/notify/:id", controller.SendNewProductNotification)
	app.Get("/newsletter/history", controller.GetNotificationHistory)

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
	app.Get("/tracker/count", controller.TrackerCountHandler)
	app.Get("/config/google-client-id", controller.GetGoogleClientID)
	app.Put("/auth/change-password", controller.ChangePasswordHandler)
	app.Put("/auth/profile", controller.UpdateProfileHandler)
	app.Post("/auth/verify-email", controller.AuthHandler)

	// Product routes
	app.Get("/products", controller.GetProducts)
	app.Get("/products/:id", controller.GetProduct)
	app.Post("/products", controller.CreateProduct)
	app.Put("/products/:id", controller.UpdateProduct)
	app.Delete("/products/:id", controller.DeleteProduct)
	app.Patch("/products/:id/featured", controller.ToggleFeatured)

	// Newsletter routes
	app.Post("/newsletter/subscribe", controller.SubscribeNewsletter)
	app.Get("/newsletter/unsubscribe", controller.UnsubscribeNewsletter)
	app.Get("/newsletter/subscribers", controller.GetSubscribers)
	app.Post("/newsletter/notify/:id", controller.SendNewProductNotification)
	app.Get("/newsletter/history", controller.GetNotificationHistory)
}
