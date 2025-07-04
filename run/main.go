package main

import (
	"log"
	"github.com/sakhaclothing/config"
	"github.com/sakhaclothing/route"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	config.ConnectDB()
	route.SetupRoutes(app)

	log.Println("Server berjalan di :8080")
	log.Fatal(app.Listen(":8080"))
}
