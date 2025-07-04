package controller

import (
	"context"
	"log"
	"net/http"

	"github.com/sakhaclothing/config"
	"github.com/sakhaclothing/model"

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

func UpdateProfileByID(c *fiber.Ctx) error {
	// Ambil user ID dari URL param
	paramID := c.Params("id")

	// Ambil user ID dari JWT token
	tokenUserID := c.Locals("user_id").(string)
	usernameFromToken := c.Locals("username").(string)

	// Debug log untuk melihat nilai yang dibandingkan
	log.Printf("DEBUG - paramID: %s, tokenUserID: %s", paramID, tokenUserID)

	// Validasi ID format
	objID, err := primitive.ObjectIDFromHex(paramID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	// Validasi format token user ID juga
	tokenObjID, err := primitive.ObjectIDFromHex(tokenUserID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Token user ID tidak valid"})
	}

	// Cek apakah ID dari token dan URL cocok (bandingkan sebagai ObjectID)
	if objID != tokenObjID {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Tidak diizinkan mengubah data user lain",
			"debug": fiber.Map{
				"param_id":      paramID,
				"token_user_id": tokenUserID,
			},
		})
	}

	// Ambil data user dari input body
	var input model.User
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Data tidak valid"})
	}

	// Validasi jika mau tambah pengecekan username juga
	if input.Username != "" && input.Username != usernameFromToken {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Username tidak cocok dengan token"})
	}

	update := bson.M{
		"$set": bson.M{
			"username": input.Username,
			"email":    input.Email,
			"fullname": input.Fullname,
		},
	}

	_, err = config.DB.Collection("users").UpdateByID(context.Background(), objID, update)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal update profil"})
	}

	return c.JSON(fiber.Map{
		"message": "Profil berhasil diupdate",
	})
}

func GetProfile(c *fiber.Ctx) error {
	// Ambil user ID dari URL param
	paramID := c.Params("id")

	// Ambil user ID dari JWT token
	tokenUserID := c.Locals("user_id").(string)

	// Debug log
	log.Printf("DEBUG GetProfile - paramID: %s, tokenUserID: %s", paramID, tokenUserID)

	// Validasi ID format dari parameter
	objID, err := primitive.ObjectIDFromHex(paramID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	// Validasi format token user ID
	tokenObjID, err := primitive.ObjectIDFromHex(tokenUserID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Token user ID tidak valid"})
	}

	// Cek apakah user hanya bisa melihat profil sendiri
	if objID != tokenObjID {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Tidak diizinkan melihat data user lain",
			"debug": fiber.Map{
				"param_id":      paramID,
				"token_user_id": tokenUserID,
			},
		})
	}

	var user model.User
	err = config.DB.Collection("users").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	user.Password = "" // jangan tampilkan password
	return c.JSON(user)
}

// Debug endpoint untuk melihat informasi token
func DebugToken(c *fiber.Ctx) error {
	tokenUserID := c.Locals("user_id").(string)
	usernameFromToken := c.Locals("username").(string)

	return c.JSON(fiber.Map{
		"token_user_id": tokenUserID,
		"username":      usernameFromToken,
		"message":       "Token info berhasil diambil",
	})
}
