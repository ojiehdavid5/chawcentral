package model

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model
	UserID        uint `gorm:"uniqueIndex"`
	User          User
	Balance       float64
	AccountNumber string // From Kora
	AccountName   string
	BankName      string
	KoraRef       string // Reference from Kora API
}


