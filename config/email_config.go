package config

import (
	"os"
)

// EmailConfig holds email configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

// GetEmailConfig returns email configuration from environment variables
func GetEmailConfig() EmailConfig {
	return EmailConfig{
		SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPUsername: getEnv("SMTP_USERNAME", "sakhaclothing@gmail.com"),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""), // Set this in environment
		FromEmail:    getEnv("FROM_EMAIL", "sakhaclothing@gmail.com"),
		FromName:     getEnv("FROM_NAME", "Sakha Clothing"),
	}
}

// getEnv helper function
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
