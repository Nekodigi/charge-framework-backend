package firestore

import (
	"context"
	"fmt"
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

func (fs *Firestore) CreateUniqueID() string {
	ctx := context.Background()
	doc, _, err := fs.Client.Collection("unique_id").Add(ctx, map[string]interface{}{
		"creator": "service_framework",
	})
	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Printf("Created new id:%s\n")
	}
	return doc.ID
}
