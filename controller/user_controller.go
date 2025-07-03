package controller

import (
	"context"
	"log"
	"net/http"

	"github.com/WeChat-Easy-Chat/config"
	"github.com/WeChat-Easy-Chat/model"

	"github.com/gofiber/fiber/v2"
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
