package bot

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"

	"github.com/ojiehdavid5/campusbyte/config"
	"github.com/ojiehdavid5/campusbyte/model"
	"github.com/ojiehdavid5/campusbyte/kora"
)

// StartBot launches the Telegram bot and handles updates

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

	go func() {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		updates := bot.GetUpdatesChan(u)

		for update := range updates {

			// 🟩 Handle /start
			if update.Message != nil {
				if update.Message.IsCommand() && update.Message.Command() == "start" {
					HandleStartCommand(bot, update)
				}
			}

			// 🟩 Handle callback queries (buttons)
			if update.CallbackQuery != nil {
				data := update.CallbackQuery.Data
				chatID := update.CallbackQuery.Message.Chat.ID

				// ✅ 1️⃣ Handle normal static buttons
				switch data {
				case "view_menu":
					var menuItems []model.MenuItem
					result := config.DB.Find(&menuItems)

					if result.Error != nil || len(menuItems) == 0 {
						bot.Send(tgbotapi.NewMessage(chatID, "😕 No menu items available right now. Check back later!"))
						continue
					}

					for _, item := range menuItems {
						text := fmt.Sprintf("🍴 *%s*\n💰 ₦%.2f\n%s", item.Name, item.Price, item.Description)

						// 👇 note: we use "add_cart_" here
						addToCartBtn := tgbotapi.NewInlineKeyboardButtonData("➕ Add to Cart", fmt.Sprintf("add_cart_%d", item.ID))
						keyboard := tgbotapi.NewInlineKeyboardMarkup(
							tgbotapi.NewInlineKeyboardRow(addToCartBtn),
						)

						msg := tgbotapi.NewMessage(chatID, text)
						msg.ParseMode = "Markdown"
						msg.ReplyMarkup = keyboard
						bot.Send(msg)
					}

				case "top_up": 
					amount := 1000.0
					userID := update.CallbackQuery.From.ID

					// Find the user in DB
					var user model.User
					if err := config.DB.Where("telegram_id = ?", userID).First(&user).Error; err != nil {
						bot.Send(tgbotapi.NewMessage(chatID, "⚠️ Could not find your account. Try /start again."))
						continue
					}

					ref := fmt.Sprintf("wallet_topup_%d_%d", user.ID, time.Now().Unix())
					fmt.Println("This is the reference ", ref)

checkoutURL, err := kora.CreateKoraPayment(user, amount, "wallet")
fmt.Println("error ", err)
								if err != nil {
						bot.Send(tgbotapi.NewMessage(chatID, "⚠️ Could not initiate payment. Try again."))
						return
					}
					fmt.Println("This is the url ", checkoutURL)
					msg := fmt.Sprintf("💳 Tap below to top-up ₦%.2f:\n%s", amount, checkoutURL)
					bot.Send(tgbotapi.NewMessage(chatID, msg))
				case "clear_cart":
					userID := update.CallbackQuery.From.ID

					// Find the user
					var user model.User
					if err := config.DB.Where("telegram_id = ?", userID).First(&user).Error; err != nil {
						bot.Send(tgbotapi.NewMessage(chatID, "⚠️ Could not find your account. Try /start again."))
						return
					}

					// Delete all their cart items
					if err := config.DB.Where("user_id = ?", user.ID).Delete(&model.Cart{}).Error; err != nil {
						bot.Send(tgbotapi.NewMessage(chatID, "❌ Failed to clear your cart. Please try again."))
						return
					}

					bot.Send(tgbotapi.NewMessage(chatID, "🧹 Your cart has been cleared successfully!"))

				case "view_cart":
					userID := update.CallbackQuery.From.ID

					// Find user in DB
					var user model.User
					config.DB.Where("telegram_id = ?", userID).First(&user)

					// Fetch user's cart items
					var cartItems []model.Cart
					config.DB.Preload("MenuItem").Where("user_id = ?", user.ID).Find(&cartItems)

					if len(cartItems) == 0 {
						bot.Send(tgbotapi.NewMessage(chatID, "🛒 Your cart is empty. Add something delicious from the menu!"))
						continue
					}

					total := 0.0
					cartText := "🛒 *Your Cart:*\n\n"

					for _, c := range cartItems {
						itemTotal := c.MenuItem.Price * float64(c.Quantity)
						total += itemTotal
						cartText += fmt.Sprintf("🍴 %s x%d — ₦%.2f\n", c.MenuItem.Name, c.Quantity, itemTotal)
					}

					cartText += fmt.Sprintf("\n💰 *Total:* ₦%.2f", total)

					// Add checkout buttons
					checkoutBtn := tgbotapi.NewInlineKeyboardButtonData("🧾 Checkout", "checkout")
					clearBtn := tgbotapi.NewInlineKeyboardButtonData("❌ Clear Cart", "clear_cart")

					keyboard := tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(checkoutBtn, clearBtn),
					)

					msg := tgbotapi.NewMessage(chatID, cartText)
					msg.ParseMode = "Markdown"
					msg.ReplyMarkup = keyboard
					bot.Send(msg)
				} // <- closes switch

				// ✅ 2️⃣ Handle dynamic "add_cart_" actions (Step 4) — moved outside switch
				if strings.HasPrefix(data, "add_cart_") {
					itemIDStr := strings.TrimPrefix(data, "add_cart_")
					itemID, _ := strconv.Atoi(itemIDStr)
					userID := update.CallbackQuery.From.ID

					var user model.User
					config.DB.Where("telegram_id = ?", userID).First(&user)

					var cartItem model.Cart
					err := config.DB.Where("user_id = ? AND menu_item_id = ?", user.ID, itemID).First(&cartItem).Error
					if err != nil {
						// New cart entry
						cartItem = model.Cart{
							UserID:     user.ID,
							MenuItemID: uint(itemID),
							Quantity:   1,
						}
						config.DB.Create(&cartItem)
						bot.Send(tgbotapi.NewMessage(chatID, "✅ Added to cart!"))
					} else {
						// Update quantity
						cartItem.Quantity++
						config.DB.Save(&cartItem)
						bot.Send(tgbotapi.NewMessage(chatID, "🛒 Quantity updated!"))
					}
				}
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

	var msgText string

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

		msgText = fmt.Sprintf("👋 Hey %s! Welcome to CampusBite — your tradefair food assistant!", user.FirstName)
	} else {
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

	// Send the message
	bot.Send(msg)
}
