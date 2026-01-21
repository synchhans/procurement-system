package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/synchhans/procurement-system/internal/database"
	"github.com/synchhans/procurement-system/internal/models"
	"gorm.io/gorm"
)

type PurchaseInput struct {
	SupplierID uint `json:"supplier_id"`
	Items      []struct {
		ItemID uint `json:"item_id"`
		Qty    int  `json:"qty"`
	} `json:"items"`
}

func CreatePurchase(c *fiber.Ctx) error {
	var input PurchaseInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	userID := c.Locals("user_id").(uint)

	var grandTotal float64
	var purchasing models.Purchasing

	log.Printf("[ORDER] Memulai transaksi baru oleh User ID: %d", userID)

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		purchasing.SupplierID = input.SupplierID
		purchasing.UserID = userID

		if err := tx.Create(&purchasing).Error; err != nil {
			return err
		}

		for _, itemInput := range input.Items {
			var item models.Item
			if err := tx.First(&item, itemInput.ItemID).Error; err != nil {
				return fmt.Errorf("item with ID %d not found", itemInput.ItemID)
			}

			subTotal := item.Price * float64(itemInput.Qty)

			if subTotal > 100000000000000 {
				tx.Rollback()
				return c.Status(400).JSON(fiber.Map{
					"error": fmt.Sprintf("Total amount for item %s exceeds limit (Max 100 Trillion)", item.Name),
				})
			}
			grandTotal += subTotal

			detail := models.PurchasingDetail{
				PurchasingID: purchasing.ID,
				ItemID:       item.ID,
				Qty:          itemInput.Qty,
				SubTotal:     subTotal,
			}
			if err := tx.Create(&detail).Error; err != nil {
				return err
			}

			if err := tx.Model(&item).Update("stock", gorm.Expr("stock + ?", itemInput.Qty)).Error; err != nil {
				return err
			}
		}

		if err := tx.Model(&purchasing).Update("grand_total", grandTotal).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Printf("[ERROR] Transaksi Gagal & Rollback: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("[SUKSES] Transaksi ID %s berhasil disimpan. Total: %.2f", purchasing.ID, purchasing.GrandTotal)

	database.DB.Preload("User").Preload("Supplier").Preload("Details.Item").First(&purchasing, purchasing.ID)

	go sendWebhook(purchasing)

	return c.Status(201).JSON(purchasing)
}

func sendWebhook(data interface{}) {
	webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		log.Println("[WEBHOOK] URL kosong, melewati pengiriman webhook.")
		return
	}

	log.Printf("[WEBHOOK] Mengirim data ke: %s...", webhookURL)

	jsonData, _ := json.Marshal(data)
	client := &http.Client{Timeout: 30 * time.Second}

	req, _ := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[WEBHOOK] Gagal mengirim: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("[WEBHOOK] Terkirim! Status Server: %s", resp.Status)
}

func UpdateItem(c *fiber.Ctx) error {
	id := c.Params("id")
	var item models.Item

	if err := database.DB.First(&item, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Item not found"})
	}

	var input models.Item
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	database.DB.Model(&item).Updates(input)

	return c.JSON(item)
}

func DeleteItem(c *fiber.Ctx) error {
	id := c.Params("id")
	var item models.Item

	if err := database.DB.First(&item, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Item not found"})
	}

	database.DB.Delete(&item)

	return c.JSON(fiber.Map{"message": "Item deleted successfully"})
}

func UpdateSupplier(c *fiber.Ctx) error {
	id := c.Params("id")
	var supplier models.Supplier

	if err := database.DB.First(&supplier, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Supplier not found"})
	}

	var input models.Supplier
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	database.DB.Model(&supplier).Updates(input)

	return c.JSON(supplier)
}

func DeleteSupplier(c *fiber.Ctx) error {
	id := c.Params("id")
	var supplier models.Supplier

	if err := database.DB.First(&supplier, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Supplier not found"})
	}

	database.DB.Delete(&supplier)

	return c.JSON(fiber.Map{"message": "Supplier deleted successfully"})
}
