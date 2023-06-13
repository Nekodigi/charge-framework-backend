package firestore

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Nekodigi/charge-framework-backend/models"
	"github.com/Nekodigi/charge-framework-backend/utils"
)

func (fs *Firestore) GetServiceById(serviceId string) models.Service {
	ctx := context.Background()
	doc, err := fs.client.Collection("payment").Doc(serviceId).Get(ctx)
	if err != nil {
		fmt.Errorf("Error getting service: %v", err)
	}
	//bind to user
	var service models.Service
	// j, _ := json.Marshal(doc.Data())
	// fmt.Printf("Document data: %#v\n", string(j))
	// err = json.Unmarshal(j, &service)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	doc.DataTo(&service)
	//fmt.Printf("doc data: %#v\n", doc.Data())
	for key, plan := range doc.Data()["plan"].(map[string]interface{}) {
		//fmt.Printf("Plan data: %#v\n", key, plan)
		plan := plan.(map[string]interface{})
		var p models.Plan
		j, _ := json.Marshal(plan)
		_ = json.Unmarshal(j, &p)
		service.Plan[key] = p
	}
	service.UpdateAt = doc.Data()["updateAt"].(time.Time)
	//fmt.Printf("Service data: %#v\n", service)
	return service
}

func (fs *Firestore) UpdateService(service models.Service) {
	ctx := context.Background()
	_, err := fs.client.Collection("payment").Doc(service.Id).Set(ctx, service)
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
