package model

import (
	"time"

	"gorm.io/gorm"
)

type Card struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	CustomerID  uint    `gorm:"not null" json:"customer_id"`
	CardNumber  string  `gorm:"size:16;not null;uniqueIndex" json:"card_number"`
	CardType    string  `gorm:"size:20;not null" json:"card_type"`
	ExpiryDate  string  `gorm:"size:5;not null" json:"expiry_date"`
	CVV         string  `gorm:"size:3;not null" json:"cvv"`
	CreditLimit float64 `gorm:"not null" json:"credit_limit"`
	Balance     float64 `gorm:"not null" json:"balance"`
	IsActive    bool    `gorm:"not null;default:true" json:"is_active"`
}
