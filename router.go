package function

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/firestore"
	tgbotapi "gitlab.com/kingofsystem/telegram-bot-api/v5"
)

type CommandName string
type CommandHandler func(context.Context, *tgbotapi.Update, *tgbotapi.BotAPI, *firestore.Client) error

type Router struct {
	handlers  map[CommandName]CommandHandler
	bot       *tgbotapi.BotAPI
	db_client *firestore.Client
}

func NewRouter(bot *tgbotapi.BotAPI, db_client *firestore.Client) *Router {
	return &Router{
		bot:       bot,
		db_client: db_client,
		handlers:  make(map[CommandName]CommandHandler, 16),
	}
}

func (router *Router) AddHandler(command CommandName, handler CommandHandler) error {
	if _, ok := router.handlers[command]; ok {
		return fmt.Errorf("command %s already exists", command)
	}
	router.handlers[command] = handler
	return nil
}

func (router *Router) HandleUpdate(ctx context.Context, update *tgbotapi.Update) error {
	if update.Message != nil && update.Message.IsCommand() {
		return router.HandleCommand(ctx, update)
	} else if update.Message != nil && update.Message.Text != "" {
		HandleQuestionAnswer(ctx, update, router.bot, router.db_client)
	} else {
		return fmt.Errorf("upddate with id %v is unhandled", update.UpdateID)
	}

	return nil
}

func (router *Router) HandleCommand(ctx context.Context, update *tgbotapi.Update) error {
	command_name := strings.Split(update.Message.Text, " ")[0]

	command, ok := router.handlers[CommandName(command_name)]
	if !ok {
		return fmt.Errorf("command %s not found", command_name)
	}
	command(ctx, update, router.bot, router.db_client)
	return nil
}
