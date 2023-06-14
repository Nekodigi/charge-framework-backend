package firestore

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/Nekodigi/charge-framework-backend/models"
	"github.com/Nekodigi/charge-framework-backend/utils"
)

func (fs *Firestore) GetUserById(serviceId string, userId string) models.User {
	ctx := context.Background()
	doc, err := fs.Client.Collection("payment").Doc(serviceId).Collection("users").Doc(userId).Get(ctx)
	if err != nil {
		log.Println(err)
		return fs.CreateUser(serviceId, userId)
	}
	//bind to user
	var user models.User
	doc.DataTo(&user)
	user.ServiceId = serviceId
	return user
}

func (fs *Firestore) GetUserByIdTx(tx *firestore.Transaction, serviceId string, userId string) models.User {
	doc, err := tx.Get(fs.Client.Collection("payment").Doc(serviceId).Collection("users").Doc(userId))
	if err != nil {
		return fs.CreateUserTx(tx, serviceId, userId)
	}
	var user models.User
	doc.DataTo(&user)
	user.ServiceId = serviceId
	return user
}

func (fs *Firestore) GetUserBySubId(subscription_id string) models.User {
	ctx := context.Background()
	services, err := fs.Client.Collection("payment").Documents(ctx).GetAll()
	if err != nil {
		log.Fatalln(err)
	}
	var user models.User
	for _, service := range services {
		docs, _ := service.Ref.Collection("users").Where("subscription", "==", subscription_id).Documents(ctx).GetAll()
		if len(docs) > 0 {
			docs[0].DataTo(&user)
			break
		}
	}
	return user
}

func (fs *Firestore) GetUserBySubIdTx(tx *firestore.Transaction, subscription_id string) models.User {
	services, err := tx.Documents(fs.Client.Collection("payment")).GetAll()
	if err != nil {
		log.Fatalln(err)
	}
	var user models.User
	for _, service := range services {
		docs, _ := tx.Documents(service.Ref.Collection("users").Where("subscription", "==", subscription_id)).GetAll()
		if len(docs) > 0 {
			docs[0].DataTo(&user)
			break
		}
	}
	return user
}

func (fs *Firestore) UpdateUser(user models.User) {
	ctx := context.Background()
	_, err := fs.Client.Collection("payment").Doc(user.ServiceId).Collection("users").Doc(user.Id).Set(ctx, user)
	if err != nil {
		log.Fatalln(err)
	}
}

func (fs *Firestore) UpdateUserTx(tx *firestore.Transaction, user models.User) {
	err := tx.Set(fs.Client.Collection("payment").Doc(user.ServiceId).Collection("users").Doc(user.Id), user)
	if err != nil {
		log.Fatalln(err)
	}
}

func (fs *Firestore) NewUser(userId string, srv models.Service) models.User {
	user := models.User{
		Id:           userId,
		ServiceId:    srv.Id,
		Plan:         "free",
		Status:       "",
		AllocQuota:   srv.Plan["free"].Quota,
		RemainQuota:  srv.Plan["free"].Quota,
		Subscription: "",
		UpdateAt:     utils.EndOfThisMonth(),
	}
	return user
}

func (fs *Firestore) CreateUser(serviceId string, userId string) models.User {
	srv := fs.GetServiceById(serviceId)
	ctx := context.Background()
	user := fs.NewUser(userId, srv)
	_, err := fs.Client.Collection("payment").Doc(serviceId).Collection("users").Doc(userId).Set(ctx, user)
	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Printf("Created new user:%s", user.Id)
	}
	return user
}
func (fs *Firestore) CreateUserTx(tx *firestore.Transaction, serviceId string, userId string) models.User {
	srv := fs.GetServiceByIdTx(tx, serviceId)
	user := fs.NewUser(userId, srv)
	err := tx.Set(fs.Client.Collection("payment").Doc(serviceId).Collection("users").Doc(userId), user)
	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Printf("Created new user:%s", user.Id)
	}
	return user
}

func (fs *Firestore) UpdateUserQuota(user *models.User) bool {
	if time.Now().After(user.UpdateAt) {
		user.RemainQuota = user.AllocQuota
		user.UpdateAt = utils.EndOfThisMonth()
		fs.UpdateUser(*user)
		return true
	}
	return false
}
