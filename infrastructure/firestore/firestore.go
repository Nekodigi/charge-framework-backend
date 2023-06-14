package firestore

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

var (
	fs *Firestore
)

type (
	Firestore struct {
		app    *firebase.App
		Client *firestore.Client
	}
)

func NewFirestore() *Firestore {
	ctx := context.Background()
	//sa := option.WithCredentialsFile("credentials/serviceAccount.json")

	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	fs = &Firestore{
		app:    app,
		Client: client,
	}

	return fs
}
