package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sakhaclothing/config"
	"github.com/sakhaclothing/model"
	"go.mongodb.org/mongo-driver/bson"
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

// GET /tracker/count
func TrackerCountHandler(c *fiber.Ctx) error {
	count, err := config.DB.Collection("tracker").CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get count"})
	}
	return c.JSON(fiber.Map{"count": count})
}
