package controller

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sakhaclothing/config"
	"github.com/sakhaclothing/model"
	"github.com/sakhaclothing/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/api/idtoken"
)

func GoogleLoginHandler(c *fiber.Ctx) error {
	var body struct {
		Token string `json:"token"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Ambil Google Client ID dari config
	var conf model.Config
	err := config.DB.Collection("config").FindOne(context.Background(), bson.M{"key": "google_client_id"}).Decode(&conf)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Google Client ID not found"})
	}

	// Verifikasi token ke Google
	payload, err := idtoken.Validate(context.Background(), body.Token, conf.Value)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Google token"})
	}

	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)

	// Cari user di database
	var user model.User
	err = config.DB.Collection("users").FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		// Auto-register user baru
		user = model.User{
			Email:      email,
			Fullname:   name,
			Username:   email, // Atau generate username unik
			Role:       "user",
			IsVerified: true,
		}
		res, err := config.DB.Collection("users").InsertOne(context.Background(), user)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
		}
		if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
			user.ID = oid
		}
	} else if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	// Generate JWT aplikasi-mu (gunakan fungsi yang sudah ada)
	token, err := utils.GeneratePasetoToken(user.ID.Hex(), user.Username)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.JSON(fiber.Map{"token": token})
}
