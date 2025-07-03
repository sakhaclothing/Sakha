package controller

import (
	"context"
	"log"
	"net/http"

	"github.com/WeChat-Easy-Chat/config"
	"github.com/WeChat-Easy-Chat/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllUsers(c *fiber.Ctx) error {
	cursor, err := config.DB.Collection("users").Find(context.Background(), map[string]interface{}{})
	if err != nil {
		log.Println("DB error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data pengguna",
		})
	}
	defer cursor.Close(context.Background())

	var users []model.User
	for cursor.Next(context.Background()) {
		var user model.User
		if err := cursor.Decode(&user); err != nil {
			continue
		}
		user.Password = "" // Jangan tampilkan password
		users = append(users, user)
	}

	return c.JSON(users)
}

func UpdateProfile(c *fiber.Ctx) error {
	rawUser := c.Locals("user")
	userToken, ok := rawUser.(map[string]interface{})
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized atau token tidak valid",
		})
	}

	userIDStr, ok := userToken["user_id"].(string)
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID tidak valid",
		})
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID tidak valid",
		})
	}

	var input model.User
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Data tidak valid",
		})
	}

	update := bson.M{
		"$set": bson.M{
			"username": input.Username,
			"email":    input.Email,
			"fullname": input.Fullname,
		},
	}

	_, err = config.DB.Collection("users").UpdateByID(context.Background(), userID, update)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal update profil",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Profil berhasil diupdate",
	})
}
