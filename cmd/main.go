package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/tbui1996/backend-practice/database"
)

func main() {
	database.ConnectDb()
	app := fiber.New()

	// Use CORS middleware to handle CORS headers
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3001", // Add other origins as needed
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Content-Type,Origin,Accept",
		AllowCredentials: true,
	}))

	setupRoutes(app)

	app.Listen(":3000")
}
