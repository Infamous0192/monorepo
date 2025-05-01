package bot

// import (
// 	"context"
// 	"fmt"
// 	"app/pkg/telegram/service/bot"

// 	"gopkg.in/telebot.v4"
// )

// // BotHandlerSetup defines a function type for setting up bot-specific handlers
// type BotHandlerSetup func(bot bot.BotService)

// // BotHandlerRegistry stores handler setups for different bot types
// var BotHandlerRegistry = map[string]BotHandlerSetup{
// 	"default": setupDefaultBot,
// 	"support": setupSupportBot,
// 	"news":    setupNewsBot,
// 	"payment": setupPaymentBot,
// 	// Add more bot types here
// }

// // SetupBotByType sets up handlers for a specific bot type
// func SetupBotByType(botType string, bot bot.BotService) error {
// 	setup, exists := BotHandlerRegistry[botType]
// 	if !exists {
// 		return fmt.Errorf("unknown bot type: %s", botType)
// 	}

// 	setup(bot)
// 	return nil
// }

// // setupDefaultBot sets up handlers for a default bot
// func setupDefaultBot(bot bot.BotService) {
// 	bot.Command("start", func(ctx context.Context, m *telebot.Message) error {
// 		welcomeText := fmt.Sprintf(
// 			"👋 Hello %s!\n\nI'm a default bot. Here's what I can do:\n\n"+
// 				"/help - Show available commands",
// 			m.Sender.FirstName,
// 		)
// 		return bot.SendMessage(ctx, m.Chat.ID, welcomeText)
// 	})

// 	bot.Command("help", func(ctx context.Context, m *telebot.Message) error {
// 		return bot.SendMessage(ctx, m.Chat.ID, "I'm a default bot with basic commands.")
// 	})
// }

// // setupSupportBot sets up handlers for a support bot
// func setupSupportBot(bot bot.BotService) {
// 	bot.Command("start", func(ctx context.Context, m *telebot.Message) error {
// 		welcomeText := fmt.Sprintf(
// 			"👋 Welcome to Support, %s!\n\n"+
// 				"How can I help you today?\n\n"+
// 				"/ticket - Create a support ticket\n"+
// 				"/faq - View frequently asked questions\n"+
// 				"/contact - Contact support team",
// 			m.Sender.FirstName,
// 		)
// 		return bot.SendMessage(ctx, m.Chat.ID, welcomeText)
// 	})

// 	bot.Command("ticket", func(ctx context.Context, m *telebot.Message) error {
// 		return bot.SendMessage(ctx, m.Chat.ID, "Creating a new support ticket...")
// 	})

// 	bot.Command("faq", func(ctx context.Context, m *telebot.Message) error {
// 		return bot.SendMessage(ctx, m.Chat.ID, "Here are our frequently asked questions...")
// 	})

// 	bot.Handle(func(ctx context.Context, m *telebot.Message) error {
// 		return bot.SendMessage(ctx, m.Chat.ID, "Support team will get back to you soon!")
// 	})
// }

// // setupNewsBot sets up handlers for a news bot
// func setupNewsBot(bot bot.BotService) {
// 	bot.Command("start", func(ctx context.Context, m *telebot.Message) error {
// 		welcomeText := fmt.Sprintf(
// 			"📰 Welcome to NewsBot, %s!\n\n"+
// 				"Stay updated with the latest news:\n\n"+
// 				"/subscribe - Subscribe to news categories\n"+
// 				"/latest - Get latest news\n"+
// 				"/categories - View available categories",
// 			m.Sender.FirstName,
// 		)
// 		return bot.SendMessage(ctx, m.Chat.ID, welcomeText)
// 	})

// 	bot.Command("subscribe", func(ctx context.Context, m *telebot.Message) error {
// 		menu := &telebot.ReplyMarkup{ResizeKeyboard: true}
// 		menu.Inline(
// 			menu.Row(
// 				telebot.Btn{Text: "🌍 World", Unique: "sub_world"},
// 				telebot.Btn{Text: "💼 Business", Unique: "sub_business"},
// 			),
// 			menu.Row(
// 				telebot.Btn{Text: "🏃 Sports", Unique: "sub_sports"},
// 				telebot.Btn{Text: "🎬 Entertainment", Unique: "sub_entertainment"},
// 			),
// 		)

// 		_, err := bot.GetBot().Send(m.Chat, "Choose categories to subscribe:", menu)
// 		return err
// 	})

// 	bot.Action("sub_world", func(ctx context.Context, c *telebot.Callback) error {
// 		return bot.SendMessage(ctx, c.Message.Chat.ID, "Subscribed to World news!")
// 	})

// 	bot.Handle(func(ctx context.Context, m *telebot.Message) error {
// 		return bot.SendMessage(ctx, m.Chat.ID, "Use commands to interact with NewsBot!")
// 	})
// }

// // setupPaymentBot sets up handlers for a payment bot
// func setupPaymentBot(bot bot.BotService) {
// 	bot.Command("start", func(ctx context.Context, m *telebot.Message) error {
// 		welcomeText := fmt.Sprintf(
// 			"💰 Welcome to PaymentBot, %s!\n\n"+
// 				"I can help you with payments and invoices:\n\n"+
// 				"/invoice - Create a new invoice\n"+
// 				"/payments - View your payment history\n"+
// 				"/help - Get help with payments",
// 			m.Sender.FirstName,
// 		)
// 		return bot.SendMessage(ctx, m.Chat.ID, welcomeText)
// 	})

// 	bot.Command("invoice", func(ctx context.Context, m *telebot.Message) error {
// 		return bot.SendMessage(ctx, m.Chat.ID, "To create an invoice, please specify the amount and description.")
// 	})

// 	bot.Command("payments", func(ctx context.Context, m *telebot.Message) error {
// 		return bot.SendMessage(ctx, m.Chat.ID, "Your payment history will be displayed here.")
// 	})

// 	// Handle pre-checkout queries
// 	bot.GetBot().Handle(&telebot.PreCheckoutQuery{}, func(c telebot.Context) error {
// 		// Accept all pre-checkout queries for now
// 		// In a real implementation, you would validate the payment details
// 		return c.Accept()
// 	})

// 	// Handle successful payments
// 	bot.GetBot().Handle(&telebot.Payment{}, func(c telebot.Context) error {
// 		msg := c.Message()

// 		// Log the successful payment
// 		fmt.Printf("Received payment: %+v\n", msg.Payment)

// 		// In a real implementation, you would process the payment
// 		// using a payment service
// 		return c.Send("Thank you for your payment!")
// 	})

// 	bot.Handle(func(ctx context.Context, m *telebot.Message) error {
// 		return bot.SendMessage(ctx, m.Chat.ID, "Use commands to interact with PaymentBot!")
// 	})
// }
