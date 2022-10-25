package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	ID        int `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time
	Username  string `json:"username"`
	Email     string `form:"email" json:email`
	Password  string `json:"password"`
}
