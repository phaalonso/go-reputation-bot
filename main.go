package main

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	telegramKey := os.Getenv("TELEGRAM_API_KEY")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
		bot.WithDebug(),
	}

	b, err := bot.New(telegramKey, opts...)

	if err != nil {
		log.Fatal(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startCommandHandler)

	b.Start(ctx)
}

func startCommandHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Welcome to the reputation bot!\nAdd me into a chat so I can begin to track the reputation of your members",
	})
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	fmt.Println(update)

	if update.MessageReaction == nil || len(update.MessageReaction.NewReaction) == 0 {
		return
	}

	fmt.Println(update)
}
