package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"-"`
	Role     string `gorm:"default:staff" json:"role"`
}

type Supplier struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Name    string `gorm:"not null" json:"name"`
	Email   string `gorm:"unique" json:"email"`
	Address string `json:"address"`
}

type Item struct {
	ID    uint    `gorm:"primaryKey" json:"id"`
	Name  string  `gorm:"not null" json:"name"`
	Stock int     `gorm:"default:0" json:"stock"`
	Price float64 `gorm:"type:decimal(20,2);not null" json:"price"`
}

type Purchasing struct {
	ID         uuid.UUID          `gorm:"type:uuid;primaryKey" json:"id"`
	Date       time.Time          `json:"date"`
	SupplierID uint               `json:"supplier_id"`
	Supplier   Supplier           `gorm:"foreignKey:SupplierID" json:"supplier"`
	UserID     uint               `json:"user_id"`
	User       User               `gorm:"foreignKey:UserID" json:"user"`
	GrandTotal float64            `gorm:"type:decimal(20,2)" json:"grand_total"`
	Details    []PurchasingDetail `gorm:"foreignKey:PurchasingID" json:"details"`
}

type PurchasingDetail struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	PurchasingID uuid.UUID `gorm:"type:uuid" json:"purchasing_id"`
	ItemID       uint      `json:"item_id"`
	Item         Item      `gorm:"foreignKey:ItemID" json:"item"`
	Qty          int       `json:"qty"`
	SubTotal     float64   `gorm:"type:decimal(20,2)" json:"sub_total"`
}

func (p *Purchasing) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	p.Date = time.Now()
	return
}
