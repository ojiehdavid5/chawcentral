package model

import "gorm.io/gorm"

type MenuItem struct {
	gorm.Model
	Name        string
	Description string
	Price       float64
	ImageURL    string // optional if you want images later
}
