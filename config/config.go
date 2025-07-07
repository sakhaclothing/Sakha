package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDB() {
	// Load .env (optional, hanya untuk local development)
	_ = godotenv.Load() // jangan log error kalau .env gak ada, aman abaikan

	// Ambil MONGOSTRING dari environment
	mongoURI := os.Getenv("MONGOSTRING")
	if mongoURI == "" {
		log.Fatal("ENV MONGOSTRING tidak ditemukan")
	}

	// Koneksi ke MongoDB
	clientOpts := options.Client().ApplyURI(mongoURI)
	client, err := mongo.NewClient(clientOpts)
	if err != nil {
		log.Fatalf("Gagal buat Mongo client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		log.Fatalf("Gagal konek ke MongoDB: %v", err)
	}

	DB = client.Database("sakha")

	// Buat unique index untuk username
	CreateUniqueIndexes()

	log.Println("✅ Terhubung ke MongoDB")
}

// CreateUniqueIndexes membuat index unik untuk username
func CreateUniqueIndexes() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Buat unique index untuk username
	usernameIndexModel := mongo.IndexModel{
		Keys: map[string]interface{}{
			"username": 1,
		},
		Options: options.Index().SetUnique(true).SetName("username_unique"),
	}

	_, err := DB.Collection("users").Indexes().CreateOne(ctx, usernameIndexModel)
	if err != nil {
		log.Printf("⚠️  Gagal membuat unique index untuk username: %v", err)
		log.Println("   (Ini normal jika index sudah ada)")
	} else {
		log.Println("✅ Unique index untuk username berhasil dibuat")
	}

	// Buat unique index untuk email
	emailIndexModel := mongo.IndexModel{
		Keys: map[string]interface{}{
			"email": 1,
		},
		Options: options.Index().SetUnique(true).SetName("email_unique"),
	}
	_, err = DB.Collection("users").Indexes().CreateOne(ctx, emailIndexModel)
	if err != nil {
		log.Printf("⚠️  Gagal membuat unique index untuk email: %v", err)
		log.Println("   (Ini normal jika index sudah ada)")
	} else {
		log.Println("✅ Unique index untuk email berhasil dibuat")
	}
}

func SetEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Gagal memuat .env file, lanjutkan dengan os.Getenv")
	}
}
