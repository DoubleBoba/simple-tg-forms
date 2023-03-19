package function

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/sourcegraph/conc"
	tgbotapi "gitlab.com/kingofsystem/telegram-bot-api/v5"
)

func init() {
	functions.HTTP("process-tg-update", ProcessTelegramUpdate)
}

type ContextKey string

var wg_ctx_key = ContextKey("wg_ctx_key")

// Entry function that called when telegram hook is triggered.
func ProcessTelegramUpdate(response_writer http.ResponseWriter, request *http.Request) {
	wg := conc.NewWaitGroup()
	ctx := context.WithValue(context.Background(), wg_ctx_key, wg)

	// Asynchronyously obtain database connection, and don't block further execution.
	// If you will need connection at some point in the future, you can take it from the channel.
	firestore_client_future := make(chan *firestore.Client)
	wg.Go(func() {
		firestore_client_future <- createFirestoreClient(ctx)
	})

	// Do some other work.
	// check that secret token matches one from from enviroment variable
	header_value, ok := request.Header["X-Telegram-Bot-Api-Secret-Token"]
	if !ok {
		panic("Secret token not found")
	}
	if header_value[0] != os.Getenv("WEBHOOK_SECRET_TOKEN") {
		panic("Secret token does not match")
	}

	var update tgbotapi.Update
	if err := json.NewDecoder(request.Body).Decode(&update); err != nil {
		panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		panic(err)
	}
	bot.Debug = true

	// At this point i found out that i actually want to use database connection and i can obtain it.
	firestore_client := <-firestore_client_future
	defer firestore_client.Close()

	router := NewRouter(bot, firestore_client)
	router.AddHandler("/start", StartCommand)
	router.AddHandler("/ping", PingCommand)

	wg.Go(func() {
		save_update(ctx, firestore_client, &update)
	})

	err = router.HandleUpdate(ctx, &update)
	if err != nil {
		panic(err)
	}
	fmt.Println("Update handled!")
	wg.Wait()
}
