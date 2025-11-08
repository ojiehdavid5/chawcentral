package model

import "gorm.io/gorm"

type Wallet struct {
    gorm.Model
    WalletID   uint    `json:"wallet_id"`
    Balance  float64 `json:"balance"`
}
