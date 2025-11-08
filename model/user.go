package model


import (
	"time"
)

type User struct {
	ID         uint      `gorm:"primaryKey"`
	TelegramID int64     `gorm:"uniqueIndex"`
	FirstName  string
	LastName   string
	Username   string
	CreatedAt  time.Time
	Wallet Wallet `gorm:"foreignKey:WalletID"`

}