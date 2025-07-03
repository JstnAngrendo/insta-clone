package entities

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

type AccessToken struct {
	ID        string `gorm:"primaryKey"`
	UserID    uint
	CreatedAt time.Time
	ExpiresAt time.Time
	Revoked   bool
}
