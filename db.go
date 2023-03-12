package function

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	tgbotapi "gitlab.com/kingofsystem/telegram-bot-api/v5"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func createFirestoreClient(ctx context.Context) chan *firestore.Client {
	result_channel := make(chan *firestore.Client)
	go func() {
		// Sets your Google Cloud Platform project ID.
		projectID := os.Getenv("PROJECT_ID")

		client, err := firestore.NewClient(ctx, projectID)
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
		// Close client when done with
		// defer client.Close()
		result_channel <- client
	}()
	return result_channel
}

func save_update(
	ctx context.Context,
	firestore_client *firestore.Client,
	update *tgbotapi.Update,
) chan bool {
	done_channel := make(chan bool)

	go func() {
		_, err := firestore_client.Collection("updates").Doc(fmt.Sprint(update.UpdateID)).Set(
			ctx,
			update,
		)
		if err != nil {
			log.Fatal(err)
			panic(1)
		}

		done_channel <- true
	}()

	return done_channel
}

func GetUser(
	ctx context.Context,
	firestore_client *firestore.Client,
	userID int64,
) chan *tgbotapi.User {
	// Get user by ID
	done_channel := make(chan *tgbotapi.User)

	go func() {
		user_map, err := firestore_client.Collection("users").Doc(fmt.Sprint(userID)).Get(ctx)
		fmt.Printf("Get user by ID: err=%v user=%v\n", err, user_map)
		if status.Code(err) == codes.NotFound {
			done_channel <- nil
		} else if err != nil {
			log.Fatal(err)
			panic(1)
		}

		var user tgbotapi.User
		user_map.DataTo(user)

		fmt.Printf("User struct is: %v", user)

		done_channel <- &user
	}()

	return done_channel
}

func CreateOrUpdateUser(
	ctx context.Context,
	firestore_client *firestore.Client,
	user *tgbotapi.User,
) chan bool {
	done_channel := make(chan bool)

	go func() {
		_, err := firestore_client.Collection("users").Doc(fmt.Sprint(user.ID)).Set(ctx, user)
		if err != nil {
			log.Fatal(err)
			panic(1)
		}
		done_channel <- true
	}()

	return done_channel
}

func AttachNewFormToUser(
	ctx context.Context,
	firestore_client *firestore.Client,
	user_id int64,
) chan *Form {
	result_channel := make(chan *Form)

	go func() {
		form, err := NewForm("form.yml", user_id)

		if err != nil {
			log.Fatal(err)
			panic(1)
		}

		doc_ref, _, err := firestore_client.Collection("forms").Add(ctx, form)

		if err != nil {
			log.Fatal(err)
			panic(1)
		}

		form.ID = doc_ref.ID
		result_channel <- form

	}()

	return result_channel
}

func GetUserForm(
	ctx context.Context,
	firestore_client *firestore.Client,
	user_id int64,
) chan *Form {
	result_channel := make(chan *Form)

	go func() {
		iter := firestore_client.Collection("forms").Where(
			"UserID", "==", user_id,
		).Where(
			"IsFilled", "==", false,
		).Documents(ctx)

		doc, err := iter.Next()

		if err == iterator.Done {
			result_channel <- nil
			return
		} else if err != nil {
			log.Fatal(err)
			panic(1)
		}

		var form Form
		err = doc.DataTo(&form)

		if err != nil {
			log.Fatal(err)
			panic(1)
		}

		form.ID = doc.Ref.ID
		result_channel <- &form
	}()

	return result_channel
}

func UpdateForm(
	ctx context.Context,
	firestore_client *firestore.Client,
	form *Form,
) chan bool {
	done_channel := make(chan bool)

	go func() {
		_, err := firestore_client.Collection("forms").Doc(form.ID).Set(ctx, form)
		if err != nil {
			log.Fatal(err)
			panic(1)
		}
		done_channel <- true
	}()

	return done_channel
}
