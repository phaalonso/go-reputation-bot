package main

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	models "github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
	repository "github.com/phaalonso/go-reputation-bot/pkg/models"
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
	b.RegisterHandler(bot.HandlerTypeMessageText, "/reputation", bot.MatchTypeExact, reputationHanler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/global_reputation", bot.MatchTypeExact, totalReputationHandler)

	b.DeleteMyCommands(ctx, &bot.DeleteMyCommandsParams{})
	//b.SetMyCommands(ctx, &bot.SetMyCommandsParams{
	//	Scope:        nil,
	//	LanguageCode: "",
	//})

	b.Start(ctx)
}

func startCommandHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Welcome to the reputation bot!\nAdd me into a chat so I can begin to track the reputation of your members",
	})
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message.ReplyToMessage == nil {
		return
	}

	fmt.Println(update.Message.Text)

	chatId := update.Message.Chat.ID
	userId := update.Message.ReplyToMessage.From.ID

	rep := repository.UpdateOrCreateReputation(chatId, userId)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   fmt.Sprintf("Esse usuário possui %d de reputação", rep.Reputation),
	})
}

func reputationHanler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID

	uReputation, err := repository.GetUserReputationInChat(chatId, userId)

	if err != nil {
		log.Fatal(err)
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("Your reputation in this chat is %d", uReputation.Reputation),
	})
}

func totalReputationHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	userId := update.Message.From.ID

	reputation, err := repository.GetTotalUserReputation(userId)

	if err != nil {
		log.Fatal(err)
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("Your global reputation is %d", reputation),
	})
}
