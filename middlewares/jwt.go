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

func JWTMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Authorization header missing"})
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	_, claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid"})
	}
	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid - user_id tidak ditemukan"})
	}
	username, ok := claims["username"].(string)
	if !ok || username == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid - username tidak ditemukan"})
	}
	c.Locals("user_id", userID)
	c.Locals("username", username)
	return c.Next()
}
