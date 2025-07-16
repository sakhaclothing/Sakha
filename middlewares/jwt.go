package middlewares

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sakhaclothing/utils"
)

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(401).JSON(fiber.Map{"error": "Token tidak ditemukan atau format salah"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := utils.ValidatePasetoToken(tokenString)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid"})
		}

		userID, ok := token.Get("user_id").(string)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid - user_id tidak ditemukan"})
		}

		username, ok := token.Get("username").(string)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid - username tidak ditemukan"})
		}

		c.Locals("user_id", userID)
		c.Locals("username", username)

		return c.Next()
	}
}
