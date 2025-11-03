package utils

import (
	"github.com/ojiehdavid5/campusbyte/config"
	"github.com/ojiehdavid5/campusbyte/model"
	"log"
)

func SeedMenu() {
	menuItems := []model.MenuItem{
		{Name: "Jollof Rice", Description: "Served with chicken & plantain", Price: 1500},
		{Name: "Burger & Fries", Description: "Double beef burger with crispy fries", Price: 2500},
		{Name: "Shawarma", Description: "Chicken shawarma with extra cream", Price: 2000},
		{Name: "Fried Rice", Description: "With spicy turkey wings", Price: 1800},
	}

	for _, item := range menuItems {
		var existing model.MenuItem
		if err := config.DB.Where("name = ?", item.Name).First(&existing).Error; err != nil {
			config.DB.Create(&item)
			log.Printf("🍴 Added menu item: %s", item.Name)
		}
	}
}
