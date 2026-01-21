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
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if err := godotenv.Load(); err != nil {
		log.Println("[INFO] File .env tidak ditemukan, menggunakan environment system.")
	} else {
		log.Println("[INFO] Konfigurasi .env berhasil dimuat.")
	}

	database.ConnectDB()

	app := fiber.New()
	app.Use(cors.New())

	api := app.Group("/api")

	api.Post("/register", handlers.Register)
	api.Post("/login", handlers.Login)

	protected := api.Group("", middleware.Protected())

	protected.Get("/suppliers", handlers.GetSuppliers)
	protected.Post("/suppliers", handlers.CreateSupplier)
	protected.Put("/suppliers/:id", handlers.UpdateSupplier)
	protected.Delete("/suppliers/:id", handlers.DeleteSupplier)

	protected.Get("/items", handlers.GetItems)
	protected.Post("/items", handlers.CreateItem)
	protected.Put("/items/:id", handlers.UpdateItem)
	protected.Delete("/items/:id", handlers.DeleteItem)

	protected.Post("/purchase", handlers.CreatePurchase)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("[INFO] Server Procurement berjalan di port :%s", port)
	log.Printf("[INFO] Siap menerima request...")

	log.Fatal(app.Listen(":" + port))
}
