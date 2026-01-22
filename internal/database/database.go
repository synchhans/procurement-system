package database

import (
	"fmt"
	"log"
	"os"

	"github.com/synchhans/procurement-system/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Database connection successful")

	db.AutoMigrate(&models.User{}, &models.Supplier{}, &models.Item{}, &models.Purchasing{}, &models.PurchasingDetail{})

	DB = db

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("[FATAL] Gagal terhubung ke Database: %v", err)
	}

	log.Println("[SUKSES] Terhubung ke Database PostgreSQL.")

	log.Println("[INFO] Menjalankan migrasi database...")
	err = DB.AutoMigrate(
		&models.User{},
		&models.Supplier{},
		&models.Item{},
		&models.Purchasing{},
		&models.PurchasingDetail{},
	)

	if err != nil {
		log.Printf("[ERROR] Gagal migrasi: %v", err)
	} else {
		log.Println("[SUKSES] Struktur database sinkron.")
	}
}
