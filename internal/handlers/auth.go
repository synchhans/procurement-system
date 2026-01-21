package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/synchhans/procurement-system/internal/database"
	"github.com/synchhans/procurement-system/internal/models"
	"github.com/synchhans/procurement-system/pkg/utils"
	"gorm.io/gorm"
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
		return c.Status(500).JSON(fiber.Map{"error": "System error"})
	}

	user := models.User{
		Username: registerInput.Username,
		Password: hashedPassword,
		Role:     "staff",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		log.Printf("[REGISTER] Error DB: %v", err)

		return c.Status(409).JSON(fiber.Map{"error": "Username not available or system error"})
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
	err := database.DB.Where("username = ?", input.Username).First(&user).Error

	if err == gorm.ErrRecordNotFound {
		user.Password = "$2a$14$29.G1.k8k9q/././././././././././././././././././."
	}

	match := utils.CheckPasswordHash(input.Password, user.Password)

	if err != nil || !match {
		log.Printf("[AUTH] Percobaan login gagal untuk username: %s (IP: %s)", input.Username, c.IP())

		return c.Status(401).JSON(fiber.Map{"error": "Username atau Password salah"})
	}

	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not generate token"})
	}

	log.Printf("[AUTH] User '%s' (ID: %d) berhasil login", user.Username, user.ID)

	return c.JSON(fiber.Map{"token": token, "user": fiber.Map{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
	}})
}
