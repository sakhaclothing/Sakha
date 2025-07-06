package controller

import (
	"context"
	"strings"

	"github.com/sakhaclothing/config"
	"github.com/sakhaclothing/model"
	"github.com/sakhaclothing/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func AuthHandler(c *fiber.Ctx) error {
	action := c.Params("action")
	var input model.User
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Data tidak valid"})
	}

	switch action {
	case "register":
		// Validasi input
		if input.Username == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Username tidak boleh kosong"})
		}
		if input.Password == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Password tidak boleh kosong"})
		}
		if input.Email == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Email tidak boleh kosong"})
		}
		if input.Fullname == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Fullname tidak boleh kosong"})
		}

		// Normalisasi username (lowercase untuk konsistensi)
		input.Username = strings.ToLower(strings.TrimSpace(input.Username))

		// Validasi format username (hanya huruf, angka, dan underscore)
		if !isValidUsername(input.Username) {
			return c.Status(400).JSON(fiber.Map{"error": "Username hanya boleh berisi huruf, angka, dan underscore"})
		}

		// Cek apakah username sudah ada (case insensitive)
		var existingUser model.User
		err := config.DB.Collection("users").FindOne(context.Background(), bson.M{
			"username": bson.M{"$regex": "^" + input.Username + "$", "$options": "i"},
		}).Decode(&existingUser)

		if err == nil {
			return c.Status(409).JSON(fiber.Map{"error": "Username sudah digunakan"})
		} else if err != mongo.ErrNoDocuments {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal memeriksa username"})
		}

		// Hash password
		hashed, err := utils.HashPassword(input.Password)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal hash password"})
		}
		input.Password = hashed

		// Insert user baru
		result, err := config.DB.Collection("users").InsertOne(context.Background(), input)
		if err != nil {
			// Cek apakah error karena duplicate key (unique constraint)
			if mongo.IsDuplicateKeyError(err) {
				return c.Status(409).JSON(fiber.Map{"error": "Username sudah digunakan"})
			}
			return c.Status(500).JSON(fiber.Map{"error": "Gagal register user"})
		}

		return c.Status(201).JSON(fiber.Map{
			"message": "Berhasil register",
			"user_id": result.InsertedID,
		})

	case "check-username":
		// Endpoint untuk mengecek ketersediaan username
		if input.Username == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Username tidak boleh kosong"})
		}

		// Normalisasi username
		username := strings.ToLower(strings.TrimSpace(input.Username))

		// Validasi format username
		if !isValidUsername(username) {
			return c.Status(400).JSON(fiber.Map{
				"available": false,
				"error":     "Username hanya boleh berisi huruf, angka, dan underscore (3-20 karakter)",
			})
		}

		// Cek apakah username sudah ada
		var existingUser model.User
		err := config.DB.Collection("users").FindOne(context.Background(), bson.M{
			"username": bson.M{"$regex": "^" + username + "$", "$options": "i"},
		}).Decode(&existingUser)

		if err == nil {
			return c.JSON(fiber.Map{
				"available": false,
				"message":   "Username sudah digunakan",
			})
		} else if err == mongo.ErrNoDocuments {
			return c.JSON(fiber.Map{
				"available": true,
				"message":   "Username tersedia",
			})
		} else {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal memeriksa username"})
		}

	case "login":
		// Validasi input
		if input.Username == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Username tidak boleh kosong"})
		}
		if input.Password == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Password tidak boleh kosong"})
		}

		// Normalisasi username untuk login
		username := strings.ToLower(strings.TrimSpace(input.Username))

		var user model.User
		err := config.DB.Collection("users").FindOne(context.Background(), bson.M{
			"username": bson.M{"$regex": "^" + username + "$", "$options": "i"},
		}).Decode(&user)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(404).JSON(fiber.Map{"error": "Username atau password salah"})
			}
			return c.Status(500).JSON(fiber.Map{"error": "Gagal mencari user"})
		}

		if !utils.CheckPasswordHash(input.Password, user.Password) {
			return c.Status(401).JSON(fiber.Map{"error": "Username atau password salah"})
		}

		token, err := utils.GenerateJWT(user.ID.Hex(), user.Username)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal membuat token"})
		}

		return c.JSON(fiber.Map{
			"message": "Login berhasil",
			"token":   token,
			"user": fiber.Map{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
				"fullname": user.Fullname,
			},
		})
	case "profile":
		// Ambil token dari header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Token tidak ditemukan"})
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		_, claims, err := utils.ValidateToken(tokenStr)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid"})
		}
	
		// Ambil user_id dari claims JWT
		userId, ok := claims["user_id"].(string)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid (user_id tidak ditemukan)"})
		}
	
		// Convert userId string ke ObjectID
		objID, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "User ID tidak valid"})
		}
	
		// Cari user di database berdasarkan _id
		var user model.User
		err = config.DB.Collection("users").FindOne(context.Background(), bson.M{
			"_id": objID,
		}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
			}
			return c.Status(500).JSON(fiber.Map{"error": "Gagal mencari user"})
		}
	
		// Return data user (tanpa password)
		return c.JSON(fiber.Map{
			"id":       user.ID.Hex(), // kirim sebagai string
			"username": user.Username,
			"email":    user.Email,
			"fullname": user.Fullname,
		})

	default:
		return c.Status(400).JSON(fiber.Map{"error": "Action tidak dikenali"})
	}
}

// isValidUsername memvalidasi format username
func isValidUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}

	// Hanya boleh berisi huruf, angka, dan underscore
	for _, char := range username {
		if !((char >= 'a' && char <= 'z') ||
			(char >= '0' && char <= '9') ||
			char == '_') {
			return false
		}
	}

	return true
}
