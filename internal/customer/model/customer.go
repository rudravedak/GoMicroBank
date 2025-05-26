package model

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	FirstName string `gorm:"size:100;not null" json:"first_name"`
	LastName  string `gorm:"size:100;not null" json:"last_name"`
	Email     string `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Phone     string `gorm:"size:20" json:"phone"`
	Address   string `gorm:"size:255" json:"address"`

	Balance float64  `gorm:"type:decimal(10,2);default:0" json:"balance"`
	Cards   []string `gorm:"type:text[]" json:"cards"` // Array of card numbers
}
