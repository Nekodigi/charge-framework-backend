package firestore

import (
	"context"
	"log"
	"time"

	"github.com/Nekodigi/charge-framework-backend/models"
	"github.com/Nekodigi/charge-framework-backend/utils"
)

func (fs *Firestore) GetUserById(serviceId string, userId string) models.User {
	ctx := context.Background()
	doc, err := fs.client.Collection("payment").Doc(serviceId).Collection("users").Doc(userId).Get(ctx)
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

func (fs *Firestore) GetUserBySubId(subscription_id string) models.User {
	ctx := context.Background()
	services, err := fs.client.Collection("payment").Documents(ctx).GetAll()
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

func (fs *Firestore) UpdateUser(user models.User) {
	ctx := context.Background()
	_, err := fs.client.Collection("payment").Doc(user.ServiceId).Collection("users").Doc(user.Id).Set(ctx, user)
	if err != nil {
		log.Fatalln(err)
	}
}

func (fs *Firestore) CreateUser(serviceId string, userId string) models.User {
	srv := fs.GetServiceById(serviceId)
	ctx := context.Background()
	user := models.User{
		Id:           userId,
		ServiceId:    serviceId,
		Plan:         "free",
		Status:       "",
		AllocQuota:   srv.Plan["free"].Quota,
		RemainQuota:  srv.Plan["free"].Quota,
		Subscription: "",
		UpdateAt:     utils.EndOfThisMonth(),
	}
	_, err := fs.client.Collection("payment").Doc(serviceId).Collection("users").Doc(userId).Set(ctx, user)
	if err != nil {
		log.Fatalln(err)
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
