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

	// Start a database transaction
	tx := database.DB.Db.Begin()

	// Create the book
	if err := tx.Create(&book).Error; err != nil {
		// Rollback the transaction on error
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Commit the transaction
	tx.Commit()

	return c.Status(fiber.StatusOK).JSON(book)
}

// DeleteAllBooks deletes all books from the database
func DeleteAllBooks(c *fiber.Ctx) error {
	// Hard delete all books using raw SQL query
	if result := database.DB.Db.Exec("DELETE FROM books"); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error deleting books",
			"error":   result.Error.Error(),
		})
	} else if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
			"message": "No books found to delete",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "All books deleted successfully",
	})
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
	id, _ := c.ParamsInt("id")
	book := models.Book{}
	err := database.DB.Db.Where("id = ?", id).First(&book).Error

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Book ID"})
	}

	// Check if a cart with ID 1 already exists
	var cart models.Cart
	result := database.DB.Db.First(&cart, 1)

	if result.RowsAffected == 0 {
		// If a cart with ID 1 doesn't exist, create a new one
		cart = models.Cart{}
		cart.ID = 1
		database.DB.Db.Create(&cart)
	}

	// Append the book to the cart's Books association (if the book exists)
	if book.ID != 0 {
		database.DB.Db.Model(&cart).Association("Books").Append(&book)
	}

	// Save the changes to the database
	database.DB.Db.Save(&cart)

	return c.Status(200).JSON(cart)
}

// api.Delete("/cart/:id", handlers.RemoveFromCart)
func RemoveFromCart(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")

	// Find the specific book to remove from the database
	book := models.Book{}
	err := database.DB.Db.First(&book, id).Error

	if err != nil {
		return c.Status(404).JSON(nil)
	}

	// Find the specific cart (ID always 1) from the database
	cart := models.Cart{}
	err = database.DB.Db.Preload("Books").First(&cart, 1).Error

	if err != nil {
		return c.Status(404).JSON(nil)
	}

	// Remove the book from the cart's Books association
	database.DB.Db.Model(&cart).Association("Books").Delete(&book)

	return c.Status(200).JSON(cart)
}

func CreateEmptyCart(c *fiber.Ctx) error {
	// Check if a cart with ID 1 already exists
	var existingCart models.Cart
	err := database.DB.Db.Preload("Books").First(&existingCart, 1).Error

	if err == nil {
		// If a cart with ID 1 already exists, return an error
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "A cart with ID 1 already exists",
		})
	}

	// Create a new empty cart with ID 1
	cart := models.Cart{}
	cart.ID = 1
	database.DB.Db.Create(&cart)

	return c.Status(fiber.StatusOK).JSON(cart)
}

func TemporaryResetHandler(c *fiber.Ctx) error {
	// Drop the "carts" table
	if err := database.DB.Db.Migrator().DropTable(&models.Cart{}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error dropping Cart table",
			"error":   err.Error(),
		})
	}

	// Recreate the "carts" table
	if err := database.DB.Db.AutoMigrate(&models.Cart{}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error recreating Cart table",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Cart table dropped and recreated successfully",
	})
}
