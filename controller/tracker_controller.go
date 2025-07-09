package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sakhaclothing/config"
	"github.com/sakhaclothing/model"
)

// POST /tracker
func TrackerHandler(c *fiber.Ctx) error {
	var tracker model.Tracker
	if err := c.BodyParser(&tracker); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid data"})
	}
	tracker.TanggalAmbil = time.Now()
	_, err := config.DB.Collection("tracker").InsertOne(context.Background(), tracker)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save tracker data"})
	}
	return c.JSON(fiber.Map{"message": "Tracker data saved"})
}
