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
	log.Println("✅ Terhubung ke MongoDB")
}

func SetEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Gagal memuat .env file, lanjutkan dengan os.Getenv")
	}
}