package bot

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)
import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
				fmt.Println(chatID,text)
			}
		}
	}()

	return bot
}
