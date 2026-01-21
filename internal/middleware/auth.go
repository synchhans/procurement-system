package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Missing authorization header"})
		}
		parts := strings.Split(authHeader, "Bearer ")
		if len(parts) != 2 {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid token format"})
		}
		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Locals("user_id", uint(claims["user_id"].(float64)))
		c.Locals("role", claims["role"].(string))

		return c.Next()
	}
}
