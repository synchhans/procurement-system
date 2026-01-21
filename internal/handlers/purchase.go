package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
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
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	database.DB.Preload("User").Preload("Supplier").Preload("Details.Item").First(&purchasing, purchasing.ID)

	go sendWebhook(purchasing)

	return c.Status(201).JSON(purchasing)
}

func sendWebhook(data interface{}) {
	webhookURL := os.Getenv("WEBHOOK_URL")

	fmt.Println("\n--- WEBHOOK DEBUG START ---")
	fmt.Println("Target URL:", webhookURL)

	if webhookURL == "" {
		fmt.Println("ERROR: WEBHOOK_URL is empty or not loaded from .env")
		fmt.Println("--- WEBHOOK DEBUG END ---")
		return
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("ERROR JSON Marshal:", err)
		return
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("ERROR Creating Request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR Sending Request:", err)
		fmt.Println("--- WEBHOOK DEBUG END ---")
		return
	}
	defer resp.Body.Close()

	fmt.Println("Webhook Response Status:", resp.Status)
	fmt.Println("--- WEBHOOK DEBUG END ---")
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
