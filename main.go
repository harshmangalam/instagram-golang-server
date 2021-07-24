package main

import (
	"fmt"
	"instagram/config"
	"instagram/database"
	"instagram/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowCredentials: true,
	}))

	config.LoadEnv()

	database.CreateConnection()

	routes.InitRoutes(app)

	log.Fatal(app.Listen(fmt.Sprintf(":%s", config.Get("PORT"))))

	defer database.Client.Close()
}
