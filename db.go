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

func createFirestoreClient(ctx context.Context) *firestore.Client {
	// Sets your Google Cloud Platform project ID.
	projectID := os.Getenv("PROJECT_ID")

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		panic(err)
	}

	return client
}

func save_update(
	ctx context.Context,
	firestore_client *firestore.Client,
	update *tgbotapi.Update,
) {
	_, err := firestore_client.Collection("updates").Doc(fmt.Sprint(update.UpdateID)).Set(
		ctx,
		update,
	)
	if err != nil {
		log.Fatal(err)
		panic(1)
	}
}

func GetUser(
	ctx context.Context,
	firestore_client *firestore.Client,
	userID int64,
) (*tgbotapi.User, error) {
	// Get user by ID
	user_map, err := firestore_client.Collection("users").Doc(fmt.Sprint(userID)).Get(ctx)
	fmt.Printf("Get user by ID: err=%v user=%v\n", err, user_map)
	if status.Code(err) == codes.NotFound {
		return nil, err
	} else if err != nil {
		log.Fatal(err)
		panic(1)
	}

	var user tgbotapi.User
	user_map.DataTo(&user)

	fmt.Printf("User struct is: %v", user)

	return &user, nil
}

func CreateOrUpdateUser(
	ctx context.Context,
	firestore_client *firestore.Client,
	user *tgbotapi.User,
) {
	_, err := firestore_client.Collection("users").Doc(fmt.Sprint(user.ID)).Set(ctx, user)
	if err != nil {
		log.Fatal(err)
		panic(1)
	}
}

func AttachNewFormToUser(
	ctx context.Context,
	firestore_client *firestore.Client,
	user_id int64,
) *Form {
	form, err := NewForm("form.yml", user_id)

	if err != nil {
		panic(err)
	}

	doc_ref, _, err := firestore_client.Collection("forms").Add(ctx, form)

	if err != nil {
		log.Fatal(err)
		panic(1)
	}

	form.ID = doc_ref.ID
	return form
}

func GetUserForm(
	ctx context.Context,
	firestore_client *firestore.Client,
	user_id int64,
) *Form {
	iter := firestore_client.Collection("forms").Where(
		"UserID", "==", user_id,
	).Where(
		"IsFilled", "==", false,
	).Documents(ctx)

	doc, err := iter.Next()

	if err == iterator.Done {
		return nil
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
	return &form
}

func UpdateForm(
	ctx context.Context,
	firestore_client *firestore.Client,
	form *Form,
) {
	_, err := firestore_client.Collection("forms").Doc(form.ID).Set(ctx, form)
	if err != nil {
		log.Fatal(err)
		panic(1)
	}
}
