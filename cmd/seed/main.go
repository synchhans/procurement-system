package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/synchhans/procurement-system/internal/database"
	"github.com/synchhans/procurement-system/internal/models"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	database.ConnectDB()

	suppliers := []models.Supplier{
		{Name: "Global Tech", Email: "contact@globaltech.com", Address: "Jakarta"},
		{Name: "Indo Mandiri", Email: "sales@indomandiri.com", Address: "Bandung"},
	}

	for _, s := range suppliers {
		database.DB.FirstOrCreate(&s, models.Supplier{Email: s.Email})
	}

	items := []models.Item{
		{Name: "MacBook Pro", Stock: 10, Price: 25000000},
		{Name: "Dell XPS 13", Stock: 5, Price: 18000000},
		{Name: "Logitech Mouse", Stock: 50, Price: 500000},
	}

	for _, it := range items {
		database.DB.FirstOrCreate(&it, models.Item{Name: it.Name})
	}

	log.Println("Database seeded successfully!")
}
