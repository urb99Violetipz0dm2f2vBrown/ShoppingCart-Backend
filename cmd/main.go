package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tbui1996/backend-practice/database"
)

func main() {
	database.ConnectDb()
	app := fiber.New()

	setupRoutes(app)

	app.Listen(":3000")
}
