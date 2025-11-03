package model

import (
	"time"
)

// Each item in a user's cart
type Cart struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     uint      // Reference to User.ID
	MenuItemID uint      // Reference to MenuItem.ID
	Quantity   int       `gorm:"default:1"`
	CreatedAt  time.Time
	UpdatedAt  time.Time

	// Relations
	User     User
	MenuItem MenuItem
}
