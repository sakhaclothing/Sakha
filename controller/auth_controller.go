package controller

import (
	"context"
	"github.com/sakhaclothing/config"
	"github.com/sakhaclothing/model"
	"github.com/sakhaclothing/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func AuthHandler(c *fiber.Ctx) error {
	action := c.Params("action")
	var input model.User
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Data tidak valid"})
	}

	switch action {
	case "register":
		exists := config.DB.Collection("users").FindOne(context.Background(), bson.M{"username": input.Username})
		if exists.Err() == nil {
			return c.Status(409).JSON(fiber.Map{"error": "Username sudah digunakan"})
		}
		hashed, err := utils.HashPassword(input.Password)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal hash password"})
		}
		input.Password = hashed
		_, err = config.DB.Collection("users").InsertOne(context.Background(), input)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal register"})
		}
		return c.JSON(fiber.Map{"message": "Berhasil register"})

	case "login":
		var user model.User
		err := config.DB.Collection("users").FindOne(context.Background(), bson.M{"username": input.Username}).Decode(&user)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
		}
		if !utils.CheckPasswordHash(input.Password, user.Password) {
			return c.Status(401).JSON(fiber.Map{"error": "Password salah"})
		}
		token, err := utils.GenerateJWT(user.ID.Hex(), user.Username)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal membuat token"})
		}

		return c.JSON(fiber.Map{"token": token})

	default:
		return c.Status(400).JSON(fiber.Map{"error": "Action tidak dikenali"})
	}
}
