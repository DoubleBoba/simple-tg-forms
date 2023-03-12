package function

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	tgbotapi "gitlab.com/kingofsystem/telegram-bot-api/v5"
)

func init() {
	functions.HTTP("process-tg-update", ProcessTelegramUpdate)
}

func ProcessTelegramUpdate(response_writer http.ResponseWriter, request *http.Request) {
	ctx := context.Background()

	create_client_channel := createFirestoreClient(ctx)

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

	firestore_client := <-create_client_channel
	defer firestore_client.Close()

	router := NewRouter(bot, firestore_client)
	router.AddHandler("/start", StartCommand)
	router.AddHandler("/ping", PingCommand)

	save_update_channel := save_update(ctx, firestore_client, &update)
	defer ReadChannel(ctx, save_update_channel)

	err = router.HandleUpdate(ctx, &update)
	if err != nil {
		panic(err)
	}
	fmt.Println("Update handled!")
}
