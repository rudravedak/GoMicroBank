package model

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	CustomerID  uint    `gorm:"not null" json:"customer_id"`
	CardID      uint    `json:"card_id"` // Optional, for card payments
	Amount      float64 `gorm:"not null" json:"amount"`
	PaymentType string  `gorm:"size:10;not null" json:"payment_type"` // "CARD" or "CASH"
	Status      string  `gorm:"size:20;not null" json:"status"`       // "PENDING", "PROCESSING", "COMPLETED", "FAILED", "CANCELLED"
	Description string  `json:"description"`
}
