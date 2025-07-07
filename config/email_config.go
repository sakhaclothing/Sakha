package config

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
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

type SMTPConfig struct {
	SMTPHost     string `bson:"SMTP_HOST"`
	SMTPPort     string `bson:"SMTP_PORT"`
	SMTPUsername string `bson:"SMTP_USERNAME"`
	SMTPPassword string `bson:"SMTP_PASSWORD"`
	FromEmail    string `bson:"FROM_EMAIL"`
	FromName     string `bson:"FROM_NAME"`
}

// GetSMTPConfig loads SMTP config from MongoDB collection 'configurations' with _id: 'smtp'
func GetSMTPConfig() (*SMTPConfig, error) {
	var cfg SMTPConfig
	err := DB.Collection("configurations").FindOne(context.Background(), bson.M{"_id": "smtp"}).Decode(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
