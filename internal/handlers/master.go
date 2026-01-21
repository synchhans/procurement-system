package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/synchhans/procurement-system/internal/database"
	"github.com/synchhans/procurement-system/internal/models"
)

// Supplier Handlers
func GetSuppliers(c *fiber.Ctx) error {
	var suppliers []models.Supplier
	database.DB.Find(&suppliers)
	return c.JSON(suppliers)
}

func CreateSupplier(c *fiber.Ctx) error {
	var supplier models.Supplier
	if err := c.BodyParser(&supplier); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	database.DB.Create(&supplier)
	return c.Status(201).JSON(supplier)
}

// Item Handlers
func GetItems(c *fiber.Ctx) error {
	var items []models.Item
	database.DB.Find(&items)
	return c.JSON(items)
}

func CreateItem(c *fiber.Ctx) error {
	var item models.Item
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	database.DB.Create(&item)
	return c.Status(201).JSON(item)
}
