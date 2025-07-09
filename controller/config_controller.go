package controller

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sakhaclothing/config"
	"github.com/sakhaclothing/model"
)

func GetGoogleClientID(c *fiber.Ctx) error {
	var conf model.Config
	err := config.DB.Collection("config").FindOne(context.Background(), map[string]interface{}{
		"key": "google_client_id",
	}).Decode(&conf)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Google Client ID not found"})
	}
	return c.JSON(fiber.Map{"client_id": conf.Value})
}
