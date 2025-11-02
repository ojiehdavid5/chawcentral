package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/ojiehdavid5/campusbyte/bot"
	"github.com/ojiehdavid5/campusbyte/config"
)

func main() {
	app := fiber.New()
	config.ConnectDB()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Start Telegram bot in background
	telegramBot := bot.StartBot()
	if telegramBot == nil {
		panic("❌ Failed to start bot")
	}
	fmt.Println("✅ Telegram bot started")

	app.Listen(":3000")
}
