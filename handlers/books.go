package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/tbui1996/backend-practice/database"
	"github.com/tbui1996/backend-practice/models"
)

func Ping(c *fiber.Ctx) error {
	return c.SendString("Pong")
}
func ListBooks(c *fiber.Ctx) error {
	//get params
	title := c.Params("title")
	author := c.Params("author")
	genre := c.Params("genre")
	books := []models.Book{}

	//make case insensitive
	title = strings.ToUpper(title)
	author = strings.ToUpper(author)
	genre = strings.ToUpper(genre)

	//query based on parameters
	if strings.EqualFold(title, "%20") && strings.EqualFold(author, "%20") && strings.EqualFold(genre, "%20") {
		database.DB.Db.Find(&books)
	} else {
		database.DB.Db.Where("UPPER(title) = ?", title).Or("UPPER(author) = ?", author).Or("UPPER(genre) = ?", genre).Find(&books)
	}

	return c.Status(200).JSON(books)
}

func CreateBook(c *fiber.Ctx) error {
	book := new(models.Book)
	if err := c.BodyParser(book); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	database.DB.Db.Create(&book)

	return c.Status(200).JSON(book)
}

func ListCart(c *fiber.Ctx) error {
	cart := []models.Cart{}
	err := database.DB.Db.Model(&cart).Preload("Books").Find(&cart).Error
	if err != nil {
		return c.Status(400).JSON(nil)
	}

	return c.Status(200).JSON(cart)
}

// api.Post("/cart/:id", handlers.AddToCart)
func AddToCart(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	cart := []models.Cart{}
	book := []models.Book{}
	err := database.DB.Db.Where("id = ?", id).First(&book).Error

	if err != nil {
		return c.Status(404).JSON(nil)
	}

	database.DB.Db.Model(&cart).Association("Books").Append(&book)
	database.DB.Db.Save(&cart)
	return c.Status(200).JSON(cart)
}

// api.Delete("/cart/:id", handlers.RemoveFromCart)
func RemoveFromCart(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	cart := []models.Cart{}
	book := []models.Book{}

	err := database.DB.Db.Where("id = ?", id).First(&book).Error

	if err != nil {
		return c.Status(404).JSON(nil)
	}

	database.DB.Db.Model(&cart).Association("Books").Delete(&book)
	return c.Status(200).JSON(cart)
}
