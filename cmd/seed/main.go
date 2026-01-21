package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/joho/godotenv"
	"github.com/synchhans/procurement-system/internal/database"
	"github.com/synchhans/procurement-system/internal/models"
	"github.com/synchhans/procurement-system/pkg/utils"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, relying on system env")
	}

	database.ConnectDB()
	log.Println("--- STARTING SEEDING ---")

	password := "password123"
	hashedPassword, _ := utils.HashPassword(password)

	defaultUser := models.User{
		Username: "admin",
		Password: hashedPassword,
		Role:     "staff",
	}

	if err := database.DB.Where("username = ?", defaultUser.Username).FirstOrCreate(&defaultUser).Error; err != nil {
		log.Printf("Error seeding user: %v\n", err)
	} else {
		log.Printf("✅ User Created: %s / %s\n", defaultUser.Username, password)
	}

	suppliers := []models.Supplier{
		{Name: "PT. Global Teknologi Data", Email: "sales@gtd.co.id", Address: "Jakarta Selatan, SCBD Lot 8"},
		{Name: "CV. Sinar Jaya Makmur", Email: "admin@sinarjaya.com", Address: "Bandung, Jl. Soekarno Hatta No 10"},
		{Name: "Mega Electronics Ltd", Email: "support@mega-elec.sg", Address: "Singapore, Changi Business Park"},
		{Name: "PT. Hardware Indonesia", Email: "contact@hardindo.co.id", Address: "Surabaya, Rungkut Industri"},
		{Name: "Distributor IT Pusat", Email: "info@itpusat.net", Address: "Jakarta Pusat, Mangga Dua Mall"},
	}

	for _, s := range suppliers {
		if err := database.DB.Where("email = ?", s.Email).FirstOrCreate(&s).Error; err != nil {
			log.Printf("Error seeding supplier %s: %v\n", s.Name, err)
		}
	}
	log.Printf("✅ Suppliers Created: %d Data\n", len(suppliers))

	brands := []string{"MacBook", "Dell", "HP", "Lenovo", "Asus", "Acer", "Samsung", "Logitech", "Sony"}
	types := []string{"Pro", "Air", "XPS", "ThinkPad", "ROG", "Pavilion", "Ultra", "MX Master", "Bravia"}
	categories := []string{"Laptop", "Monitor", "Mouse", "Keyboard", "Headset", "Webcam"}

	rand.Seed(time.Now().UnixNano())

	count := 0
	for i := 1; i <= 50; i++ {
		brand := brands[rand.Intn(len(brands))]
		typ := types[rand.Intn(len(types))]
		cat := categories[rand.Intn(len(categories))]

		itemName := fmt.Sprintf("%s %s %s - Gen %d", brand, typ, cat, i)

		price := float64((rand.Intn(500) + 1) * 50000)

		stock := rand.Intn(90) + 10

		item := models.Item{
			Name:  itemName,
			Stock: stock,
			Price: price,
		}

		if err := database.DB.Where("name = ?", item.Name).FirstOrCreate(&item).Error; err == nil {
			count++
		}
	}

	log.Printf("✅ Items Generated: %d Data (Total ideal untuk test scrolling)\n", count)
	log.Println("--- SEEDING FINISHED ---")
}
