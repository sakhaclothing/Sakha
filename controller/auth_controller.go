package controller

import (
	"context"
	"strings"
	"time"

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

	switch action {
	case "profile":
		// Profile action doesn't need body parsing
		// Ambil token dari header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Token tidak ditemukan"})
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := utils.ValidatePasetoToken(tokenStr)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid"})
		}
		userId, ok := token.Get("user_id").(string)
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
			"role":     user.Role,
		})

	case "forgot-password":
		// Ultra-minimal forgot password handler
		var input model.ForgotPasswordRequest
		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).SendString("Invalid data")
		}

		if input.Email == "" {
			return c.Status(400).SendString("Email required")
		}

		// Check if email exists
		var user model.User
		err := config.DB.Collection("users").FindOne(context.Background(), bson.M{
			"email": bson.M{"$regex": "^" + input.Email + "$", "$options": "i"},
		}).Decode(&user)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.SendString("Email sent if registered")
			}
			return c.Status(500).SendString("Database error")
		}

		// Generate token
		token, expiresAt, err := utils.GenerateResetTokenWithExpiry()
		if err != nil {
			return c.Status(500).SendString("Token error")
		}

		// Save token
		resetToken := model.PasswordReset{
			Email:     input.Email,
			Token:     token,
			ExpiresAt: expiresAt,
			Used:      false,
			CreatedAt: time.Now(),
		}

		_, err = config.DB.Collection("password_resets").InsertOne(context.Background(), resetToken)
		if err != nil {
			return c.Status(500).SendString("Save error")
		}

		// Kirim link ke frontend, bukan ke endpoint API backend
		resetLink := "https://sakhaclothing.shop/reset-password/?token=" + token
		_ = utils.SendPasswordResetEmail(input.Email, token, resetLink)

		return c.SendString("Email sent")

	case "reset-password":
		// Handle reset password request
		var input model.ResetPasswordRequest
		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Data tidak valid"})
		}

		// Validasi input
		if input.Token == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Token tidak boleh kosong"})
		}
		if input.Password == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Password tidak boleh kosong"})
		}

		// Validasi panjang password
		if len(input.Password) < 6 {
			return c.Status(400).JSON(fiber.Map{"error": "Password minimal 6 karakter"})
		}

		// Cari token reset di database
		var resetToken model.PasswordReset
		err := config.DB.Collection("password_resets").FindOne(context.Background(), bson.M{
			"token": input.Token,
			"used":  false,
		}).Decode(&resetToken)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(400).JSON(fiber.Map{"error": "Token reset tidak valid atau sudah digunakan"})
			}
			return c.Status(500).JSON(fiber.Map{"error": "Gagal memeriksa token"})
		}

		// Cek apakah token sudah expired
		if utils.IsTokenExpired(resetToken.ExpiresAt) {
			return c.Status(400).JSON(fiber.Map{"error": "Token reset sudah expired"})
		}

		// Cari user berdasarkan email
		var user model.User
		err = config.DB.Collection("users").FindOne(context.Background(), bson.M{
			"email": resetToken.Email,
		}).Decode(&user)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
			}
			return c.Status(500).JSON(fiber.Map{"error": "Gagal mencari user"})
		}

		// Hash password baru
		hashedPassword, err := utils.HashPassword(input.Password)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal hash password"})
		}

		// Update password user
		_, err = config.DB.Collection("users").UpdateOne(
			context.Background(),
			bson.M{"_id": user.ID},
			bson.M{"$set": bson.M{"password": hashedPassword}},
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal update password"})
		}

		// Mark token sebagai used
		_, err = config.DB.Collection("password_resets").UpdateOne(
			context.Background(),
			bson.M{"_id": resetToken.ID},
			bson.M{"$set": bson.M{"used": true}},
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal mark token sebagai used"})
		}

		return c.JSON(fiber.Map{
			"message": "Password berhasil direset",
		})

	case "register", "check-username", "login":
		// For other actions, parse the body
		var input model.User
		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Data tidak valid"})
		}

		// --- Tambahkan verifikasi Turnstile di login/register ---
		if action == "register" || action == "login" || action == "reset-password" {
			token := c.FormValue("cf-turnstile-response")
			remoteip := c.IP()
			if !utils.VerifyTurnstile(token, remoteip) {
				return c.Status(400).JSON(fiber.Map{"error": "Verifikasi CAPTCHA gagal"})
			}
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

			// Cek apakah email sudah ada (case insensitive)
			var existingEmailUser model.User
			err = config.DB.Collection("users").FindOne(context.Background(), bson.M{
				"email": bson.M{"$regex": "^" + input.Email + "$", "$options": "i"},
			}).Decode(&existingEmailUser)
			if err == nil {
				return c.Status(409).JSON(fiber.Map{"error": "Email sudah digunakan"})
			} else if err != mongo.ErrNoDocuments {
				return c.Status(500).JSON(fiber.Map{"error": "Gagal memeriksa email"})
			}

			// Hash password
			hashed, err := utils.HashPassword(input.Password)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Gagal hash password"})
			}
			input.Password = hashed

			// Set user belum terverifikasi
			input.IsVerified = false
			input.Role = "user" // Set role default

			// Insert user baru
			result, err := config.DB.Collection("users").InsertOne(context.Background(), input)
			if err != nil {
				// Cek apakah error karena duplicate key (unique constraint)
				if mongo.IsDuplicateKeyError(err) {
					return c.Status(409).JSON(fiber.Map{"error": "Username sudah digunakan"})
				}
				return c.Status(500).JSON(fiber.Map{"error": "Gagal register user"})
			}

			// Generate OTP
			otp, err := utils.GenerateOTP()
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Gagal generate OTP"})
			}
			expiresAt := time.Now().Add(10 * time.Minute)
			verification := model.EmailVerification{
				Email:     input.Email,
				OTP:       otp,
				ExpiresAt: expiresAt,
				Used:      false,
				CreatedAt: time.Now(),
			}
			_, err = config.DB.Collection("email_verifications").InsertOne(context.Background(), verification)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Gagal simpan OTP"})
			}

			// Kirim OTP ke email user
			subject := "Verifikasi Email - Sakha Clothing"
			body := "<p>Kode OTP verifikasi email Anda: <b>" + otp + "</b></p><p>Kode berlaku 10 menit.</p>"
			smtpConfig, smtpErr := config.GetSMTPConfig()
			if smtpErr != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil konfigurasi SMTP"})
			}
			err = utils.SendSMTPEmail(smtpConfig, input.Email, subject, body)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Gagal mengirim email OTP"})
			}

			return c.Status(201).JSON(fiber.Map{
				"message": "Berhasil register. Silakan cek email untuk verifikasi (OTP)",
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
				return c.Status(400).JSON(fiber.Map{"error": "Username atau email tidak boleh kosong"})
			}
			if input.Password == "" {
				return c.Status(400).JSON(fiber.Map{"error": "Password tidak boleh kosong"})
			}

			loginField := "username"
			loginValue := strings.ToLower(strings.TrimSpace(input.Username))
			if strings.Contains(loginValue, "@") {
				loginField = "email"
			}

			var user model.User
			err := config.DB.Collection("users").FindOne(context.Background(), bson.M{
				loginField: bson.M{"$regex": "^" + loginValue + "$", "$options": "i"},
			}).Decode(&user)

			if err != nil {
				if err == mongo.ErrNoDocuments {
					return c.Status(404).JSON(fiber.Map{"error": "Username/email atau password salah"})
				}
				return c.Status(500).JSON(fiber.Map{"error": "Gagal mencari user"})
			}

			if !utils.CheckPasswordHash(input.Password, user.Password) {
				return c.Status(401).JSON(fiber.Map{"error": "Username/email atau password salah"})
			}

			token, err := utils.GeneratePasetoToken(user.ID.Hex(), user.Username)
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
					"role":     user.Role,
				},
			})
		}

	case "check-email":
		var input struct {
			Email string `json:"email"`
		}
		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).JSON(fiber.Map{"valid": false, "error": "Data tidak valid"})
		}
		if input.Email == "" {
			return c.Status(400).JSON(fiber.Map{"valid": false, "error": "Email wajib diisi"})
		}
		var user model.User
		err := config.DB.Collection("users").FindOne(context.Background(), bson.M{
			"email": bson.M{"$regex": "^" + input.Email + "$", "$options": "i"},
		}).Decode(&user)
		if err == mongo.ErrNoDocuments {
			return c.JSON(fiber.Map{"valid": false, "error": "Email tidak terdaftar"})
		} else if err != nil {
			return c.Status(500).JSON(fiber.Map{"valid": false, "error": "Gagal cek email"})
		}
		return c.JSON(fiber.Map{"valid": true})

	case "verify-email":
		var input struct {
			Email string `json:"email"`
			OTP   string `json:"otp"`
		}
		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Data tidak valid"})
		}
		if input.Email == "" || input.OTP == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Email dan OTP wajib diisi"})
		}

		// Cari OTP yang belum digunakan dan belum expired
		var verification model.EmailVerification
		err := config.DB.Collection("email_verifications").FindOne(context.Background(), bson.M{
			"email": input.Email,
			"otp":   input.OTP,
			"used":  false,
		}).Decode(&verification)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(400).JSON(fiber.Map{"error": "OTP salah atau sudah digunakan"})
			}
			return c.Status(500).JSON(fiber.Map{"error": "Gagal cek OTP"})
		}
		if utils.IsTokenExpired(verification.ExpiresAt) {
			return c.Status(400).JSON(fiber.Map{"error": "OTP sudah expired"})
		}

		// Update user menjadi verified
		_, err = config.DB.Collection("users").UpdateOne(
			context.Background(),
			bson.M{"email": input.Email},
			bson.M{"$set": bson.M{"is_verified": true}},
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal update user"})
		}

		// Mark OTP as used
		_, err = config.DB.Collection("email_verifications").UpdateOne(
			context.Background(),
			bson.M{"_id": verification.ID},
			bson.M{"$set": bson.M{"used": true}},
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal update OTP"})
		}

		return c.JSON(fiber.Map{"message": "Email berhasil diverifikasi"})

	case "update-role":
		// Ambil token dari header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Token tidak ditemukan"})
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := utils.ValidatePasetoToken(tokenStr)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid"})
		}
		userId, ok := token.Get("user_id").(string)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid (user_id tidak ditemukan)"})
		}

		// Convert userId string ke ObjectID
		objID, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "User ID tidak valid"})
		}

		// Cari user yang melakukan request
		var adminUser model.User
		err = config.DB.Collection("users").FindOne(context.Background(), bson.M{
			"_id": objID,
		}).Decode(&adminUser)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
			}
			return c.Status(500).JSON(fiber.Map{"error": "Gagal mencari user"})
		}

		// Cek apakah user adalah admin
		if adminUser.Role != "admin" {
			return c.Status(403).JSON(fiber.Map{"error": "Akses ditolak. Hanya admin yang dapat mengubah role"})
		}

		// Parse input untuk update role
		var input struct {
			UserID string `json:"user_id"`
			Role   string `json:"role"`
		}
		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Data tidak valid"})
		}

		if input.UserID == "" {
			return c.Status(400).JSON(fiber.Map{"error": "User ID tidak boleh kosong"})
		}
		if input.Role == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Role tidak boleh kosong"})
		}

		// Validasi role yang diizinkan
		allowedRoles := []string{"user", "admin", "moderator"}
		roleValid := false
		for _, allowedRole := range allowedRoles {
			if input.Role == allowedRole {
				roleValid = true
				break
			}
		}
		if !roleValid {
			return c.Status(400).JSON(fiber.Map{"error": "Role tidak valid. Role yang diizinkan: user, admin, moderator"})
		}

		// Convert target user ID string ke ObjectID
		targetUserID, err := primitive.ObjectIDFromHex(input.UserID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "User ID target tidak valid"})
		}

		// Cek apakah target user ada
		var targetUser model.User
		err = config.DB.Collection("users").FindOne(context.Background(), bson.M{
			"_id": targetUserID,
		}).Decode(&targetUser)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(404).JSON(fiber.Map{"error": "User target tidak ditemukan"})
			}
			return c.Status(500).JSON(fiber.Map{"error": "Gagal mencari user target"})
		}

		// Update role user target
		_, err = config.DB.Collection("users").UpdateOne(
			context.Background(),
			bson.M{"_id": targetUserID},
			bson.M{"$set": bson.M{"role": input.Role}},
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal update role user"})
		}

		return c.JSON(fiber.Map{
			"message":  "Role user berhasil diupdate",
			"user_id":  input.UserID,
			"new_role": input.Role,
		})

	case "get-users":
		// Ambil token dari header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Token tidak ditemukan"})
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := utils.ValidatePasetoToken(tokenStr)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid"})
		}
		userId, ok := token.Get("user_id").(string)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid (user_id tidak ditemukan)"})
		}

		// Convert userId string ke ObjectID
		objID, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "User ID tidak valid"})
		}

		// Cari user yang melakukan request
		var adminUser model.User
		err = config.DB.Collection("users").FindOne(context.Background(), bson.M{
			"_id": objID,
		}).Decode(&adminUser)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
			}
			return c.Status(500).JSON(fiber.Map{"error": "Gagal mencari user"})
		}

		// Cek apakah user adalah admin
		if adminUser.Role != "admin" {
			return c.Status(403).JSON(fiber.Map{"error": "Akses ditolak. Hanya admin yang dapat melihat daftar user"})
		}

		// Ambil semua user (tanpa password)
		cursor, err := config.DB.Collection("users").Find(context.Background(), bson.M{})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil daftar user"})
		}
		defer cursor.Close(context.Background())

		var users []fiber.Map
		for cursor.Next(context.Background()) {
			var user model.User
			if err := cursor.Decode(&user); err != nil {
				continue
			}
			users = append(users, fiber.Map{
				"id":          user.ID.Hex(),
				"username":    user.Username,
				"email":       user.Email,
				"fullname":    user.Fullname,
				"role":        user.Role,
				"is_verified": user.IsVerified,
			})
		}

		return c.JSON(fiber.Map{
			"users": users,
			"total": len(users),
		})

	case "change-password":
		// PUT /user/password
		return ChangePasswordHandler(c)

	case "update-profile":
		// PUT /user/profile
		return UpdateProfileHandler(c)

	default:
		return c.Status(400).JSON(fiber.Map{"error": "Action tidak dikenali"})
	}

	return c.Status(400).JSON(fiber.Map{"error": "Action tidak dikenali"})
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

// ChangePasswordHandler adalah handler terpisah untuk ganti password user
func ChangePasswordHandler(c *fiber.Ctx) error {
	// Ambil token dari header Authorization
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Token tidak ditemukan"})
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := utils.ValidatePasetoToken(tokenStr)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid"})
	}
	userId, ok := token.Get("user_id").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid (user_id tidak ditemukan)"})
	}

	objID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "User ID tidak valid"})
	}

	// Ambil input
	type ChangePasswordInput struct {
		OldPassword        string `json:"old_password"`
		NewPassword        string `json:"new_password"`
		ConfirmNewPassword string `json:"confirm_new_password"`
	}
	var input ChangePasswordInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Data tidak valid"})
	}

	if input.NewPassword != input.ConfirmNewPassword {
		return c.Status(400).JSON(fiber.Map{"error": "Konfirmasi password baru tidak cocok"})
	}
	if len(input.NewPassword) < 6 {
		return c.Status(400).JSON(fiber.Map{"error": "Password baru minimal 6 karakter"})
	}

	// Cari user di database
	var user model.User
	err = config.DB.Collection("users").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	// Cek password lama
	if !utils.CheckPasswordHash(input.OldPassword, user.Password) {
		return c.Status(400).JSON(fiber.Map{"error": "Password lama salah"})
	}

	// Hash password baru
	hashed, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal hash password baru"})
	}

	// Update password di database
	_, err = config.DB.Collection("users").UpdateOne(
		context.Background(),
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{"password": hashed}},
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal update password"})
	}

	return c.JSON(fiber.Map{"message": "Password berhasil diganti"})
}

// UpdateProfileHandler untuk update username, fullname, email
func UpdateProfileHandler(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Token tidak ditemukan"})
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := utils.ValidatePasetoToken(tokenStr)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid"})
	}
	userId, ok := token.Get("user_id").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid (user_id tidak ditemukan)"})
	}
	objID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "User ID tidak valid"})
	}

	// Ambil input
	type UpdateProfileInput struct {
		Username string `json:"username"`
		Fullname string `json:"fullname"`
		Email    string `json:"email"`
	}
	var input UpdateProfileInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Data tidak valid"})
	}

	// Cek user lama
	var user model.User
	err = config.DB.Collection("users").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	// Jika email berubah, set is_verified=false dan TODO: generate & kirim OTP
	update := bson.M{
		"username": input.Username,
		"fullname": input.Fullname,
	}
	if input.Email != user.Email {
		update["email"] = input.Email
		update["is_verified"] = false
		// Generate OTP
		otp, err := utils.GenerateOTP()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal generate OTP"})
		}
		expiresAt := time.Now().Add(10 * time.Minute)
		verification := model.EmailVerification{
			Email:     input.Email,
			OTP:       otp,
			ExpiresAt: expiresAt,
			Used:      false,
			CreatedAt: time.Now(),
		}
		_, err = config.DB.Collection("email_verifications").InsertOne(context.Background(), verification)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal simpan OTP"})
		}
		// Kirim OTP ke email user
		subject := "Verifikasi Email Baru - Sakha Clothing"
		body := "<p>Kode OTP verifikasi email Anda: <b>" + otp + "</b></p><p>Kode berlaku 10 menit.</p>"
		smtpConfig, smtpErr := config.GetSMTPConfig()
		if smtpErr != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil konfigurasi SMTP"})
		}
		err = utils.SendSMTPEmail(smtpConfig, input.Email, subject, body)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal mengirim email OTP"})
		}
		_, err = config.DB.Collection("users").UpdateOne(
			context.Background(),
			bson.M{"_id": user.ID},
			bson.M{"$set": update},
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal update email"})
		}
		return c.JSON(fiber.Map{"message": "Kode OTP dikirim ke email baru. Silakan verifikasi."})
	}

	// Update data user (tanpa ganti email)
	_, err = config.DB.Collection("users").UpdateOne(
		context.Background(),
		bson.M{"_id": user.ID},
		bson.M{"$set": update},
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal update profil"})
	}

	return c.JSON(fiber.Map{"message": "Profil berhasil diperbarui!"})
}
