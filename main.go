package main

import (
	// "fmt"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/ojiehdavid5/campusbyte/bot"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// ✅ Start Telegram bot
	telegramBot := bot.StartBot()
	if telegramBot == nil {
		panic("❌ Failed to start bot")
	}
	fmt.Println("✅ Telegram bot started")

	app.Listen(":3000")
}
