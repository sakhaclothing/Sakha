package utils

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"os"
	"time"

	"encoding/base64"

	"github.com/golang-jwt/jwt/v4"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte(getEnv("JWT_SECRET", "sakha_secret"))
var pasetoSecretKey = []byte(getEnv("PASETO_SECRET", "sakha_paseto_secret_32byteslong!!"))

// Get env with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// HashPassword hashes plain password
func HashPassword(pw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), 14)
	return string(hash), err
}

// CheckPasswordHash compares plain with hash
func CheckPasswordHash(pw, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw)) == nil
}

// GenerateJWT creates token for user
func GenerateJWT(id, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  id,
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString(jwtKey)
}

// ValidateToken parses token and returns claims
func ValidateToken(tokenStr string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, nil, jwt.ErrInvalidKey
	}
	return token, claims, nil
}

// GenerateResetToken generates a random token for password reset
func GenerateResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateResetTokenWithExpiry generates a reset token with expiry time
func GenerateResetTokenWithExpiry() (string, time.Time, error) {
	token, err := GenerateResetToken()
	if err != nil {
		return "", time.Time{}, err
	}

	// Token expires in 1 hour
	expiresAt := time.Now().Add(1 * time.Hour)
	return token, expiresAt, nil
}

// IsTokenExpired checks if a token has expired
func IsTokenExpired(expiresAt time.Time) bool {
	return time.Now().After(expiresAt)
}

// GenerateOTP generates a 6-digit numeric OTP for email verification
func GenerateOTP() (string, error) {
	var digits = []byte("0123456789")
	otp := make([]byte, 6)
	for i := range otp {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		otp[i] = digits[n.Int64()]
	}
	return string(otp), nil
}

// GeneratePasetoToken creates a PASETO token for user
func GeneratePasetoToken(id, username string) (string, error) {
	now := time.Now()
	exp := now.Add(24 * time.Hour)
	jsonToken := paseto.JSONToken{
		Expiration: exp,
		IssuedAt:   now,
		NotBefore:  now,
		Subject:    id,
		Audience:   "sakhaclothing",
		Issuer:     "sakhaclothing-backend",
		Jti:        randomJTI(),
		// Custom claims
		Set: map[string]interface{}{
			"user_id":  id,
			"username": username,
		},
	}
	footer := ""
	return paseto.NewV2().Encrypt(pasetoSecretKey, jsonToken, footer)
}

// ValidatePasetoToken parses and validates a PASETO token
func ValidatePasetoToken(tokenStr string) (*paseto.JSONToken, error) {
	var jsonToken paseto.JSONToken
	var footer string
	err := paseto.NewV2().Decrypt(tokenStr, pasetoSecretKey, &jsonToken, &footer)
	if err != nil {
		return nil, err
	}
	if jsonToken.Expiration.Before(time.Now()) {
		return nil, jwt.ErrTokenExpired
	}
	return &jsonToken, nil
}

// randomJTI generates a random string for JTI
func randomJTI() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
