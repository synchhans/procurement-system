package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/synchhans/procurement-system/internal/database"
	"github.com/synchhans/procurement-system/internal/models"
	"github.com/synchhans/procurement-system/pkg/utils"
)

func Register(c *fiber.Ctx) error {
	var registerInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&registerInput); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if registerInput.Username == "" || registerInput.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Username and Password are required"})
	}

	hashedPassword, err := utils.HashPassword(registerInput.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not hash password"})
	}

	user := models.User{
		Username: registerInput.Username,
		Password: hashedPassword,
		Role:     "staff",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create user"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "User created successfully"})
}

func Login(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	var user models.User
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not generate token"})
	}

	return c.JSON(fiber.Map{"token": token, "user": fiber.Map{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
	}})
}
