package function

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	tgbotapi "gitlab.com/kingofsystem/telegram-bot-api/v5"
)

func StartCommand(ctx context.Context, update *tgbotapi.Update, bot *tgbotapi.BotAPI, db_client *firestore.Client) error {

	get_user_channel := GetUser(ctx, db_client, update.Message.From.ID)
	user := <-get_user_channel

	fmt.Printf("User: %v\n", user)

	ReadChannel(
		ctx,
		CreateOrUpdateUser(ctx, db_client, update.Message.From),
	)

	if user == nil {
		ReadChannel(
			ctx,
			StartFormFill(ctx, update.Message.From, bot, db_client),
		)
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
) chan bool {
	done_channel := make(chan bool)

	go func() {
		welcome_message := tgbotapi.NewMessage(user.ID, "Welcome!")
		bot.Send(welcome_message)

		form := <-AttachNewFormToUser(ctx, db_client, user.ID)
		ReadChannel(
			ctx,
			SendFormQuestion(ctx, bot, form.GetCurrentQuestion(), user),
		)
		done_channel <- true
	}()

	return done_channel
}

func SendFormQuestion(ctx context.Context, bot *tgbotapi.BotAPI, question *Question, user *tgbotapi.User) chan bool {
	done_channel := make(chan bool)

	go func() {
		message := tgbotapi.NewMessage(user.ID, question.QuestionText)
		fmt.Println("Sending question...")
		_, err := bot.Send(message)
		if err != nil {
			panic(err)
		}
		done_channel <- true
	}()

	return done_channel
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
) chan bool {
	done_channel := make(chan bool)

	go func() {
		form := <-GetUserForm(ctx, db_client, update.Message.From.ID)
		if form == nil {
			fmt.Printf("Form for user %v not found\n", update.Message.From.ID)
			done_channel <- false
			return
		}
		fmt.Printf("Form for user with id=%v is found\n", update.Message.From.ID)

		current_question := form.GetCurrentQuestion()
		current_question.Answer = update.Message.Text

		if form.IsLastQuestion() {
			form.IsFilled = true
		} else {
			form.CurrentQuestionIndex += 1
			defer ReadChannel(
				ctx,
				SendFormQuestion(ctx, bot, form.GetCurrentQuestion(), update.Message.From),
			)
		}

		ReadChannel(
			ctx,
			UpdateForm(ctx, db_client, form),
		)

		done_channel <- true
	}()

	return done_channel
}
