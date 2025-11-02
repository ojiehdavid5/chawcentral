package bot

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/ojiehdavid5/campusbyte/model"
	"github.com/ojiehdavid5/campusbyte/config"
)

// StartBot launches the Telegram bot and handles updates
func StartBot() *tgbotapi.BotAPI {
	_ = godotenv.Load()
	token := os.Getenv("TELEGRAM_APITOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_APITOKEN is not set")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	bot.Debug = true
	fmt.Printf("✅ Authorized on account %s\n", bot.Self.UserName)

	// Run updates in a goroutine so it doesn't block Fiber
	go func() {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		updates := bot.GetUpdatesChan(u)

		for update := range updates {
			if update.Message == nil {
				continue
			}

			if update.Message.IsCommand() && update.Message.Command() == "start" {
				HandleStartCommand(bot, update)
			}
		}
	}()

	return bot
}

// HandleStartCommand handles /start — save user & greet them
func HandleStartCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	user := update.Message.From
	telegramID := user.ID

	var existing model.User
	result := config.DB.Where("telegram_id = ?", telegramID).First(&existing)

var msgText string // declare message text outside if/else

if result.Error != nil {
	// Create a new user if not found
	newUser := model.User{
		TelegramID: telegramID,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Username:   user.UserName,
	}
	config.DB.Create(&newUser)
	log.Printf("🆕 New user registered: %s (%d)\n", newUser.FirstName, newUser.TelegramID)

	// First-time welcome message
	msgText = fmt.Sprintf("👋 Hey %s! Welcome to CampusBite — your tradefair food assistant!", user.FirstName)
} else {
	// Returning user message
	msgText = fmt.Sprintf("👋 Welcome back, %s! Ready to order something delicious?", user.FirstName)
}



// Create inline buttons
menuButton := tgbotapi.NewInlineKeyboardButtonData("🍔 View Menu", "view_menu")
topUpButton := tgbotapi.NewInlineKeyboardButtonData("💳 Top Up Wallet", "top_up")
cartButton := tgbotapi.NewInlineKeyboardButtonData("🛒 View Cart", "view_cart")

// Arrange them in rows
keyboard := tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(menuButton, topUpButton, cartButton),
)

// Attach keyboard to message
msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
msg.ReplyMarkup = keyboard

// Send message

// Send the message

bot.Send(msg)
}
