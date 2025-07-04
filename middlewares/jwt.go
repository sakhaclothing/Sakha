package middlewares

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		jwtKey := []byte(getEnv("JWT_SECRET", "wechat_secret"))

		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(401).JSON(fiber.Map{"error": "Token tidak ditemukan atau format salah"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// Validasi algoritma
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(401, "Metode signing tidak valid")
			}
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid"})
		}

		// Pastikan user_id dan username dalam format string
		userID, ok := claims["user_id"].(string)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid - user_id tidak ditemukan"})
		}

		username, ok := claims["username"].(string)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid - username tidak ditemukan"})
		}

		c.Locals("user_id", userID)
		c.Locals("username", username)

		return c.Next()
	}
}
