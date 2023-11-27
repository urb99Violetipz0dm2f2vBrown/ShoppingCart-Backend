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
	cart := models.Cart{}
	err := database.DB.Db.Model(&cart).Preload("Books").Find(&cart).Error
	if err != nil {
		return c.Status(400).JSON(nil)
	}

	return c.Status(200).JSON(cart)
}

// Function to add a Book to the Cart
func AddToCart(c *fiber.Ctx) error {
	book := []models.Book{}
	var cart models.Cart
	bookID, err := c.ParamsInt("id")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Book ID"})
	}

	// Retrieve the Book from the database based on the given ID
	if err := database.DB.Db.First(&book, bookID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Book not found"})
	}

	// Retrieve the existing cart
	if err := database.DB.Db.Preload("Books").First(&cart).Error; err != nil {
		// Create a new Cart if it doesn't exist
		cart = models.Cart{}
		database.DB.Db.Create(&cart)
	}

	// Add the Book to the Cart
	database.DB.Db.Model(&cart).Association("Books").Append(&book)

	// Save the changes to the database
	database.DB.Db.Save(&cart)

	// Return the updated Cart as a response
	return c.Status(200).JSON(cart)
}

// api.Delete("/cart/:id", handlers.RemoveFromCart)
func RemoveFromCart(c *fiber.Ctx) error {
	bookID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Book ID"})
	}

	// Retrieve the cart
	var cart models.Cart
	if err := database.DB.Db.Preload("Books").First(&cart).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Cart not found"})
	}

	// Find the book in the cart's Books slice
	var foundBook models.Book
	for i, b := range cart.Books {
		if b.ID == uint(bookID) {
			foundBook = cart.Books[i]
			break
		}
	}

	if foundBook.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Book not found in the cart"})
	}

	// Remove the book from the cart
	database.DB.Db.Model(&cart).Association("Books").Delete(&foundBook)

	// Save the changes to the database
	database.DB.Db.Save(&cart)

	return c.Status(200).JSON(cart)
}
