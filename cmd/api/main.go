package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/synchhans/procurement-system/internal/database"
	"github.com/synchhans/procurement-system/internal/handlers"
	"github.com/synchhans/procurement-system/internal/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	database.ConnectDB()

	app := fiber.New()

	// app.Use(cors.New())
	app.Use(cors.New())

	// Auth routes
	api := app.Group("/api")
	api.Post("/register", handlers.Register)
	api.Post("/login", handlers.Login)

	// Protected routes
	protected := api.Group("", middleware.Protected())

	// Master Data
	protected.Get("/suppliers", handlers.GetSuppliers)
	protected.Post("/suppliers", handlers.CreateSupplier)
	protected.Put("/suppliers/:id", handlers.UpdateSupplier)    // <--- BARU
	protected.Delete("/suppliers/:id", handlers.DeleteSupplier) // <--- BARU

	protected.Get("/items", handlers.GetItems)
	protected.Post("/items", handlers.CreateItem)
	protected.Put("/items/:id", handlers.UpdateItem)    // <--- BARU
	protected.Delete("/items/:id", handlers.DeleteItem) // <--- BARU

	// Transactions
	protected.Post("/purchase", handlers.CreatePurchase)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}
