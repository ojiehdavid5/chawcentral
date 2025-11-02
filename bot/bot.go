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

func StartBot() *tgbotapi.BotAPI { // Original function name
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
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)
	go func() {

		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		updates := bot.GetUpdatesChan(u)

		for update := range updates {
			if update.Message != nil {
				// Handle incoming messages (text, commands)
				chatID := update.Message.Chat.ID
				text := update.Message.Text
				fmt.Println(chatID, text)
			}
		}
	}()

	return bot
}
func HandleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	user := update.Message.From
	telegramID := user.ID

	// Check if user already exists
	var existing model.User
	result := config.DB.Where("telegram_id = ?", telegramID).First(&existing)

	if result.Error != nil {
		// If not found, create new user
		newUser := model.User{
			TelegramID: telegramID,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			Username:   user.UserName,
		}
		config.DB.Create(&newUser)
		log.Printf("New user registered: %s (%d)\n", newUser.FirstName, newUser.TelegramID)
	}

	// Send a welcome message
	msgText := fmt.Sprintf("👋 Hey %s! Welcome to CampusBite — your tradefair food assistant!", user.FirstName)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	bot.Send(msg)
}
