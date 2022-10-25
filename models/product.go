package models

import (
	"gorm.io/gorm"
	"time"
)

type Product struct {
	gorm.Model
	ID        uint `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time
	Name      string  `json:"name"`
	Quantity  int     `json:quantity`
	Price     float64 `json:price`
	Image     string  `json:image`
}
