package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name           string    `gorm:"size:255;not null" json:"name"`
	Email          string    `gorm:"size:255;not null;unique" json:"email"`
	Password       string    `gorm:"size:255;not null" json:"-"`
	LastLoginAt    time.Time `json:"last_login_at"`
	FailedLogins   int       `gorm:"default:0" json:"-"`
	LockedUntil    time.Time `json:"-"`
	RegistrationIP string    `gorm:"size:45" json:"-"`
	IsActive       bool      `gorm:"default:true" json:"is_active"`
} 