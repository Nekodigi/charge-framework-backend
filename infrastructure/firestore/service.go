package firestore

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/Nekodigi/charge-framework-backend/models"
	"github.com/Nekodigi/charge-framework-backend/utils"
)

func DocToService(doc *firestore.DocumentSnapshot) models.Service {
	var service models.Service
	doc.DataTo(&service)
	//fmt.Printf("doc data: %#v\n", doc.Data())
	for key, plan := range doc.Data()["plan"].(map[string]interface{}) {
		//fmt.Printf("Plan data: %#v\n", key, plan)
		plan := plan.(map[string]interface{})
		var p models.Plan
		fmt.Println(plan)
		j, err := json.Marshal(plan)
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Println(j)
		err = json.Unmarshal(j, &p)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(p)
		service.Plan[key] = p
	}
	service.UpdateAt = doc.Data()["updateAt"].(time.Time)
	return service
}

func (fs *Firestore) GetServiceList() []models.Service {
	ctx := context.Background()
	iter := fs.Client.Collection("payment").Documents(ctx)
	var services []models.Service
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		services = append(services, DocToService(doc))
	}
	return services
}

func (fs *Firestore) GetServiceById(serviceId string) models.Service {
	ctx := context.Background()
	doc, err := fs.Client.Collection("payment").Doc(serviceId).Get(ctx)
	if err != nil {
		fmt.Errorf("Error getting service: %v", err)
	}
	return DocToService(doc)
}

func (fs *Firestore) GetServiceByIdTx(tx *firestore.Transaction, serviceId string) models.Service {
	doc, err := tx.Get(fs.Client.Collection("payment").Doc(serviceId))
	if err != nil {
		fmt.Errorf("Error getting service: %v", err)
	}
	fmt.Println(doc, serviceId)
	return DocToService(doc)
}

func (fs *Firestore) UpdateService(service models.Service) {
	ctx := context.Background()
	_, err := fs.Client.Collection("payment").Doc(service.Id).Set(ctx, service)
	if err != nil {
		log.Fatalln(err)
	}
}

func (fs *Firestore) UpdateServiceTx(tx *firestore.Transaction, service models.Service) {
	err := tx.Set(fs.Client.Collection("payment").Doc(service.Id), service)
	if err != nil {
		log.Fatalln(err)
	}
}

func (fs *Firestore) UpdateServiceQuota(service *models.Service) bool {
	if time.Now().After(service.UpdateAt) {
		if service.RemainQuota < service.AllocQuota {
			service.RemainQuota = service.AllocQuota
			service.UpdateAt = utils.EndOfThisMonth()
			fs.UpdateService(*service)
			return true
		}
	}
	return false
}
