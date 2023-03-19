package function

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/sourcegraph/conc"
	tgbotapi "gitlab.com/kingofsystem/telegram-bot-api/v5"
)

func StartCommand(ctx context.Context, update *tgbotapi.Update, bot *tgbotapi.BotAPI, db_client *firestore.Client) error {
	wg := ctx.Value(wg_ctx_key).(*conc.WaitGroup)

	user, err := GetUser(ctx, db_client, update.Message.From.ID)
	if err != nil {
		panic(err)
	}

	fmt.Printf("User: %v\n", user)

	wg.Go(func() {
		CreateOrUpdateUser(ctx, db_client, update.Message.From)
	})

	if user == nil {
		wg.Go(func() {
			StartFormFill(ctx, update.Message.From, bot, db_client)
		})
	} else {
		reply := tgbotapi.NewMessage(update.Message.Chat.ID, "I know you already.")
		reply.ReplyToMessageID = update.Message.MessageID
		bot.Send(reply)
	}

	return nil
}

func StartFormFill(
	ctx context.Context,
	user *tgbotapi.User,
	bot *tgbotapi.BotAPI,
	db_client *firestore.Client,
) {
	welcome_message := tgbotapi.NewMessage(user.ID, "Welcome!")
	bot.Send(welcome_message)
	form := AttachNewFormToUser(ctx, db_client, user.ID)
	SendFormQuestion(ctx, bot, form.GetCurrentQuestion(), user)
}

func SendFormQuestion(ctx context.Context, bot *tgbotapi.BotAPI, question *Question, user *tgbotapi.User) {
	message := tgbotapi.NewMessage(user.ID, question.QuestionText)
	fmt.Println("Sending question...")
	_, err := bot.Send(message)
	if err != nil {
		panic(err)
	}
}

func PingCommand(ctx context.Context, update *tgbotapi.Update, bot *tgbotapi.BotAPI, db_client *firestore.Client) error {
	reply := tgbotapi.NewMessage(update.Message.Chat.ID, "Pong!")
	reply.ReplyToMessageID = update.Message.MessageID
	bot.Send(reply)

	return nil
}

func HandleQuestionAnswer(
	ctx context.Context,
	update *tgbotapi.Update,
	bot *tgbotapi.BotAPI,
	db_client *firestore.Client,
) {
	wg := ctx.Value(wg_ctx_key).(*conc.WaitGroup)

	form := GetUserForm(ctx, db_client, update.Message.From.ID)
	if form == nil {
		fmt.Printf("Form for user %v not found\n", update.Message.From.ID)
		return
	}
	fmt.Printf("Form for user with id=%v is found\n", update.Message.From.ID)

	current_question := form.GetCurrentQuestion()
	current_question.Answer = update.Message.Text

	if form.IsLastQuestion() {
		form.IsFilled = true
	} else {
		form.CurrentQuestionIndex += 1

		wg.Go(func() {
			SendFormQuestion(ctx, bot, form.GetCurrentQuestion(), update.Message.From)
		})
	}
	wg.Go(func() {
		UpdateForm(ctx, db_client, form)
	})
}
