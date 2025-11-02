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
	}

	// Send personalized welcome message
	msgText := fmt.Sprintf("👋 Hey %s! Welcome to CampusBite — your tradefair food assistant!", user.FirstName)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	bot.Send(msg)
}
